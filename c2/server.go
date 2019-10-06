package c2

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/kcarretto/paragon/transport"

	"go.uber.org/zap"
	"gocloud.dev/pubsub"
)

// Server handles agent messages and replies with new tasks for the agent to execute.
type Server struct {
	Logger      *zap.Logger
	AgentOutput *pubsub.Topic

	mu    sync.RWMutex
	tasks []queuedTask
}

// HandleMessage received from the agent, and write a reply to the provided writer.
func (srv *Server) HandleMessage(w io.Writer, msg transport.Response) error {
	// TODO: Get available tasks. srv.GetTasks(msg.Metadata)
	tasks := srv.GetTasks(msg.Metadata)
	reply := transport.Payload{
		Tasks: tasks,
	}

	data, err := json.Marshal(reply)
	if err != nil {
		return fmt.Errorf("failed to marshal server reply message: %w", err)
	}

	if _, err = w.Write(data); err != nil {
		return fmt.Errorf("failed to send reply message to agent: %w", err)
	}

	return nil
}
