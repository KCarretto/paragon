package main

import (
	"context"
	"time"

	"github.com/kcarretto/paragon/api/codec"
	"github.com/kcarretto/paragon/api/events"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"gocloud.dev/pubsub"
)

func onClaim(ctx context.Context, logger *zap.Logger, claimTopic *pubsub.Topic) func(events.TaskClaimed) {
	return func(event events.TaskClaimed) {
		data, err := proto.Marshal(&event)
		if err != nil {
			logger.Error("Failed to marshal task claimed event", zap.Error(err))
		}
		if err := claimTopic.Send(ctx, &pubsub.Message{
			Body: data,
		}); err != nil {
			logger.Error("Failed to send task claimed event", zap.Error(err))
		}
	}
}

func onAgentMsg(ctx context.Context, logger *zap.Logger, checkinTopic *pubsub.Topic, execTopic *pubsub.Topic) func(codec.AgentMessage) {
	return func(msg codec.AgentMessage) {
		// Publish agent checkin
		recvTime := time.Now()
		checkinEvent := events.AgentCheckin{
			// TODO: PublicIP
			SeenTime: recvTime.Unix(),
			Agent:    msg.Metadata,
		}
		checkinData, err := proto.Marshal(&checkinEvent)
		if err != nil {
			logger.Error("Failed to marshal agent checkin event", zap.Error(err))
		}
		if err := checkinTopic.Send(ctx, &pubsub.Message{
			Body: checkinData,
		}); err != nil {
			logger.Error("Failed to send agent checkin event", zap.Error(err))
		}

		// Publish result events
		for _, result := range msg.GetResults() {
			event := events.TaskExecuted{
				Id:            result.GetId(),
				Output:        result.GetOutput(),
				Error:         result.GetError(),
				ExecStartTime: result.GetExecStartTime().GetSeconds(),
				ExecStopTime:  result.GetExecStopTime().GetSeconds(),
				RecvTime:      recvTime.Unix(),
			}
			data, err := proto.Marshal(&event)
			if err != nil {
				logger.Error("Failed to marshal task executed event", zap.Error(err))
				continue
			}
			if err := execTopic.Send(ctx, &pubsub.Message{
				Body: data,
			}); err != nil {
				logger.Error("Failed to send task executed event", zap.Error(err))
			}
		}

	}
}
