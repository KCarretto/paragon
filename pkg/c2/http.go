package c2

import (
	"fmt"
	"io"
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

// ServeStatus is an HTTP handler that returns the C2 Server status
func (srv *Server) ServeStatus(w http.ResponseWriter, req *http.Request) {
	if _, err := io.WriteString(w, `{"status":"OK"}`); err != nil {
		http.Error(
			w,
			fmt.Sprintf("failed to write response data: %s", err.Error()),
			http.StatusInternalServerError,
		)
	}
	w.Header().Set("Content-Type", "application/json")
}
