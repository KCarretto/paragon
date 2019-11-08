package agent

import (
	"context"
	"time"

	"github.com/kcarretto/paragon/api/codec"

	"go.uber.org/zap"
)

// A Sender is responsible for transporting messages to a server.
type Sender interface {
	Send(ServerMessageWriter, Message) error
}

// A Receiver is responsible for handling messages sent by a server.
type Receiver interface {
	Receive(MessageWriter, ServerMessage)
}

// An Agent communicates with server(s) using the configured transport.
type Agent struct {
	Receiver
	Log        *zap.Logger
	Metadata   *codec.AgentMetadata
	Transports []Transport

	MaxIdleTime time.Duration
	lastSend    time.Time
}

// Send messages to a server using the configured transports. Returns ErrNoTransports if all fail or
// if none are configured.
func (agent Agent) Send(w ServerMessageWriter, msg Message) error {
	// Don't send empty messages unless it has been at least MaxIdleTime since the last send.
	if msg.IsEmpty() && time.Since(agent.lastSend) <= agent.MaxIdleTime {
		return nil
	}

	agent.Log.Debug("Agent sending message", zap.Reflect("agent_msg", msg))

	// Attempt to send using available transports.
	for _, transport := range agent.Transports {
		if err := transport.Send(w, msg); err != nil {
			transport.Log.Error(
				"Failed to send message using transport",
				zap.Error(err),
				zap.Reflect("transport", transport),
			)
			continue
		}

		// When send is successful, update the timestamp
		agent.lastSend = time.Now()

		// Sleep the transport's interval on success
		transport.Sleep()

		return nil
	}

	return ErrNoTransports
}

// Run the agent, sending agent messages to a server using configured transports.
func (agent Agent) Run(ctx context.Context) error {
	agent.Log.Debug("Starting agent execution")
	agentMsg := Message{
		Metadata: agent.Metadata,
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var srvMsg ServerMessage

			if err := agent.Send(&srvMsg, agentMsg); err != nil {
				return err
			}
			agent.Log.Debug("Agent received message", zap.Reflect("srv_msg", srvMsg))

			agentMsg = Message{
				Metadata: agent.Metadata,
			}
			agent.Receive(&agentMsg, srvMsg)
		}
	}
}
