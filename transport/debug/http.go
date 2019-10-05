package debug

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/kcarretto/paragon/transport"
	"go.uber.org/zap"
)

func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

func (t *Transport) handleIndex(w http.ResponseWriter, req *http.Request) {

	w.Write([]byte(debugHTML))
	return
}

func (t *Transport) handleQueue(w http.ResponseWriter, req *http.Request) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Logger.Error("Failed to read request body", zap.Error(err))
		msg := fmt.Sprintf("failed to read request body: %s", err.Error())
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	fmt.Printf("\n\nTASK BYTES: %v\n\n", data)
	t.Logger.Debug("Request to queue task", zap.String("task", string(data)))
	var task transport.Task
	if err := json.Unmarshal(data, &task); err != nil {
		t.Logger.Error("Failed to parse request body", zap.Error(err))
		msg := fmt.Sprintf("failed to parse request body to json: %s", err.Error())
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	serverMsg, err := json.Marshal(
		transport.Payload{
			Tasks: []transport.Task{task},
		},
	)
	if err != nil {
		t.Logger.Error("Failed to marshal server message", zap.Error(err))
		msg := fmt.Sprintf("failed to marshal server message: %s", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	t.WritePayload(serverMsg)
}

func (t *Transport) handleResponses(w http.ResponseWriter, req *http.Request) {
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

	messages := []transport.Response{}
	for i := offset; i < len(t.messages) && i <= offset+limit; i++ {
		messages = append(messages, t.messages[i])
	}

	t.Logger.Debug(
		"Retrieving responses",
		zap.Int("offset", offset),
		zap.Int("limit", limit),
		zap.Int("sent_messages", len(messages)),
		zap.Int("total_messages", len(t.messages)),
	)

	respData, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal responses to json: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(respData); err != nil {
		t.Logger.Error("Failed to write response data to client", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to write response data to client: %s", err.Error()), http.StatusInternalServerError)
	}

}

// listenAndServe configures and runs an http server for agent debugging.
func (t *Transport) listenAndServe() {
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

	router.HandleFunc("/", use(t.handleIndex, middleware...))
	router.HandleFunc("/queue", use(t.handleQueue, middleware...))
	router.HandleFunc("/messages", use(t.handleResponses, middleware...))
	http.Handle("/", router)

	httpAddr := os.Getenv("DEBUG_API_ADDR")
	if httpAddr == "" {
		httpAddr = "127.0.0.1:8080"
	}

	t.Logger.Info("HTTP DEBUG ENABLED", zap.String("http_addr", httpAddr))
	http.ListenAndServe(httpAddr, router)
}
