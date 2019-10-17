package c2_test

// import (
// 	"context"
// 	"encoding/json"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"go.uber.org/zap"

// 	"github.com/kcarretto/paragon/c2"
// 	"github.com/kcarretto/paragon/c2/mocks"
// 	"github.com/kcarretto/paragon/transport"
// 	"github.com/stretchr/testify/require"
// )

// func TestHandleMessage(t *testing.T) {
// 	expectedTask := transport.Task{ID: "HelloThere"}
// 	expectedReply := transport.Payload{
// 		Tasks: []transport.Task{expectedTask},
// 	}
// 	expected, err := json.Marshal(expectedReply)
// 	require.NoError(t, err)

// 	// Prepare mock
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	replyWriter := mocks.NewMockWriter(ctrl)
// 	replyWriter.EXPECT().Write(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
// 		require.Equal(t, string(expected), string(p))
// 		return len(p), nil
// 	})

// 	srv := &c2.Server{
// 		Log: zap.NewNop(),
// 	}
// 	srv.QueueTask(expectedTask, func(agent transport.Metadata) bool { return true })

// 	err = srv.HandleMessage(context.Background(), replyWriter, transport.Response{})
// 	require.NoError(t, err)
// }
