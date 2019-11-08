package http_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/kcarretto/paragon/pkg/agent"
	transport "github.com/kcarretto/paragon/pkg/agent/http"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	// Prepare expected messages
	expectedAgentMsg := agent.Message{}
	expectedServerMsg := agent.ServerMessage{}

	// Initialize new test server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Read body
		body, err := ioutil.ReadAll(req.Body)
		require.NoError(t, err)

		// Unmarshal message
		var msg agent.Message
		err = json.Unmarshal(body, &msg)
		require.NoError(t, err)
		require.Equal(t, expectedAgentMsg, msg)

		// Marshal & write response
		resp, err := json.Marshal(expectedServerMsg)
		require.NoError(t, err)
		w.Write(resp)
	}))
	defer srv.Close()

	// Prepare URL
	srvURL, err := url.Parse(srv.URL)
	require.NoError(t, err)

	// Prepare Sender
	sender := transport.Sender{
		URL: srvURL,
	}

	// Send message & verify response
	var msg agent.ServerMessage
	err = sender.Send(&msg, expectedAgentMsg)
	require.NoError(t, err)
	require.Equal(t, expectedServerMsg, msg)
}
