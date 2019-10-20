package c2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kcarretto/paragon/api/codec"
)

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
}
