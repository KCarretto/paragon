package c2_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kcarretto/paragon/proto/codec"
	"github.com/kcarretto/paragon/pkg/c2"

	"github.com/stretchr/testify/require"
)

func TestHTTP(t *testing.T) {
	// Prepare HTTP Request from agent
	msg := codec.AgentMessage{}
	data, err := json.Marshal(msg)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(data))

	// Prepare server
	srv := &c2.Server{
		Queue: &c2.Queue{},
	}

	// Serve the http request
	srv.ServeHTTP(w, req)

	body, err := ioutil.ReadAll(w.Result().Body)
	defer w.Result().Body.Close()
	require.NoError(t, err)

	var reply codec.ServerMessage
	err = json.Unmarshal(body, &reply)
	require.NoError(t, err)
}
