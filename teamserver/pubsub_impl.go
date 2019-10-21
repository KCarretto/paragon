package teamserver

import (
	"context"
	"strconv"
	"time"

	"gocloud.dev/pubsub"

	"github.com/golang/protobuf/proto"
	"github.com/kcarretto/paragon/api/codec"
	"github.com/kcarretto/paragon/api/events"
	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/target"
)

func (srv *Server) agentCheckin(ctx context.Context, event events.AgentCheckin) error {
	targetExists, err := srv.EntClient.Target.Query().
		Where(
			target.MachineUUIDEQ(
				event.Agent.GetMachineUUID(),
			),
		).Exist(ctx)
	if err != nil {
		return err
	}

	if targetExists {
		targetEnt, err := srv.EntClient.Target.Query().
			Where(
				target.MachineUUIDEQ(
					event.Agent.MachineUUID,
				),
			).
			Only(ctx)
		if err != nil {
			return err
		}
		_, err = targetEnt.Update().SetLastSeen(time.Unix(event.GetSeenTime(), 0)).
			Save(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	// otherwise make a new target
	_, err = srv.EntClient.Target.
		Create().
		SetMachineUUID(event.Agent.MachineUUID).
		SetPublicIP(event.PublicIP).
		SetHostname(event.Agent.Hostname).
		SetPrimaryIP(event.Agent.PrimaryIP).
		SetPrimaryMAC(event.Agent.PrimaryMAC).
		SetLastSeen(time.Unix(event.SeenTime, 0)).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) queueTask(ctx context.Context, task *ent.Task) error {
	target := task.QueryTarget().FirstX(ctx)
	targetID := strconv.Itoa(target.ID)
	agentMetadata := codec.AgentMetadata{
		AgentID:     targetID,
		MachineUUID: target.MachineUUID,
		SessionID:   task.SessionID,
		Hostname:    target.Hostname,
		PrimaryIP:   target.PrimaryIP,
		PrimaryMAC:  target.PrimaryMAC,
	}
	taskID := strconv.Itoa(task.ID)
	event := events.TaskQueued{
		Id:      taskID,
		Content: task.Content,
		Filter:  &agentMetadata,
	}
	body, err := proto.Marshal(&event)
	if err != nil {
		return err
	}
	err = srv.QueuedTopic.Send(ctx, &pubsub.Message{
		Body: body,
	})
	if err != nil {
		return err
	}
	return nil
}

func (srv *Server) taskClaimed(ctx context.Context, event events.TaskClaimed) error {
	taskID, err := strconv.Atoi(event.GetId())
	if err != nil {
		return err
	}
	task, err := srv.EntClient.Task.Get(ctx, taskID)
	if err != nil {
		return err
	}
	task, err = task.Update().
		SetClaimTime(time.Now()).
		Save(ctx)
	if err != nil {
		return err
	}
	target := task.QueryTarget().FirstX(ctx)
	target, err = target.Update().
		SetLastSeen(time.Now()).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (srv *Server) taskExecuted(ctx context.Context, event events.TaskExecuted) error {
	taskID, err := strconv.Atoi(event.GetId())
	if err != nil {
		return err
	}
	task, err := srv.EntClient.Task.Get(ctx, taskID)
	if err != nil {
		return err
	}
	execStartTime := time.Unix(event.GetExecStartTime(), 0)
	execStopTime := time.Unix(event.GetExecStopTime(), 0)
	output := event.GetOutput()
	task, err = task.Update().
		SetExecStartTime(execStartTime).
		SetExecStopTime(execStopTime).
		SetOutput(output).
		SetError(event.GetError()).
		Save(ctx)
	if err != nil {
		return err
	}
	target := task.QueryTarget().FirstX(ctx)
	target, err = target.Update().
		SetLastSeen(time.Now()).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}
