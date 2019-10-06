package c2_test

import (
	"testing"

	"github.com/kcarretto/paragon/c2"
	"github.com/kcarretto/paragon/transport"
	"github.com/stretchr/testify/require"
)

func TestQueueTask(t *testing.T) {
	expectedTask := transport.Task{ID: "TestID1"}

	srv := &c2.Server{}
	srv.QueueTask(
		expectedTask,
		func(a transport.Metadata) bool {
			return true
		},
	)

	require.Equal(t, 1, srv.TaskCount())

	tasks := srv.GetTasks(transport.Metadata{})
	require.Equal(t, 1, len(tasks))
	require.Equal(t, expectedTask, tasks[0])

	require.Equal(t, 0, srv.TaskCount())
}

func TestFilteredTask(t *testing.T) {
	expectedTask := transport.Task{ID: "TestID1"}
	unexpectedTask := transport.Task{ID: "Nope"}

	srv := &c2.Server{}
	srv.QueueTask(
		expectedTask,
		func(transport.Metadata) bool {
			return true
		},
	)
	require.Equal(t, 1, srv.TaskCount())

	srv.QueueTask(
		unexpectedTask,
		func(transport.Metadata) bool {
			return false
		},
	)
	require.Equal(t, 2, srv.TaskCount())

	tasks := srv.GetTasks(transport.Metadata{})
	require.Equal(t, 1, len(tasks))
	require.Equal(t, expectedTask, tasks[0])

	require.Equal(t, 1, srv.TaskCount())
}
