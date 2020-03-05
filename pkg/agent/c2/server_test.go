package c2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/graphql/models"
	"github.com/kcarretto/paragon/pkg/agent/c2"
	"github.com/kcarretto/paragon/pkg/agent/transport"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func convertTasks(taskNodes ...*ent.Task) (tasks []*transport.Task) {
	for _, task := range taskNodes {
		tasks = append(tasks, &transport.Task{
			Id:      int64(task.ID),
			Content: task.Content,
		})
	}
	return
}

func TestWriteAgentMessage(t *testing.T) {
	expectedTasks := []*ent.Task{
		&ent.Task{
			ID:      1234,
			Content: "ABCDE",
		},
	}
	results := []*transport.TaskResult{
		&transport.TaskResult{
			Id:     1233,
			Output: "Some stuff",
		},
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// Prepare mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock teamserver
	teamserver := NewMockTeamserver(ctrl)
	teamserver.EXPECT().ClaimTasks(gomock.Any(), gomock.Any()).Return(expectedTasks, nil)
	teamserver.EXPECT().SubmitTaskResult(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, req models.SubmitTaskResultRequest) error {
		assert.Equal(t, 1233, req.ID, "submitted invalid response id")
		assert.Equal(t, "Some stuff", *req.Output, "submitted invalid response output")
		return fmt.Errorf("oh no, some error submitting results!")
	})
	// Ensure correct response is provided to agent
	w := NewMockServerMessageWriter(ctrl)
	w.EXPECT().WriteServerMessage(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, msg transport.ServerMessage) {
			assert.Equal(t, 1, len(msg.Tasks), "server returned invalid tasks")
			if len(msg.Tasks) < 1 {
				return
			}

			assert.Equal(t, int64(expectedTasks[0].ID), msg.Tasks[0].Id)
			assert.Equal(t, expectedTasks[0].Content, msg.Tasks[0].Content)
		})

	// Build server with mock dependencies
	srv := &c2.Server{
		Teamserver: teamserver,
		Log:        logger,
	}

	// WriteAgentMessage
	msg := transport.AgentMessage{
		Metadata: &transport.AgentMetadata{
			PrimaryIP: "127.0.0.1",
		},
		Results: results,
	}
	err = srv.WriteAgentMessage(context.Background(), w, msg)
	require.NoError(t, err)
}
