package protocol

// MetaParams represents static hardware and system metadata reported via agent.meta.
type MetaParams struct {
	Version          string `json:"version"`            // Agent version (e.g. v1.0.0)
	Hostname         string `json:"hostname"`           // Node hostname
	OS               string `json:"os"`                 // OS distribution details (e.g. Debian 12)
	Kernel           string `json:"kernel"`             // OS kernel release version
	Arch             string `json:"arch"`               // CPU architecture (e.g. amd64, arm64)
	CPUModel         string `json:"cpu_model"`          // CPU model name
	CPUCores         int    `json:"cpu_cores"`          // Logical CPU cores / vCPUs
	CPUPhysicalCores int    `json:"cpu_physical_cores"` // Physical CPU cores count
	MemoryTotal      uint64 `json:"memory_total"`       // Total physical memory in bytes
	SwapTotal        uint64 `json:"swap_total"`         // Total swap space in bytes
	DiskTotal        uint64 `json:"disk_total"`         // Aggregated total disk storage in bytes
	Virtualization   string `json:"virtualization"`     // Virtualization type (kvm, docker, bare-metal, etc.)
	BootTime         uint64 `json:"boot_time"`          // System boot time in seconds (Unix epoch)
}

// MetricsParams represents the dynamic telemetry payload reported via agent.metrics.
type MetricsParams struct {
	Timestamp int64          `json:"timestamp"`         // Timestamp in Unix milliseconds
	Metrics   MetricsPayload `json:"metrics"`           // Core physical hardware metrics
	Plugins   map[string]any `json:"plugins,omitempty"` // Dynamic data reported by external plugins
}

// MetricsPayload contains the 14 core vitality dynamic metrics.
type MetricsPayload struct {
	CPUPercent   float64 `json:"cpu_percent"`   // Overall CPU utilization percentage (0.0 - 100.0)
	MemoryUsed   uint64  `json:"memory_used"`   // Memory used in bytes (cgroup aware)
	SwapUsed     uint64  `json:"swap_used"`     // Swap memory used in bytes
	DiskUsed     uint64  `json:"disk_used"`     // Aggregated physical disk used in bytes
	NetInBytes   uint64  `json:"net_in_bytes"`  // Total inbound network bytes since boot
	NetOutBytes  uint64  `json:"net_out_bytes"` // Total outbound network bytes since boot
	NetInRate    uint64  `json:"net_in_rate"`   // Instantaneous inbound rate in B/s
	NetOutRate   uint64  `json:"net_out_rate"`  // Instantaneous outbound rate in B/s
	Load1        float64 `json:"load_1"`        // 1-minute system load average
	Load5        float64 `json:"load_5"`        // 5-minute system load average
	Load15       float64 `json:"load_15"`       // 15-minute system load average
	ProcessCount int     `json:"process_count"` // Total active process count
	TCPCount     int     `json:"tcp_count"`     // Total active TCP connections count
	UDPCount     int     `json:"udp_count"`     // Total active UDP sockets count
}
