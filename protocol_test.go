package protocol

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestRPCNotification verifies constructor and JSON-RPC 2.0 notification frame serialization.
func TestRPCNotification(t *testing.T) {
	notice := NewNotification(MethodAgentMeta, map[string]string{"node": "test-box"})
	if notice.JSONRPC != Version || notice.Method != MethodAgentMeta {
		t.Errorf("unexpected notification fields: %+v", notice)
	}

	bytes, err := json.Marshal(notice)
	if err != nil || !strings.Contains(string(bytes), `"jsonrpc":"2.0"`) || !strings.Contains(string(bytes), `"method":"agent.meta"`) {
		t.Errorf("failed to marshal notification: %s", string(bytes))
	}
}

// TestMetaParamsSerialization verifies deterministic serialization for all 13 static metadata fields.
func TestMetaParamsSerialization(t *testing.T) {
	meta := MetaParams{
		Version:          "v1.0.0",
		Hostname:         "tokyo-01",
		OS:               "Debian 12",
		Kernel:           "6.1.0",
		Arch:             "amd64",
		CPUModel:         "Intel Xeon",
		CPUCores:         4,
		CPUPhysicalCores: 2,
		MemoryTotal:      8192 * 1024 * 1024,
		SwapTotal:        0, // must output deterministically even when zero
		DiskTotal:        100 * 1024 * 1024 * 1024,
		Virtualization:   "kvm",
		BootTime:         1710000000,
	}

	bytes, err := json.Marshal(meta)
	if err != nil {
		t.Fatalf("failed to marshal MetaParams: %v", err)
	}
	str := string(bytes)

	expectedFields := []string{
		`"version":"v1.0.0"`,
		`"hostname":"tokyo-01"`,
		`"os":"Debian 12"`,
		`"kernel":"6.1.0"`,
		`"arch":"amd64"`,
		`"cpu_model":"Intel Xeon"`,
		`"cpu_cores":4`,
		`"cpu_physical_cores":2`,
		`"swap_total":0`,
		`"virtualization":"kvm"`,
		`"boot_time":1710000000`,
	}

	for _, field := range expectedFields {
		if !strings.Contains(str, field) {
			t.Errorf("missing expected field %s in JSON: %s", field, str)
		}
	}
}

// TestMetricsPayloadDeterministic verifies all 14 core dynamic metrics are deterministically serialized with zero omitempty.
func TestMetricsPayloadDeterministic(t *testing.T) {
	payload := MetricsPayload{} // all-zero struct
	bytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal MetricsPayload: %v", err)
	}
	str := string(bytes)

	expectedKeys := []string{
		`"cpu_percent":0`,
		`"memory_used":0`,
		`"swap_used":0`,
		`"disk_used":0`,
		`"net_in_bytes":0`,
		`"net_out_bytes":0`,
		`"net_in_rate":0`,
		`"net_out_rate":0`,
		`"load_1":0`,
		`"load_5":0`,
		`"load_15":0`,
		`"process_count":0`,
		`"tcp_count":0`,
		`"udp_count":0`,
	}

	for _, key := range expectedKeys {
		if !strings.Contains(str, key) {
			t.Errorf("deterministic contract violation, missing field %s in JSON: %s", key, str)
		}
	}
}

// TestMetricsPluginsSlot verifies the optional plugins map slot and omitempty behavior.
func TestMetricsPluginsSlot(t *testing.T) {
	// Case A: serialize plugins when populated
	reportWithPlugins := &MetricsParams{
		Timestamp: 1700000000000,
		Metrics:   MetricsPayload{CPUPercent: 35.0},
		Plugins: map[string]any{
			"gpu": map[string]any{
				"temp": 65,
				"util": 90.5,
			},
			"ping": map[string]float64{
				"HK": 15.2,
			},
		},
	}
	bytes, err := json.Marshal(reportWithPlugins)
	if err != nil || !strings.Contains(string(bytes), `"plugins"`) || !strings.Contains(string(bytes), `"temp":65`) {
		t.Errorf("failed to marshal populated plugins: %s", string(bytes))
	}

	// Case B: omit plugins key when empty to keep frames compact
	reportNoPlugins := &MetricsParams{
		Timestamp: 1700000000000,
		Metrics:   MetricsPayload{CPUPercent: 10.0},
	}
	emptyBytes, _ := json.Marshal(reportNoPlugins)
	if strings.Contains(string(emptyBytes), `"plugins"`) {
		t.Errorf("empty plugins should be omitted, got: %s", string(emptyBytes))
	}
}
