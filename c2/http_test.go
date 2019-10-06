package c2_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kcarretto/paragon/c2"
	"github.com/kcarretto/paragon/transport"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHTTP(t *testing.T) {
	expectedTask := transport.Task{
		ID: "ABC",
	}
	unexpectedTask := transport.Task{
		ID: "XYZ",
	}

	msg := transport.Response{}
	data, err := json.Marshal(msg)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(data))

	srv := &c2.Server{
		Logger: zap.NewNop(),
	}
	srv.QueueTask(expectedTask, func(transport.Metadata) bool { return true })
	srv.QueueTask(unexpectedTask, func(transport.Metadata) bool { return false })

	srv.ServeHTTP(w, req)

	body, err := ioutil.ReadAll(w.Result().Body)
	defer w.Result().Body.Close()
	require.NoError(t, err)

	var reply transport.Payload
	err = json.Unmarshal(body, &reply)
	require.NoError(t, err)
	require.Equal(t, 1, len(reply.Tasks))
	require.Equal(t, expectedTask, reply.Tasks[0])
}
