package teamserver

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gocloud.dev/pubsub"

	"github.com/golang/protobuf/proto"
	"github.com/kcarretto/paragon/api/codec"
	"github.com/kcarretto/paragon/api/events"
	"github.com/kcarretto/paragon/ent"
)

// Server handles c2 messages and replies with new tasks for the c2 to send out.
type Server struct {
	Log                  *zap.Logger
	EntClient            *ent.Client
	QueuedTopic          *pubsub.Topic
	ClaimedSubscription  *pubsub.Subscription
	ExecutedSubscription *pubsub.Subscription
}

type rawTask struct {
	Content string
}

func (srv *Server) handleQueueTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var t rawTask
		err := decoder.Decode(&t)
		if err != nil {
			http.Error(w, "improper task json sent", http.StatusBadRequest)
			return
		}
		ctx := context.Background()
		newTask, err := srv.EntClient.Task.Create().SetContent(t.Content).Save(ctx)
		if err != nil {
			http.Error(w, "unable to create new task", http.StatusInternalServerError)
			return
		}
		if err := srv.QueueTask(ctx, newTask); err != nil {
			http.Error(w, "unable to queue task to topic", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Run begins the handlers for processing the subscriptions to the `tasks.claimed` and `tasks.executed` topics
func (srv *Server) Run(ctx context.Context) {
	go srv.handleTasksClaimed(ctx)
	go srv.handleTasksExecuted(ctx)
	http.HandleFunc("/queueTask", srv.handleQueueTask)
	http.ListenAndServe(":80", nil)
}

// QueueTask sends a given task (and some associated target data) to the `tasks.queued` topic
func (srv *Server) QueueTask(ctx context.Context, task *ent.Task) error {
	target, err := task.QueryTarget().First(ctx)
	if err != nil {
		return err
	}
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

func (srv *Server) handleTasksClaimed(ctx context.Context) {
	for {
		msg, err := srv.ClaimedSubscription.Receive(ctx)
		if err != nil {
			panic(err)
		}
		var event events.TaskClaimed
		if err := proto.Unmarshal(msg.Body, &event); err != nil {
			srv.Log.Error("failed to parse protobuf", zap.Error(err))
		}
		taskID, err := strconv.Atoi(event.GetId())
		if err := proto.Unmarshal(msg.Body, &event); err != nil {
			srv.Log.Error("task id passed was not an int", zap.Error(err))
		}
		task, err := srv.EntClient.Task.Get(ctx, taskID)
		if err := proto.Unmarshal(msg.Body, &event); err != nil {
			srv.Log.Error("incorrect task id passed", zap.Error(err))
		}
		task, err = task.Update().
			SetClaimTime(time.Now()).
			Save(ctx)
		if err := proto.Unmarshal(msg.Body, &event); err != nil {
			srv.Log.Error("the task failed to be updated", zap.Error(err))
		}
		msg.Ack()
	}
}

func (srv *Server) handleTasksExecuted(ctx context.Context) {
	for {
		msg, err := srv.ClaimedSubscription.Receive(ctx)
		if err != nil {
			panic(err)
		}
		var event events.TaskExecuted
		if err := proto.Unmarshal(msg.Body, &event); err != nil {
			srv.Log.Error("failed to parse protobuff", zap.Error(err))
		}
		taskID, err := strconv.Atoi(event.GetId())
		if err := proto.Unmarshal(msg.Body, &event); err != nil {
			srv.Log.Error("task id passed was not an int", zap.Error(err))
		}
		task, err := srv.EntClient.Task.Get(ctx, taskID)
		if err := proto.Unmarshal(msg.Body, &event); err != nil {
			srv.Log.Error("incorrect task id passed", zap.Error(err))
		}
		execStartTime := time.Unix(event.GetExecStartTime(), 0)
		execStopTime := time.Unix(event.GetExecStopTime(), 0)
		output := []string{event.GetOutput()}
		task, err = task.Update().
			SetExecStartTime(execStartTime).
			SetExecStopTime(execStopTime).
			SetOutput(output).
			SetError(event.GetError()).
			Save(ctx)
		if err := proto.Unmarshal(msg.Body, &event); err != nil {
			srv.Log.Error("the task failed to be updated", zap.Error(err))
		}
		msg.Ack()
	}
}
