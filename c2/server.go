package c2

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"sync"
// 	"time"

// 	"github.com/kcarretto/paragon/api/events"
// 	"github.com/kcarretto/paragon/transport"

// 	"github.com/golang/protobuf/proto"
// 	"go.uber.org/zap"
// 	"gocloud.dev/pubsub"
// )

// // Server handles agent messages and replies with new tasks for the agent to execute.
// type Server struct {
// 	Log         *zap.Logger
// 	TaskResults *pubsub.Topic

// 	mu    sync.RWMutex
// 	tasks []queueEntry
// }

// // HandleMessage received from the agent, and write a reply to the provided writer.
// func (srv *Server) HandleMessage(ctx context.Context, w io.Writer, msg transport.Response) error {
// 	srv.publishResults(ctx, msg)

// 	tasks := srv.GetTasks(msg.Metadata)
// 	reply := transport.Payload{
// 		Tasks: tasks,
// 	}

// 	data, err := json.Marshal(reply)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal server reply message: %w", err)
// 	}

// 	if _, err = w.Write(data); err != nil {
// 		return fmt.Errorf("failed to send reply message to agent: %w", err)
// 	}

// 	return nil
// }

// func (srv *Server) publishResults(ctx context.Context, msg transport.Response) {
// 	if srv.TaskResults == nil {
// 		srv.Log.Warn("No topic set for task results")
// 		return
// 	}

// 	for _, result := range msg.Results {
// 		event := events.TaskExecuted{
// 			Id:            result.ID,
// 			Output:        string(result.Output),
// 			Error:         result.Error,
// 			ExecStartTime: result.ExecStartTime.Unix(),
// 			ExecStopTime:  result.ExecStopTime.Unix(),
// 			RecvTime:      time.Now().Unix(),
// 		}

// 		// TODO: Send all results, or only ones with data?
// 		body, err := proto.Marshal(&event)
// 		if err != nil {
// 			srv.Log.Error("Failed to marshal result to json", zap.Error(err))
// 		}
// 		msg := &pubsub.Message{
// 			Body: body,
// 			// TODO: Add c2 metadata
// 		}
// 		go func() {
// 			if err = srv.TaskResults.Send(ctx, msg); err != nil {
// 				srv.Log.Error("Failed to publish task result", zap.Error(err))
// 			}
// 		}()
// 	}
// }
