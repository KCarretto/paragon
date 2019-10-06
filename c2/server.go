package c2

import (
	"context"
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
	Log         *zap.Logger
	TaskResults *pubsub.Topic

	mu    sync.RWMutex
	tasks []queuedTask
}

// HandleMessage received from the agent, and write a reply to the provided writer.
func (srv *Server) HandleMessage(ctx context.Context, w io.Writer, msg transport.Response) error {
	srv.publishResults(ctx, msg)

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

func (srv *Server) publishResults(ctx context.Context, msg transport.Response) {
	if srv.TaskResults == nil {
		srv.Log.Warn("No topic set for task results")
		return
	}

	for _, result := range msg.Results {
		// TODO: Send all results, or only ones with data?
		body, err := json.Marshal(result)
		if err != nil {
			srv.Log.Error("Failed to marshal result to json", zap.Error(err))
		}
		event := &pubsub.Message{
			Body: body,
			// TODO: Add agent metadata
		}
		go func() {
			if err = srv.TaskResults.Send(ctx, event); err != nil {
				srv.Log.Error("Failed to publish task result", zap.Error(err))
			}
		}()
	}
}
