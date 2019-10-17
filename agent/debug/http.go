package debug

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/kcarretto/paragon/agent"
	"github.com/kcarretto/paragon/api/codec"
	"go.uber.org/zap"
)

func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

func (transport *Sender) handleIndex(w http.ResponseWriter, req *http.Request) {

	w.Write([]byte(debugHTML))
	return
}

// handleQueue handles an http request to queue a task.
func (transport *Sender) handleQueue(w http.ResponseWriter, req *http.Request) {
	// Read request data
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		transport.Log.Error("Failed to read request body", zap.Error(err))
		msg := fmt.Sprintf("failed to read request body: %s", err.Error())
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	transport.Log.Debug("Request to queue task", zap.String("task", string(data)))

	// Unmarshal request into task
	var task codec.Task
	if err := json.Unmarshal(data, &task); err != nil {
		transport.Log.Error("Failed to parse request body", zap.Error(err))
		msg := fmt.Sprintf("failed to parse request body to json: %s", err.Error())
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// Queue the task
	transport.QueueTask(&task)
}

// handleMessages handles an http request to view agent messages
func (transport *Sender) handleMessages(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()

	offsetStr := params.Get("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}

	limitStr := params.Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid offset provided %q: %s", offsetStr, err.Error()), http.StatusBadRequest)
		return
	}
	if offset < 0 {
		offset = 0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid limit provided %q: %s", limitStr, err.Error()), http.StatusBadRequest)
		return
	}
	if limit <= 0 {
		limit = 1
	}

	messages := []agent.Message{}
	for i := offset; i < len(transport.messages) && i <= offset+limit; i++ {
		messages = append(messages, transport.messages[i])
	}

	transport.Log.Debug(
		"Retrieving responses",
		zap.Int("offset", offset),
		zap.Int("limit", limit),
		zap.Int("sent_messages", len(messages)),
		zap.Int("total_messages", len(transport.messages)),
	)

	respData, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal responses to json: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(respData); err != nil {
		transport.Log.Error("Failed to write response data to client", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to write response data to client: %s", err.Error()), http.StatusInternalServerError)
	}
}

// listenAndServe configures and runs an http server for agent debugging.
func (transport *Sender) listenAndServe() {
	router := http.NewServeMux()
	var middleware []func(http.HandlerFunc) http.HandlerFunc

	// Prepare basic auth
	username := os.Getenv("DEBUG_API_USER")
	password := os.Getenv("DEBUG_API_PASSWORD")
	if username == "" {
		username = "admin"
	}
	if password != "" {
		middleware = append(middleware, prepareBasicAuth(username, password))
	}

	httpAddr := os.Getenv("DEBUG_API_ADDR")
	if httpAddr == "" {
		httpAddr = "127.0.0.1:8080"
	}

	router.HandleFunc("/", use(transport.handleIndex, middleware...))
	router.HandleFunc("/queue", use(transport.handleQueue, middleware...))
	router.HandleFunc("/messages", use(transport.handleMessages, middleware...))
	transport.srv = &http.Server{
		Addr:    httpAddr,
		Handler: router,
	}

	transport.Log.Info("HTTP DEBUG ENABLED", zap.String("http_addr", httpAddr))
	if err := transport.srv.ListenAndServe(); err != nil {
		transport.Log.Error("HTTP Debug Server encountered an error while closing", zap.Error(err))
	}
	transport.Log.Info("HTTP Debug Server Closed")
}
