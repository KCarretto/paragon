package agent

import "errors"

// ErrAgentUninitialized occurs when attempting to use an agent method before initializing the agent.
// ErrNoTransportAvailable occurs when all available transports fail to report output, or if no transports were configured.
var (
	ErrAgentUninitialized   = errors.New("must initialize agent with agent.New()")
	ErrNoTransportAvailable = errors.New("no transport available")
)
