package agent

import (
	"io"

	"github.com/kcarretto/paragon/api/codec"
)

// MessageWriter is responsible for writing output to be collected as a message from the agent.
type MessageWriter interface {
	io.Writer
	io.StringWriter
	WriteResult(*codec.Result)
}

// ServerMessageWriter is responsible for writing output to be collected as a message from the server.
type ServerMessageWriter interface {
	WriteServerMessage(*ServerMessage)
}

// Message is an alias for codec.AgentMessage with some extended functionality.
type Message codec.AgentMessage

// Write log output to be included in a message to a server.
func (msg *Message) Write(output []byte) (int, error) {
	return msg.WriteString(string(output))
}

// WriteString writes log output to be included in a message to a server.
func (msg *Message) WriteString(output string) (int, error) {
	msg.Logs = append(msg.Logs, output)
	return len(output), nil
}

// WriteResult writes execution output to be included in a message to a server.
func (msg *Message) WriteResult(result *codec.Result) {
	msg.Results = append(msg.Results, result)
}

// IsEmpty checks if the message has no significant contents
func (msg Message) IsEmpty() bool {
	return len(msg.Results) <= 0 && len(msg.Logs) <= 0
}

// ServerMessage is an alias for codec.ServerMessage with some extended functionality.
type ServerMessage codec.ServerMessage

// WriteServerMessage replaces this message with the provided message.
func (msg *ServerMessage) WriteServerMessage(srvMsg *ServerMessage) {
	msg.Tasks = append(msg.Tasks, srvMsg.Tasks...)
}
