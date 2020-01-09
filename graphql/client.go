package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/graphql/models"
)

// An Error returned by the GraphQL server
type Error struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

// Error implements the error interface by formatting an error message of the available error info.
func (err Error) Error() string {
	return fmt.Sprintf("%s (path: %s)", err.Message, strings.Join(err.Path, ", "))
}

// A Request stores execution properties of GraphQL queries and mutations.
type Request struct {
	Operation string      `json:"operationName"`
	Query     string      `json:"query"`
	Variables interface{} `json:"variables"`
}

// A Client can be used to request GraphQL queries and mutations using HTTP.
type Client struct {
	URL  string
	HTTP *http.Client
}

// Do executes a GraphQL request and unmarshals the JSON result into the destination struct.
func (client Client) Do(ctx context.Context, request Request, dst interface{}) error {
	// Encode request payload
	payload, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}

	// Build http request
	httpReq, err := http.NewRequest(http.MethodPost, client.URL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Set default http client if necessary
	if client.HTTP == nil {
		client.HTTP = http.DefaultClient
	}

	// Issue the request
	httpResp, err := client.HTTP.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	// Check response status
	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status error: %s", httpResp.Status)
	}

	tstData, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received response from teamserver: %s\n", string(tstData))

	// Decode the response
	// data := json.NewDecoder(httpResp.Body)
	data := json.NewDecoder(bytes.NewBuffer(tstData))
	if err := data.Decode(dst); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// ClaimTasks for a target that has the provided attributes, returning an array of tasks to execute.
// If no tasks are available, an empty task array is returned. If no target can be found, an error
// will be returned.
func (client Client) ClaimTasks(ctx context.Context, vars models.ClaimTaskRequest) ([]*ent.Task, error) {
	// Build request
	req := Request{
		Operation: "ClaimTasks",
		Query: `
		mutation ClaimTasks($params: ClaimTaskRequest!) {
			claimTasks(input: $params) {
			  id
			  content
			}
		}`,
		Variables: map[string]interface{}{
			"params": vars,
		},
	}

	// Prepare response
	var resp struct {
		Data struct {
			Tasks []*ent.Task `json:"claimTasks"`
		} `json:"data"`
		Errors []Error `json:"errors"`
	}

	// Execute mutation
	if err := client.Do(ctx, req, &resp); err != nil {
		return nil, err
	}

	fmt.Printf("Response from teamserver: %+v\n", resp)

	// Check for errors
	if resp.Errors != nil {
		return nil, fmt.Errorf("mutation failed: [%+v]", resp.Errors)
	}

	// Return claimed tasks
	return resp.Data.Tasks, nil
}

// SubmitTaskResult updates a task with execution output.
func (client Client) SubmitTaskResult(ctx context.Context, vars models.SubmitTaskResultRequest) error {
	// Build request
	req := Request{
		Operation: "SubmitTaskResult",
		Query: `
		mutation SubmitTaskResult($params: SubmitTaskResultRequest!) {
			submitTaskResult(input: $params) {
			  id
			}
		}`,
		Variables: map[string]interface{}{
			"params": vars,
		},
	}

	// Prepare response
	var resp struct {
		Errors []Error `json:"errors"`
	}

	// Execute mutation
	if err := client.Do(ctx, req, &resp); err != nil {
		return err
	}

	// Check for errors
	if resp.Errors != nil {
		return fmt.Errorf("mutation failed: [%+v]", resp.Errors)
	}

	return nil
}
