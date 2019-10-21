package teamserver

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/golang/protobuf/proto"
	"github.com/kcarretto/paragon/api/events"
)

func (srv *Server) handleAgentCheckin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var message pubSubMessage
		err := decoder.Decode(&message)
		if err != nil {
			http.Error(w, "improper pubsub message json sent", http.StatusBadRequest)
			return
		}
		var event events.AgentCheckin
		if err := proto.Unmarshal(message.Message.Data, &event); err != nil {
			srv.Log.Error("failed to parse protobuf", zap.Error(err))
		}
		ctx := context.Background()
		if err = srv.agentCheckin(ctx, event); err != nil {
			http.Error(w, "Failed to upset the agent checkin to the target", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "404 not found.", http.StatusNotFound)
}

func (srv *Server) handleTaskClaimed(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var message pubSubMessage
		err := decoder.Decode(&message)
		if err != nil {
			http.Error(w, "improper pubsub message json sent", http.StatusBadRequest)
			return
		}
		var event events.TaskClaimed
		if err := proto.Unmarshal(message.Message.Data, &event); err != nil {
			srv.Log.Error("failed to parse protobuf", zap.Error(err))
		}
		ctx := context.Background()
		err = srv.taskClaimed(ctx, event)
		if err != nil {
			http.Error(w, "an error occured in updating the claimed task", http.StatusInternalServerError)
			return
		}
		return
	}

	http.Error(w, "404 not found.", http.StatusNotFound)
}

func (srv *Server) handleTaskExecuted(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var message pubSubMessage
		err := decoder.Decode(&message)
		if err != nil {
			http.Error(w, "improper pubsub message json sent", http.StatusBadRequest)
			return
		}
		var event events.TaskExecuted
		if err := proto.Unmarshal(message.Message.Data, &event); err != nil {
			srv.Log.Error("failed to parse protobuf", zap.Error(err))
		}
		ctx := context.Background()
		err = srv.taskExecuted(ctx, event)
		if err != nil {
			http.Error(w, "an error occured in updating the executed task", http.StatusInternalServerError)
			return
		}
		return
	}

	http.Error(w, "404 not found.", http.StatusNotFound)
}
