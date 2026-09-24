# Diting Protocol (`diting-protocol`)

[![Go Reference](https://pkg.go.dev/badge/github.com/diting-monitor/diting-protocol.svg)](https://pkg.go.dev/github.com/diting-monitor/diting-protocol)
[![Go Version](https://img.shields.io/github/go-mod/go-version/diting-monitor/diting-protocol)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Official JSON-RPC 2.0 protocol specifications, message envelopes, and strongly-typed data contracts for the **Diting (谛听)** host observability ecosystem.

---

## 💡 Design Philosophy

- **Zero External Dependencies**: Built strictly using the Go standard library (`encoding/json`, `fmt`, `testing`). No third-party bloat.
- **Deterministic Vitality Core**: All 14 core dynamic metrics are strictly deterministic (zero `omitempty`). Payload shapes remain consistent and predictable under all load conditions.
- **Single-Flight Handshake**: Static hardware and OS metadata are registered only once upon connection, decoupling dynamic telemetry from static facts.
- **In-Process Extensible Plugins**: Features an optional `map[string]any` slot for dynamic extensions (such as GPU telemetry, ping diagnostics, and custom collectors) without polluting core contracts.

---

## 📦 Installation

```bash
go get github.com/diting-monitor/diting-protocol
```

---

## 📡 Protocol Specification

Diting uses **JSON-RPC 2.0 Notifications** over persistent WebSocket connections. Because metrics reporting is a continuous unidirectional stream, notifications omit the `id` field and require no server ACK responses, reducing framing overhead by 50%.

### 1. Envelope (`Notification`)

```json
{
  "jsonrpc": "2.0",
  "method": "agent.metrics",
  "params": { ... }
}
```

```go
type Notification struct {
    JSONRPC string `json:"jsonrpc"` // Always "2.0"
    Method  string `json:"method"`  // e.g. "agent.meta" or "agent.metrics"
    Params  any    `json:"params"`
}
```

### 2. Static Metadata (`agent.meta`)

Reported once upon handshake or reconnect:

| Field | Type | Description |
| :--- | :--- | :--- |
| `version` | `string` | Agent binary version (e.g. `v1.0.0`) |
| `hostname` | `string` | Node hostname |
| `os` | `string` | Operating system details (e.g. `Debian GNU/Linux 12`) |
| `kernel` | `string` | Kernel release (e.g. `6.1.0-18-amd64`) |
| `arch` | `string` | System architecture (`amd64`, `arm64`, etc.) |
| `cpu_model` | `string` | CPU model string name |
| `cpu_cores` | `int` | Logical thread count (vCPU count) |
| `cpu_physical_cores` | `int` | Physical CPU core topology |
| `memory_total` | `uint64` | Total physical memory in bytes (cgroup aware) |
| `swap_total` | `uint64` | Total swap space in bytes |
| `disk_total` | `uint64` | Aggregated physical disk capacity in bytes |
| `virtualization` | `string` | Virtualization type (`kvm`, `docker`, `bare-metal`, etc.) |
| `boot_time` | `uint64` | System boot time in seconds (Unix epoch) |

### 3. Dynamic Telemetry (`agent.metrics`)

High-frequency telemetry reported periodically (default: every 2 seconds):

| Metric Field | Type | Description |
| :--- | :--- | :--- |
| `cpu_percent` | `float64` | Total CPU utilization percentage (`0.0` - `100.0`) |
| `memory_used` | `uint64` | Memory used in bytes (cgroup aware, excluding page cache) |
| `swap_used` | `uint64` | Swap space used in bytes |
| `disk_used` | `uint64` | Aggregated physical disk storage used in bytes |
| `net_in_bytes` | `uint64` | Total cumulative inbound network bytes since boot |
| `net_out_bytes` | `uint64` | Total cumulative outbound network bytes since boot |
| `net_in_rate` | `uint64` | Instantaneous inbound transfer rate in bytes/sec |
| `net_out_rate` | `uint64` | Instantaneous outbound transfer rate in bytes/sec |
| `load_1` | `float64` | 1-minute system load average |
| `load_5` | `float64` | 5-minute system load average |
| `load_15` | `float64` | 15-minute system load average |
| `process_count` | `int` | Count of active system processes |
| `tcp_count` | `int` | Count of active TCP connections |
| `udp_count` | `int` | Count of active UDP sockets |

#### Optional Plugins Slot (`plugins`)

Extensions can attach arbitrary telemetry via the `plugins` map:

```json
{
  "timestamp": 1710000000000,
  "metrics": { ... },
  "plugins": {
    "gpu": { "util": 85.5, "temp": 68 },
    "ping": { "hk": 18.2, "us": 142.0 }
  }
}
```

---

## 🛠️ Quick Example

### Agent: Creating and Sending a Notification

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/diting-monitor/diting-protocol"
)

func main() {
    params := protocol.MetricsParams{
        Timestamp: 1710000000000,
        Metrics: protocol.MetricsPayload{
            CPUPercent: 12.5,
            MemoryUsed: 1024 * 1024 * 1024,
            TCPCount:   128,
            UDPCount:   16,
        },
    }

    // Wrap into a standard JSON-RPC 2.0 notification frame
    notice := protocol.NewNotification(protocol.MethodAgentMetrics, params)
    data, _ := json.Marshal(notice)

    fmt.Println(string(data))
    // Output: {"jsonrpc":"2.0","method":"agent.metrics","params":{...}}
}
```

### Server: Handling Incoming Frames

```go
type InboundFrame struct {
    JSONRPC string          `json:"jsonrpc"`
    Method  string          `json:"method"`
    Params  json.RawMessage `json:"params"`
}

func handleMessage(msgBytes []byte) error {
    var frame InboundFrame
    if err := json.Unmarshal(msgBytes, &frame); err != nil {
        return err
    }

    switch frame.Method {
    case protocol.MethodAgentMeta:
        var meta protocol.MetaParams
        return json.Unmarshal(frame.Params, &meta)
    case protocol.MethodAgentMetrics:
        var metrics protocol.MetricsParams
        return json.Unmarshal(frame.Params, &metrics)
    }
    return nil
}
```

---

## 📄 License

This module is licensed under the [MIT License](LICENSE).
