package c2_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kcarretto/paragon/pkg/c2"
	"github.com/kcarretto/paragon/pkg/c2/mocks"
	"github.com/kcarretto/paragon/ent"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestServeHTTP(t *testing.T) {
	// Prepare mock teamserver
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	teamserver := mocks.NewMockTeamserver(ctrl)
	teamserver.EXPECT().ClaimTasks(gomock.Any(), gomock.Any()).Return([]*ent.Task{
		{
			ID: 1234,
			Content: `print('another test!')`,
		},
	}, nil)

	// Initialize C2 server
	srv := c2.Server{
		teamserver,
	}

	// Prepare HTTP Request from agent
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"metadata":{"machineUUID":"ABCD"}}`))
	resp := httptest.NewRecorder()

	// Serve the http request
	srv.ServeHTTP(resp, req)

	// Ensure expected response
	body, err := ioutil.ReadAll(resp.Result().Body)
	defer resp.Result().Body.Close()
	require.NoError(t, err)
	require.Equal(t, "{\"tasks\":[{\"id\":1234,\"content\":\"print('another test!')\"}]}\n", string(body))
}