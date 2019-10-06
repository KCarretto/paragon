package c2

import (
	"github.com/kcarretto/paragon/transport"
)

type queuedTask struct {
	task   transport.Task
	filter func(transport.Metadata) bool
}

// QueueTask prepares a task to be sent to the first agent that reports metadata matching the filter.
func (srv *Server) QueueTask(task transport.Task, filter func(transport.Metadata) bool) {
	srv.mu.Lock()
	srv.tasks = append(srv.tasks, queuedTask{
		task,
		filter,
	})
	srv.mu.Unlock()
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
