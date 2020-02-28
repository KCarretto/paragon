package c2

import (
	"context"
	"fmt"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/graphql/models"
	"github.com/kcarretto/paragon/pkg/agent/transport"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=mockgen_teamserver_test.go -package=c2_test github.com/kcarretto/paragon/pkg/agent/c2 Teamserver
//go:generate mockgen -destination=mockgen_transport_test.go -package=c2_test github.com/kcarretto/paragon/pkg/agent/transport ServerMessageWriter

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
	var srvMsg transport.ServerMessage
	defer func() {
		w.WriteServerMessage(ctx, srvMsg)
	}()

	// Submit task results
	srv.submitResults(ctx, msg.Results...)

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

	srvMsg.Tasks = convertTasks(tasks)
	srv.Log.Error("HERE ARE THE TASKS", zap.Int("length", len(srvMsg.Tasks)), zap.Reflect("tasks", srvMsg.Tasks))
	fmt.Printf("TASKS: %d\n", len(srvMsg.Tasks))
	return nil
}

// Submit task results
func (srv Server) submitResults(ctx context.Context, results ...*transport.TaskResult) {
	for _, result := range results {
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
