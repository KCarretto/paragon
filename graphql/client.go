package teamserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/graphql/models"
	"net/http"
)

// Response stores decoded JSON data returned from a GraphQL request.
type Response struct {
	Data   map[string]interface{}
	Errors []string
}

// Request stores GraphQL queries and mutations.
type Request struct {
	Query     string
	Variables map[string]interface{}
}

// A Client can be used to request GraphQL queries and mutations using HTTP.
type Client struct {
	URL  string
	HTTP *http.Client
}

// Do executes a GraphQL request.
func (client Client) Do(r Request) (*Response, error) {
	// Encode request payload
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to encode json: %w", err)
	}

	// Build http request
	req, err := http.NewRequest(http.MethodPost, client.URL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set default http client if necessary
	if client.HTTP == nil {
		client.HTTP = http.DefaultClient
	}

	// Issue the request
	resp, err := client.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Decode the response
	response := &Response{}
	data := json.NewDecoder(resp.Body)
	if err := data.Decode(response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}

func ClaimTasks(params models.ClaimTaskRequest) ([]*ent.Task, error) {
	// TODO: Create an new Request object with query string and variables from params
	// TODO: Execute request and parse response into an array of tasks, or else return an error

	return nil, nil
}
