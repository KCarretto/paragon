package c2

import (
	"context"
	"fmt"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/graphql/models"
	"github.com/kcarretto/paragon/pkg/agent/transport"
	"go.uber.org/zap"
)

//go:generate protoc -I=./proto -I=${GOPATH}/pkg/mod/github.com/gogo/googleapis@v1.3.0/ -I=${GOPATH}/pkg/mod/github.com/gogo/protobuf@v1.3.1/ --gogoslick_out=plugins=grpc,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api:. transport.proto
//go:generate mockgen -destination=mocks/io.gen.go -package=mocks io Writer,WriteCloser
//go:generate mockgen -destination=mocks/teamserver.gen.go -package=mocks github.com/kcarretto/paragon/pkg/c2 Teamserver

// Teamserver provides client methods used to interact with a teamserver.
type Teamserver interface {
	ClaimTasks(context.Context, models.ClaimTasksRequest) ([]*ent.Task, error)
	SubmitTaskResult(context.Context, models.SubmitTaskResultRequest) error
}

// Server manages communication with agents. Upon check-in, the server will claim and respond with
// any tasks available for the agent.
type Server struct {
	Teamserver
	Log *zap.Logger
}

// WriteAgentMessage is a transport-agnostic method for handling communications from an agent.
func (srv Server) WriteAgentMessage(ctx context.Context, w transport.ServerMessageWriter, msg transport.AgentMessage) error {
	// Submit task results
	for _, result := range msg.Results {
		if result == nil {
			continue
		}

		if err := srv.SubmitTaskResult(ctx, models.SubmitTaskResultRequest{
			ID:            int(result.Id),
			Output:        &result.Output,
			Error:         &result.Error,
			ExecStartTime: result.CoerceStartTime(),
			ExecStopTime:  result.CoerceStopTime(),
		}); err != nil {
			srv.Log.Error("failed to submit task result", zap.Error(err), zap.Int64("task_id", result.Id))
		}
	}

	// Determine target criteria based on reported agent metadata
	target, err := resolveTarget(msg.Metadata)
	if err != nil {
		return err
	}

	// Claim tasks for the agent
	tasks, err := srv.ClaimTasks(ctx, target)
	if err != nil {
		return fmt.Errorf("failed to claim tasks from teamserver: %w", err)
	}

	srvMsg := transport.ServerMessage{
		Tasks: convertTasks(tasks),
	}

	w.WriteServerMessage(ctx, srvMsg)
	return nil
}

// convertTasks coerces an array of ent.Task to an array of transport.Task
func convertTasks(models []*ent.Task) (tasks []*transport.Task) {
	for _, task := range models {
		if task == nil {
			continue
		}

		tasks = append(tasks, &transport.Task{
			Id:      int64(task.ID),
			Content: task.Content,
		})
	}

	return tasks
}

// resolveTarget determines parameters necessary to claim tasks based on received agent metadata.
func resolveTarget(metadata *transport.AgentMetadata) (models.ClaimTasksRequest, error) {
	// Ensure Agent provided metadata to resolve a target
	if metadata == nil {
		return models.ClaimTasksRequest{}, fmt.Errorf("agent provided no valid metadata for target resolution")
	}

	// Otherwise combine Primary IP and Primary MAC addresses
	return models.ClaimTasksRequest{
		MachineUUID: &metadata.MachineUUID,
		SessionID:   &metadata.SessionID,
		Hostname:    &metadata.Hostname,
		PrimaryIP:   &metadata.PrimaryIP,
		PrimaryMac:  &metadata.PrimaryMAC,
	}, nil
}
