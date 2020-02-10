package teamserver_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/pkg/teamserver"

	_ "github.com/mattn/go-sqlite3"
)

// An Error returned by the GraphQL server
type Error struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

func TestTeamserver(t *testing.T) {
	client := testClient()
	defer client.Close()
	svc := &teamserver.Service{
		Graph: client,
	}
	router := http.NewServeMux()
	svc.HTTP(router)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/graphql", nil)

	router.ServeHTTP(rr, req)

	// fmt.Printf("%s", rr.Body.String())
	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}

}

func testClient() *ent.Client {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}

func TestTargetsReturnsCorrectly(t *testing.T) {
	client := testClient()
	defer client.Close()
	_, err := client.Target.Create().SetHostname("test").SetName("test").SetPrimaryIP("test").Save(context.Background())
	if err != nil {
		t.Errorf("failed to create target: %w", err)
	}
	srv := &teamserver.Server{EntClient: client}
	router := teamserver.NewRouter(srv)

	rr := httptest.NewRecorder()
	query := "{\"query\": \"{ targets { id } }\"}"
	req, err := http.NewRequest("POST", "/graphql", bytes.NewBufferString(query))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(rr, req)

	// fmt.Printf("%s", rr.Body.String())
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}
	// Prepare response
	var resp struct {
		Data struct {
			Targets []*ent.Target `json:"targets"`
		} `json:"data"`
		Errors []Error `json:"errors"`
	}
	data := json.NewDecoder(rr.Body)
	if err := data.Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %w", err)
	}
	if len(resp.Data.Targets) != 1 {
		t.Errorf("Targets had more targets than were expected returned")
	}
}

// // test @master:/graphql/generated/generated.go -> MutationResolver and QueryResolver are the interfaces

// func TestTargetQuery(t *testing.T) {
// 	client := testClient()
// 	defer client.Close()
// 	_, err := client.Target.Create().SetHostname("test").SetName("test").SetPrimaryIP("test").Save(context.Background())
// 	_, err = client.Target.Create().SetHostname("test2").SetName("test2").SetPrimaryIP("test2").Save(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to create target: %w", err)
// 	}
// 	resolver := &resolve.Resolver{EntClient: client}
// 	query := resolver.Query()
// 	targets, err := query.Targets(context.Background())
// 	if err != nil {
// 		t.Errorf("Targets query failed with %w", err)
// 	}
// 	if len(targets) != 2 {
// 		t.Errorf("Targets query failed to get the targets \n expected: %d \n given: %d", 2, len(targets))
// 	}
// 	t1, t2 := false, false
// 	for _, t := range targets {
// 		if t.Name == "test" {
// 			t1 = true
// 		} else if t.Name == "test2" {
// 			t2 = true
// 		}
// 	}
// 	if !(t1 && t2) {
// 		t.Errorf("Did not get the correct names for the targets")
// 	}
// }

// func TestTargetsQuery(t *testing.T) {
// 	client := testClient()
// 	defer client.Close()
// 	target, err := client.Target.Create().SetHostname("test").SetName("test").SetPrimaryIP("test").Save(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to create target: %w", err)
// 	}
// 	resolver := &resolve.Resolver{EntClient: client}
// 	query := resolver.Query()
// 	queriedTarget, err := query.Target(context.Background(), target.ID)
// 	if err != nil {
// 		t.Errorf("Target query failed with %w", err)
// 	}
// 	if queriedTarget.ID != target.ID {
// 		t.Errorf("Target query failed to get the correct target \n expected: %#v \n given: %#v", target, queriedTarget)
// 	}

// }
// func TestJobQuery(t *testing.T) {
// 	client := testClient()
// 	defer client.Close()
// 	job, err := client.Job.Create().SetName("test").SetContent("wat").Save(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to create job: %w", err)
// 	}
// 	task, err := client.Task.Create().SetContent("test").SetJob(job).Save(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to create target: %w", err)
// 	}
// 	resolver := &resolve.Resolver{EntClient: client}
// 	query := resolver.Query()
// 	queriedTask, err := query.Task(context.Background(), task.ID)
// 	if err != nil {
// 		t.Errorf("Task query failed with %w", err)
// 	}
// 	if queriedTask.ID != task.ID {
// 		t.Errorf("Task query failed to get the correct target \n expected: %#v \n given: %#v", task, queriedTask)
// 	}

// }

// func TestJobsQuery(t *testing.T) {
// 	client := testClient()
// 	defer client.Close()
// 	job, err := client.Job.Create().SetName("test").SetContent("wat").Save(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to create job: %w", err)
// 	}
// 	task, err := client.Task.Create().SetContent("test").SetJob(job).Save(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to create target: %w", err)
// 	}
// 	resolver := &resolve.Resolver{EntClient: client}
// 	query := resolver.Query()
// 	queriedTask, err := query.Task(context.Background(), task.ID)
// 	if err != nil {
// 		t.Errorf("Task query failed with %w", err)
// 	}
// 	if queriedTask.ID != task.ID {
// 		t.Errorf("Task query failed to get the correct target \n expected: %#v \n given: %#v", task, queriedTask)
// 	}

// }
