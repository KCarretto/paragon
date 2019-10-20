package main

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/kcarretto/paragon/api/events"
	"github.com/kcarretto/paragon/c2"
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

func listenForTasks(ctx context.Context, logger *zap.Logger, srv *c2.Server, queue *pubsub.Subscription) {
	logger.Info("Started listening to task queue")
	for {
		select {
		case <-ctx.Done():
			logger.Error("Finished listening to task queue", zap.Error(ctx.Err()))
			return
		default:
			msg, err := queue.Receive(ctx)
			if err != nil || msg == nil {
				logger.DPanic(
					"Failed to receive message from task queue",
					zap.Error(err),
					zap.Reflect("msg", msg),
				)
				continue
			}

			var event events.TaskQueued
			if err := proto.Unmarshal(msg.Body, &event); err != nil {
				logger.Error("Failed to unmarshal message from task queue", zap.Error(err))
				continue
			}

			logger.Debug("Received event from task queue", zap.Reflect("queue_event", event))
			srv.ConsumeTasks(event)
			msg.Ack()
		}
	}
}
