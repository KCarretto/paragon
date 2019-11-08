package c2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kcarretto/paragon/api/codec"
	"github.com/kcarretto/paragon/api/events"

	"github.com/golang/protobuf/proto"
)

type pubsubEvent struct {
	Message pubsubMsg `json:"message"`
}
type pubsubMsg struct {
	Data []byte `json:"data"`
}

// ServeStatus is an HTTP handler that returns the C2 Server status
func (srv *Server) ServeStatus(w http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{
		"status":         "OK",
		"tasks_in_queue": len(srv.tasks),
	}

	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Could not marshal status: %s", err.Error()),
			http.StatusInternalServerError,
		)
	}

	if _, err = w.Write(resp); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write response data: %s", err.Error()), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
}

//ServeHTTP implements http.Handler to handle agent messages sent via http.
func (srv *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Read request data
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Decode agent message
	var msg codec.AgentMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Prepare server message
	srvMsg := codec.ServerMessage{
		Tasks: srv.Queue.ClaimTasks(msg.Metadata),
	}

	// Encode server message
	resp, err := json.Marshal(srvMsg)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Write server message
	if _, err = w.Write(resp); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write response data: %s", err.Error()), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")

	// Handle task results
	if srv.OnAgentMessage == nil {
		return
	}
	srv.OnAgentMessage(msg)

}

// ServeEventTaskQueued is an HTTP handler for TaskQueued events.
func (srv *Server) ServeEventTaskQueued(w http.ResponseWriter, req *http.Request) {
	// Read request data
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Decode event message
	var msg pubsubEvent
	if err := json.Unmarshal(data, &msg); err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Decode event
	var event events.TaskQueued
	if err := proto.Unmarshal(msg.Message.Data, &event); err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal event: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Add task to queue
	srv.ConsumeTasks(event)
}

// ServeEventTaskClaimed is an HTTP handler for TaskClaimed events.
func (srv *Server) ServeEventTaskClaimed(w http.ResponseWriter, req *http.Request) {
	// Read request data
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Decode event message
	var msg pubsubEvent
	if err := json.Unmarshal(data, &msg); err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Decode event
	var event events.TaskClaimed
	if err := proto.Unmarshal(msg.Message.Data, &event); err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal event: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Remove task from queue
	srv.RemoveTask(event.Id)
}
