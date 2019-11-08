package c2

import "github.com/kcarretto/paragon/api/codec"

// Server manages communication with agents.
type Server struct {
	*Queue

	OnAgentMessage func(codec.AgentMessage)
}
