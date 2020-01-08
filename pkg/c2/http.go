package c2

import (
	"net/http"
)

// ServeHTTP wraps HandleJSON to provide a server handler for HTTP(s) transport.
func (srv Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Establish request context
	ctx := req.Context()

	// Wrap JSON Handler
	if err := srv.HandleJSON(ctx, req.Body, w); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}