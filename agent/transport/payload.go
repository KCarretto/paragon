package transport

// AgentPayload represents a common structure of data reported to the server from the agent. It is not
// required that all implementations use this structure, however the implementations provided by
// this repository will rely on it.
type AgentPayload struct {
	Output []byte `json:"output"`
}

// ServerPayload represents a common structure of data sent to the agent from the server. It is not
// required that all implementations use this structure, however the implementations provided by
// this repository will likely rely on it.
type ServerPayload struct {
	Tasks []TaskPayload `json:"tasks"`
}

// TaskPayload represents the structure of tasks included in a ServerPayload.
type TaskPayload struct {
	ID      string `json:"id"`
	Content []byte `json:"content"`
}
