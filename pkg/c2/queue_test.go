package c2_test

import (
	"testing"

	"github.com/kcarretto/paragon/api/codec"
	"github.com/kcarretto/paragon/api/events"
	"github.com/kcarretto/paragon/pkg/c2"
	"github.com/stretchr/testify/require"
)

func TestClaimTasks(t *testing.T) {
	agent := codec.AgentMetadata{
		MachineUUID: "A",
	}
	queue := &c2.Queue{}
	expectedTask := codec.Task{
		Id:      "Expected",
		Content: "Shuhmoopi",
	}
	queue.ConsumeTasks(events.TaskQueued{
		Id:      expectedTask.Id,
		Content: expectedTask.Content,
		Filter: &codec.AgentMetadata{
			MachineUUID: "A",
		},
	})

	tasks := queue.ClaimTasks(&agent)
	require.Equal(t, 1, len(tasks))
	require.NotNil(t, tasks[0])
	require.Equal(t, expectedTask, *tasks[0])

	moreTasks := queue.ClaimTasks(&agent)
	require.Equal(t, 0, len(moreTasks))
}
