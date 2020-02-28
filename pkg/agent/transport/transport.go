// Package transport defines messages sent between an agent and a server.
package transport

import (
	"context"
	"encoding/json"
	"fmt"
	io "io"
	"time"
)

//go:generate protoc -I=./proto -I=${GOPATH}/pkg/mod/github.com/gogo/googleapis@v1.3.0/ -I=${GOPATH}/pkg/mod/github.com/gogo/protobuf@v1.3.1/ --gogoslick_out=plugins=grpc,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api:. transport.proto

// ErrNoTransports occurs if no transports are available.
// ErrTransportsFailed occurs if all transports fail to send a message.
var (
	ErrNoTransports     = fmt.Errorf("cannot send message because no transports are configured")
	ErrTransportsFailed = fmt.Errorf("all available transports failed to send message")
)

// An AgentMessageWriter is responsible for writing messages from the agent to the server.
type AgentMessageWriter interface {
	WriteAgentMessage(context.Context, ServerMessageWriter, AgentMessage) error
}

// A ServerMessageWriter is responsible for consuming messages from the server to the agent.
type ServerMessageWriter interface {
	WriteServerMessage(context.Context, ServerMessage)
}

// A TaskExecutor is responsible for executing tasks and reporting output.
type TaskExecutor interface {
	ExecuteTask(ctx context.Context, output io.Writer, task *Task) error
}

// Transport is a base transport that should be embedded by other transports.
type Transport struct {
}

// EncodeAgentMessage to the provided writer.
func (t Transport) EncodeAgentMessage(msg AgentMessage, w io.Writer) (err error) {
	encoder := json.NewEncoder(w)
	err = encoder.Encode(msg)
	return
}

// DecodeAgentMessage from the provided reader.
func (t Transport) DecodeAgentMessage(data io.Reader) (msg AgentMessage, err error) {
	decoder := json.NewDecoder(data)
	err = decoder.Decode(&msg)
	return
}

// DecodeServerMessage from the provided reader.
func (t Transport) DecodeServerMessage(data io.Reader) (msg ServerMessage, err error) {
	decoder := json.NewDecoder(data)
	err = decoder.Decode(&msg)
	return
}

// EncodeServerMessage to the provided writer.
func (t Transport) EncodeServerMessage(msg ServerMessage, w io.Writer) (err error) {
	encoder := json.NewEncoder(w)
	err = encoder.Encode(msg)
	return
}

// AgentMessageMultiWriter attempts to write an agent message using the configured transports until
// one succeeds or no transports remain.
type AgentMessageMultiWriter struct {
	Transports []AgentMessageWriter
}

// WriteAgentMessage using configured transports (in order) until one succeeds or no transports remain.
func (w *AgentMessageMultiWriter) WriteAgentMessage(ctx context.Context, srv ServerMessageWriter, msg AgentMessage) error {
	if w.Transports == nil || len(w.Transports) < 1 {
		return ErrNoTransports
	}

	for _, transport := range w.Transports {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := transport.WriteAgentMessage(ctx, srv, msg); err != nil {
			continue
		}

		return nil
	}

	return ErrTransportsFailed
}

// CoerceStartTime returns the result's execution start time if available, nil otherwise.
func (result *TaskResult) CoerceStartTime() *time.Time {
	if result.ExecStartTime == nil {
		return nil
	}
	t := time.Unix(result.ExecStartTime.Seconds, int64(result.ExecStartTime.Nanos))
	return &t
}

// CoerceStopTime returns the result's execution stop time if available, nil otherwise.
func (result *TaskResult) CoerceStopTime() *time.Time {
	if result.ExecStopTime == nil {
		return nil
	}
	t := time.Unix(result.ExecStopTime.Seconds, int64(result.ExecStopTime.Nanos))
	return &t
}
