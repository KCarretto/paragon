package teamserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"gocloud.dev/pubsub"

	"github.com/kcarretto/paragon/ent"
)

// Server handles c2 messages and replies with new tasks for the c2 to send out.
type Server struct {
	Log         *zap.Logger
	EntClient   *ent.Client
	QueuedTopic *pubsub.Topic
}

type rawTask struct {
	Content   string `json:"content"`
	SessionID string `json:"sessionID"`
	TargetID  int    `json:"targetID"`
}

type rawTarget struct {
	Name        string `json:"name"`
	MachineUUID string `json:"machineUUID"`
	Hostname    string `json:"hostname"`
	PrimaryIP   string `json:"primaryIP"`
	PrimaryMAC  string `json:"primaryMAC"`
}

type iDStruct struct {
	ID int `json:"id"`
}

type messageData struct {
	Data      []byte `json:"data"`
	MessageID string `json:"messageId"`
}

type pubSubMessage struct {
	Message      messageData `json:"message"`
	Subscription string      `json:"subscription"`
}

func (srv *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"status": "OK",
	}
	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "failed to marshal the json for the status", http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(resp); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write response data: %s", err.Error()), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
}

// Run begins the handlers for processing the subscriptions to the `tasks.claimed` and `tasks.executed` topics
func (srv *Server) Run() {
	http.HandleFunc("/status", srv.handleStatus)

	http.HandleFunc("/events/tasks/claimed", srv.handleTaskClaimed)
	http.HandleFunc("/events/tasks/executed", srv.handleTaskExecuted)

	http.HandleFunc("/makeTask", srv.handleMakeTask)
	http.HandleFunc("/makeTarget", srv.handleMakeTarget)
	http.HandleFunc("/getTask", srv.handleGetTask)
	http.HandleFunc("/getTarget", srv.handleGetTarget)
	http.HandleFunc("/listTargets", srv.handleListTargets)
	http.HandleFunc("/listTasksForTarget", srv.handleListTasksForTarget)
	if err := http.ListenAndServe("0.0.0.0:80", nil); err != nil {
		panic(err)
	}
}
