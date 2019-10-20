package c2

import (
	"github.com/kcarretto/paragon/api/codec"
	"github.com/kcarretto/paragon/api/events"
)

// A Queue consumes and emits events to manage state for the C2 server.
type Queue struct {
	tasks []events.TaskQueued

	OnClaim func(events.TaskClaimed)
}

// ClaimTasks for an agent with the provided metadata.
func (q Queue) ClaimTasks(agent *codec.AgentMetadata) (tasks []*codec.Task) {
	if agent == nil {
		return
	}

	for i, event := range q.tasks {
		if sessionID := event.GetSessionID(); sessionID != "" && sessionID != agent.SessionID {
			continue
		}
		if machineUUID := event.GetMachineUUID(); machineUUID != "" && machineUUID != agent.MachineUUID {
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
func (q Queue) ConsumeTasks(tasks ...events.TaskQueued) {
	q.tasks = append(q.tasks, tasks...)
}

// pop a task from the queue by index, returns nil if the task does not exist.
func (q Queue) pop(index int) (task *codec.Task) {
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
