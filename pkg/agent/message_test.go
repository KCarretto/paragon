package agent_test

import (
	"testing"

	"github.com/kcarretto/paragon/pkg/agent"
	"github.com/kcarretto/paragon/pkg/c2"
	"github.com/stretchr/testify/require"
)

func TestMessageWrite(t *testing.T) {
	expected := "test message"

	msg := &agent.Message{}
	require.True(t, msg.IsEmpty())

	n, err := msg.Write([]byte(expected))
	require.NoError(t, err)
	require.Equal(t, len(expected), n)
	require.Equal(t, 1, len(msg.Logs))
	require.Equal(t, expected, msg.Logs[0])
}

func TestMessageWriteResult(t *testing.T) {
	expected := &c2.TaskResult{}

	msg := &agent.Message{}
	msg.WriteResult(expected)
	require.Equal(t, 1, len(msg.Results))
	require.Equal(t, expected, msg.Results[0])
}

func TestServerMessageWrite(t *testing.T) {
	expectedTask := &c2.Task{}
	expected := &agent.ServerMessage{}

	msg := &agent.ServerMessage{
		Tasks: []*c2.Task{expectedTask},
	}
	expected.WriteServerMessage(msg)

	require.NotNil(t, expected)
	require.Equal(t, 1, len(expected.Tasks))
	require.Equal(t, expectedTask, expected.Tasks[0])
}
