package c2

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/kcarretto/paragon/transport"
	"go.uber.org/zap"
)

// ServeHTTP implements http.Handler to handle agent messages sent via http.
func (srv *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		srv.Logger.Error("Failed to read agent message", zap.Error(err))
		return
	}

	srv.Logger.Debug("Received agent message", zap.String("agent_msg", string(data)))

	var msg transport.Response
	if err := json.Unmarshal(data, &msg); err != nil {
		srv.Logger.Error("Failed to unmarshal agent message", zap.Error(err))
		return
	}

	if err = srv.HandleMessage(w, msg); err != nil {
		srv.Logger.Error("Agent communication failed", zap.Error(err))
		return
	}

	// TODO: Add agent metadata
	srv.Logger.Info("Replied to agent message")

}
