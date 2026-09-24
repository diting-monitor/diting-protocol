package protocol

// Notification represents a standard JSON-RPC 2.0 notification (one-way, no id, no ack).
type Notification struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

// NewNotification creates a new JSON-RPC 2.0 notification frame.
func NewNotification(method string, params any) Notification {
	return Notification{
		JSONRPC: Version,
		Method:  method,
		Params:  params,
	}
}
