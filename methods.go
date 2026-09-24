package protocol

// RPC method names and protocol specifications.
const (
	// Version defines the current JSON-RPC protocol version.
	Version = "2.0"

	// MethodAgentMeta is reported once on handshake for static system metadata.
	MethodAgentMeta = "agent.meta"

	// MethodAgentMetrics is reported periodically for dynamic telemetry.
	MethodAgentMetrics = "agent.metrics"
)
