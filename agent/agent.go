package agent

import (
	"context"
	"fmt"

	"github.com/kcarretto/paragon/api/codec"
)

type AgentMessage codec.AgentMessage
type ServerMessage codec.ServerMessage

var ErrNoTransports = fmt.Errorf("all available transports failed to send message")

type ServerMessageWriter interface{}
type AgentMessageWriter interface{}

type Sender interface {
	Send(ServerMessageWriter, AgentMessage) error
}

type Receiver interface {
	Receive(AgentMessageWriter, ServerMessage)
}

type Agent struct {
	Sender
	Receiver

	Transports []Sender
}

func (agent Agent) Send(w ServerMessageWriter, msg AgentMessage) error {
	if agent.Sender != nil {
		return agent.Sender.Send(w, msg)
	}

	for _, transport := range agent.Transports {
		if err := transport.Send(w, msg); err != nil {
			// TODO: Error Handler
			continue
		}

		return nil
	}

	return ErrNoTransports
}

func (agent Agent) Run(ctx context.Context) error {

	var agentMsg AgentMessage

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var srvMsg ServerMessage
			if err := agent.Send(&srvMsg, agentMsg); err != nil {
				return err
			}

			agent.Receive(&agentMsg, srvMsg)
		}
	}
}
