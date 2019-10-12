package c2

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/kcarretto/paragon/api/events"
	"github.com/kcarretto/paragon/transport"

	"go.uber.org/zap"
	"gocloud.dev/pubsub"
)

type queueEntry struct {
	task   transport.Task
	filter func(transport.Metadata) bool
}

// TODO: Convert underlying struct to map

// QueueTask prepares a task to be sent to the first agent that reports metadata matching the filter.
func (srv *Server) QueueTask(task transport.Task, filter func(transport.Metadata) bool) {
	srv.mu.Lock()
	srv.tasks = append(srv.tasks, queueEntry{
		task,
		filter,
	})
	srv.mu.Unlock()
}

// removeTask removes the first task from the queue with a matching ID. If no matching task is found
// in the queue, removeTask is a no-op.
func (srv *Server) removeTask(id string) {
	srv.mu.RLock()
	var index = -1
	for i, entry := range srv.tasks {
		if entry.task.ID == id {
			index = i
			break
		}
	}
	srv.mu.RUnlock()

	if index < 0 {
		return
	}

	// Remove element from the queue
	srv.tasks = append(srv.tasks[:index], srv.tasks[index+1:]...)

}

// GetTasks for an agent based on the metadata it reported and the current tasks in the server queue.
func (srv *Server) GetTasks(agent transport.Metadata) []transport.Task {
	srv.mu.Lock()
	n := 0
	tasks := make([]transport.Task, 0, len(srv.tasks))
	for _, t := range srv.tasks {
		if t.filter == nil || t.filter(agent) {
			tasks = append(tasks, t.task)
		} else {
			srv.tasks[n] = t
			n++
		}
	}
	srv.tasks = srv.tasks[:n]
	srv.mu.Unlock()

	return tasks
}

// TaskCount returns the total number of tasks left in the queue.
func (srv *Server) TaskCount() int {
	srv.mu.RLock()
	count := len(srv.tasks)
	srv.mu.RUnlock()

	return count
}

// ListenForQueuedTasks receives TaskQueued events and queues the corresponding task.
func (srv *Server) ListenForQueuedTasks(ctx context.Context, subscription *pubsub.Subscription) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := subscription.Receive(ctx)
			if err != nil {
				srv.Log.Error("Subscription failed", zap.Error(err))
				return
			}

			var event events.TaskQueued
			if err = proto.Unmarshal(msg.Body, &event); err != nil {
				srv.Log.Error("Failed to unmarshal TaskQueued event", zap.Error(err))
			}

			criteria := transport.Metadata{
				AgentID:     event.Filter.GetAgentID(),
				MachineUUID: event.Filter.GetMachineUUID(),
				SessionID:   event.Filter.GetSessionID(),
				Hostname:    event.Filter.GetHostname(),
				PrimaryIP:   event.Filter.GetPrimaryIP(),
				PrimaryMAC:  event.Filter.GetPrimaryMAC(),
			}

			task := transport.Task{
				ID:      event.GetId(),
				Content: []byte(event.GetContent()),
			}

			srv.QueueTask(task, claimFilter(criteria))
		}
	}
}

// ListenForClaimedTasks receives TaskClaimed events and dequeues the corresponding task.
func (srv *Server) ListenForClaimedTasks(ctx context.Context, subscription *pubsub.Subscription) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := subscription.Receive(ctx)
			if err != nil {
				srv.Log.Error("Subscription failed", zap.Error(err))
				return
			}

			var event events.TaskClaimed
			if err = proto.Unmarshal(msg.Body, &event); err != nil {
				srv.Log.Error("Failed to unmarshal TaskClaimed event", zap.Error(err))
			}

			srv.removeTask(event.GetId())
		}
	}
}

// claimFilter returns a filter that implements a simple algorithm for task claiming.
func claimFilter(criteria transport.Metadata) func(agent transport.Metadata) bool {
	return func(agent transport.Metadata) bool {
		// True if session ids match.
		if id := criteria.SessionID; id != "" {
			if id == agent.SessionID {
				return true
			}
		}

		// Restrict if agentID criteria is set.
		if agentID := criteria.AgentID; agentID != "" {
			if agentID != agent.AgentID {
				return false
			}
		}

		// True if machine uuids match.
		if uuid := criteria.MachineUUID; uuid != "" {
			if uuid == agent.MachineUUID {
				return true
			}
		}

		// True if MAC address matches
		if mac := criteria.PrimaryMAC; mac != "" {
			if mac == agent.PrimaryMAC {
				return true
			}
		}

		// True if IP address matches
		if addr := criteria.PrimaryIP; addr != "" {
			if addr == agent.PrimaryIP {
				return true
			}
		}

		// True if hostname matches
		if hostname := criteria.Hostname; hostname != "" {
			if hostname == agent.Hostname {
				return true
			}
		}

		// False otherwise
		return false
	}
}
