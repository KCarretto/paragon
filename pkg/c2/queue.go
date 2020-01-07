package c2

import (
	"sync"

	"github.com/kcarretto/paragon/proto/codec"
	"github.com/kcarretto/paragon/proto/events"
)

// A Queue consumes and emits events to manage state for the C2 server.
type Queue struct {
	tasks []events.TaskQueued

	mu      sync.Mutex
	OnClaim func(events.TaskClaimed)
}

// ClaimTasks for an agent with the provided metadata.
func (q *Queue) ClaimTasks(agent *codec.AgentMetadata) (tasks []*codec.Task) {
	if agent == nil {
		return
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	for i, event := range q.tasks {
		filter := event.GetFilter()
		if filter == nil {
			continue
		}

		if sessionID := filter.GetSessionID(); sessionID != "" && sessionID != agent.SessionID {
			continue
		}
		if machineUUID := filter.GetMachineUUID(); machineUUID != "" && machineUUID != agent.MachineUUID {
			continue
		}

		task := q.pop(i)
		if task == nil {
			continue
		}

		tasks = append(tasks, task)

		if q.OnClaim == nil {
			continue
		}
		q.OnClaim(events.TaskClaimed{
			Id:    task.GetId(),
			Agent: agent,
		})
	}

	return tasks
}

// ConsumeTasks adds tasks to the queue from queued events.
func (q *Queue) ConsumeTasks(tasks ...events.TaskQueued) {
	q.mu.Lock()
	q.tasks = append(q.tasks, tasks...)
	q.mu.Unlock()
}

// RemoveTask removes a task from the queue based on ID. If no task exists, RemoveTask is a no-op.
func (q *Queue) RemoveTask(id string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for i, event := range q.tasks {
		if event.GetId() == id {
			q.tasks = append(q.tasks[:i], q.tasks[i+1:]...)
			return
		}
	}
}

// pop a task from the queue by index, returns nil if the task does not exist.
func (q *Queue) pop(index int) (task *codec.Task) {
	if index < 0 || index >= len(q.tasks) {
		return nil
	}

	event := q.tasks[index]

	task = &codec.Task{
		Id:      event.GetId(),
		Content: event.GetContent(),
	}

	q.tasks = append(q.tasks[:index], q.tasks[index+1:]...)

	return task
}
