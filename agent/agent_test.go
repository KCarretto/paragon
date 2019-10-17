package agent_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kcarretto/paragon/agent"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSend(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	badSender := NewMockSender(ctrl)
	sender := NewMockSender(ctrl)
	unusedSender := NewMockSender(ctrl)

	srvMsg := agent.ServerMessage{}
	agentMsg := agent.Message{}
	badSender.EXPECT().Send(&srvMsg, agentMsg).Return(errors.New("oops ^_^"))
	sender.EXPECT().Send(&srvMsg, agentMsg).Return(nil)

	testAgent := &agent.Agent{
		Log: logger,
		Transports: []agent.Transport{
			agent.Transport{
				Sender: badSender,
				Log:    logger.Named("transport.bad"),
				Name:   "Bad Sender",
			},
			agent.Transport{
				Sender: sender,
				Log:    logger.Named("transport.good"),
				Name:   "Good Sender",
			},
			agent.Transport{
				Sender: unusedSender,
				Log:    logger.Named("transport.should_not_see"),
				Name:   "Unusued Sender",
			},
		},
	}
	err = testAgent.Send(&srvMsg, agentMsg)
	require.NoError(t, err)
}

func TestAgentRun(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sender := NewMockSender(ctrl)
	recv := NewMockReceiver(ctrl)

	sender.EXPECT().Send(gomock.Not(gomock.Nil()), agent.Message{}).Return(nil).AnyTimes()
	recv.EXPECT().Receive(gomock.Not(gomock.Nil()), agent.ServerMessage{}).AnyTimes()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	testAgent := &agent.Agent{
		Log:         logger,
		Receiver:    recv,
		MaxIdleTime: time.Second * 1,
		Transports: []agent.Transport{
			agent.Transport{
				Sender: sender,
				Log:    logger.Named("transport"),
				Name:   "Test Sender",
			},
		},
	}
	err = testAgent.Run(ctx)
	require.True(t, errors.Is(err, context.DeadlineExceeded), "Unexpected error: %v", err)
}

func TestRunErrNoTransport(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srvMsg := agent.ServerMessage{}
	agentMsg := agent.Message{}
	badSender := NewMockSender(ctrl)
	badSender.EXPECT().Send(&srvMsg, agentMsg).Return(errors.New("oops ^_^"))

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	testAgent := &agent.Agent{
		Log: logger,
		Transports: []agent.Transport{
			agent.Transport{
				Sender: badSender,
				Log:    logger.Named("transport.bad"),
				Name:   "Bad Sender",
			},
		},
	}

	err = testAgent.Run(ctx)
	require.True(t, errors.Is(err, agent.ErrNoTransports), "Unexpected error: %v", err)
}
