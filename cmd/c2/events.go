package main

import (
	"context"

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

func onExec(ctx context.Context, logger *zap.Logger, execTopic *pubsub.Topic) func(events.TaskExecuted) {
	return func(event events.TaskExecuted) {
		data, err := proto.Marshal(&event)
		if err != nil {
			logger.Error("Failed to marshal task executed event", zap.Error(err))
		}
		if err := execTopic.Send(ctx, &pubsub.Message{
			Body: data,
		}); err != nil {
			logger.Error("Failed to send task executed event", zap.Error(err))
		}
	}
}
