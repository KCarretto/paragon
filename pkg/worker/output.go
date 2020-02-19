package worker

import (
	"context"
	"log"
	"time"

	"github.com/kcarretto/paragon/graphql"
	"github.com/kcarretto/paragon/graphql/models"
)

// taskOutput acts like a stream buffer which submits task execution results as they are written.
type taskOutput struct {
	ID    int
	Ctx   context.Context
	Graph graphql.Client

	id  int
	err *string
}

// Start submits an initial result for the task, indicating it's execution start time.
func (t *taskOutput) Start() {
	now := time.Now()

	if err := t.Graph.SubmitTaskResult(t.Ctx, models.SubmitTaskResultRequest{
		ID:            t.ID,
		ExecStartTime: &now,
	}); err != nil {
		log.Printf("[ERR] Failed to submit task result for execution start: %s", err.Error())
	}
}

// Stop submits the final result for the task, indicating it's execution stop time and any errors.
func (t *taskOutput) Stop(execErr error) {
	now := time.Now()

	if execErr != nil {
		errMsg := execErr.Error()
		t.err = &errMsg
	}

	if err := t.Graph.SubmitTaskResult(t.Ctx, models.SubmitTaskResultRequest{
		ID:           t.ID,
		Error:        t.err,
		ExecStopTime: &now,
	}); err != nil {
		log.Printf("[ERR] Failed to submit task result for execution stop: %s", err.Error())
	}
}

// Write submits task output.
func (t *taskOutput) Write(p []byte) (int, error) {
	output := string(p)
	if err := t.Graph.SubmitTaskResult(t.Ctx, models.SubmitTaskResultRequest{
		ID:     t.ID,
		Output: &output,
	}); err != nil {
		log.Printf("[ERR] Failed to submit task execution result: %s", err.Error())
		return 0, err
	}

	return len(p), nil
}
