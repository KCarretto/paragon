package agent

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kcarretto/paragon/pkg/agent/transport"

	types "github.com/gogo/protobuf/types"
	"go.uber.org/zap"
)

// Agent provides a standard flow for receiving tasks, writing results, and sending output.
type Agent struct {
	transport.TaskExecutor
	transport.AgentMessageWriter

	OnRun		func()

	Log         *zap.Logger
	Metadata    transport.AgentMetadata
	MaxIdleTime time.Duration

	wg sync.WaitGroup

	machineidPrefix string
}

// Run the agent, which will block until the provided context has been canceled.
func (agent *Agent) Run(ctx context.Context) {
	agent.collectMetadata()

	if agent.OnRun != nil{
		agent.OnRun()
	}

	checkinTicker := time.NewTicker(agent.MaxIdleTime)
	defer checkinTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			agent.wg.Wait()
			return
		case <-checkinTicker.C:
			agent.sendMessage(ctx, transport.AgentMessage{Metadata: &agent.Metadata})
			break
		}
	}
}

// WriteServerMessage will spawn go routines to execute each task provided by the server.
func (agent *Agent) WriteServerMessage(ctx context.Context, msg transport.ServerMessage) {
	for _, t := range msg.Tasks {
		if t == nil {
			continue
		}

		agent.wg.Add(1)
		go func(task *transport.Task) {
			defer agent.wg.Done()

			result := agent.runTask(ctx, task)

			agent.sendMessage(ctx, transport.AgentMessage{
				Metadata: &agent.Metadata,
				Results: []*transport.TaskResult{
					&result,
				},
			})
		}(t)
	}
}

func (agent *Agent) sendMessage(ctx context.Context, msg transport.AgentMessage) {
	if err := agent.WriteAgentMessage(ctx, agent, msg); err != nil {
		agent.Log.Error("failed to send agent message",
			zap.Error(err),
		)
		return
	}
	agent.Log.Debug("successfully sent agent message")
}

func (agent *Agent) runTask(ctx context.Context, task *transport.Task) (result transport.TaskResult) {
	// Set result ID
	result.Id = task.Id

	// Set execution start time
	start := time.Now()
	result.ExecStartTime = &types.Timestamp{
		Seconds: start.Unix(),
		Nanos:   int32(start.Nanosecond()),
	}

	// Set error if task execution causes panic
	defer func() {
		if err := recover(); err != nil {
			result.Error = fmt.Sprintf("task execution resulting in panic: %v", err)
		}
	}()

	// Set output after task execution is finished
	output := new(bytes.Buffer)
	defer func() {
		result.Output = output.String()
	}()

	// Set execution stop time after task execution is finished
	defer func() {
		stop := time.Now()
		result.ExecStopTime = &types.Timestamp{
			Seconds: stop.Unix(),
			Nanos:   int32(stop.Nanosecond()),
		}
	}()

	// Execute the task, setting the error if necessary
	if err := agent.ExecuteTask(ctx, output, task); err != nil {
		result.Error = err.Error()
	}

	agent.Log.Debug("Agent completed task execution", zap.Int64("task_id", task.Id))
	return
}
