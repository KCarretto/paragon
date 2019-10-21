package teamserver

import (
	"context"
	"encoding/json"
	"net/http"
)

func (srv *Server) handleMakeTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var t rawTask
		err := decoder.Decode(&t)
		if err != nil {
			http.Error(w, "improper task json sent", http.StatusBadRequest)
			return
		}
		ctx := context.Background()
		target, err := srv.EntClient.Target.Get(ctx, t.TargetID)
		if err != nil {
			http.Error(w, "improper target id given", http.StatusBadRequest)
			return
		}
		newTask, err := srv.EntClient.Task.
			Create().
			SetContent(t.Content).
			SetSessionID(t.SessionID).
			SetTarget(target).
			Save(ctx)
		if err != nil {
			http.Error(w, "unable to create new task", http.StatusInternalServerError)
			return
		}
		resp := iDStruct{
			ID: newTask.ID,
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "unable to create response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(respBody)
		if err := srv.queueTask(ctx, newTask); err != nil {
			http.Error(w, "unable to queue task to topic", http.StatusInternalServerError)
			return
		}
		return
	}

	http.Error(w, "404 not found.", http.StatusNotFound)
}

func (srv *Server) handleGetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var message iDStruct
		err := decoder.Decode(&message)
		if err != nil {
			http.Error(w, "improper id struct sent", http.StatusBadRequest)
			return
		}
		ctx := context.Background()
		task, err := srv.EntClient.Task.Get(ctx, message.ID)

		data := map[string]interface{}{
			"id":            task.ID,
			"queueTime":     task.QueueTime.Unix(),
			"claimTime":     task.ClaimTime.Unix(),
			"execStartTime": task.ExecStartTime.Unix(),
			"execStopTime":  task.ExecStopTime.Unix(),
			"content":       task.Content,
			"output":        task.Output,
			"error":         task.Error,
			"sessionID":     task.SessionID,
		}
		resp, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "failed to marshal the json for the task", http.StatusInternalServerError)
			return
		}
		if _, err = w.Write(resp); err != nil {
			http.Error(w, "Failed to write response data", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		return
	}

	http.Error(w, "404 not found.", http.StatusNotFound)
}

func (srv *Server) handleGetTarget(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var message iDStruct
		err := decoder.Decode(&message)
		if err != nil {
			http.Error(w, "improper id struct sent", http.StatusBadRequest)
			return
		}
		ctx := context.Background()
		target, err := srv.EntClient.Target.Get(ctx, message.ID)

		data := map[string]interface{}{
			"id":          target.ID,
			"name":        target.Name,
			"machineUUID": target.MachineUUID,
			"primaryIP":   target.PrimaryIP,
			"primaryMAC":  target.PrimaryMAC,
			"hostname":    target.Hostname,
			"lastSeen":    target.LastSeen,
		}
		resp, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "failed to marshal the json for the target", http.StatusInternalServerError)
			return
		}
		if _, err = w.Write(resp); err != nil {
			http.Error(w, "Failed to write response data", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		return
	}

	http.Error(w, "404 not found.", http.StatusNotFound)
}

func (srv *Server) handleListTargets(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	targets, err := srv.EntClient.Target.Query().All(ctx)
	if err != nil {
		http.Error(w, "failed to query the ent db for targets", http.StatusInternalServerError)
		return
	}
	var targetIDs []int
	for _, target := range targets {
		targetIDs = append(targetIDs, target.ID)
	}
	data := map[string]interface{}{
		"ids": targetIDs,
	}
	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "failed to marshal the json for the ids", http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(resp); err != nil {
		http.Error(w, "Failed to write response data", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
}

func (srv *Server) handleListTasksForTarget(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var message iDStruct
		err := decoder.Decode(&message)
		if err != nil {
			http.Error(w, "improper id struct sent", http.StatusBadRequest)
			return
		}
		ctx := context.Background()
		target, err := srv.EntClient.Target.Get(ctx, message.ID)
		if err != nil {
			http.Error(w, "failed to lookup target with given id", http.StatusBadRequest)
			return
		}

		tasks, err := target.QueryTasks().All(ctx)
		var taskIDs []int
		for _, task := range tasks {
			taskIDs = append(taskIDs, task.ID)
		}
		data := map[string]interface{}{
			"ids": taskIDs,
		}
		resp, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "failed to marshal the json for the ids", http.StatusInternalServerError)
			return
		}
		if _, err = w.Write(resp); err != nil {
			http.Error(w, "Failed to write response data", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		return
	}

	http.Error(w, "404 not found.", http.StatusNotFound)
}
