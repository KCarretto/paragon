package graphql_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/graphql/models"
	"github.com/kcarretto/paragon/graphql/resolve"

	_ "github.com/mattn/go-sqlite3"
)

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

// test @master:/graphql/generated/generated.go -> MutationResolver and QueryResolver are the interfaces

func TestTargetsQuery(t *testing.T) {
	client := testClient()
	defer client.Close()
	_, err := client.Target.Create().SetHostname("test").SetName("test").SetPrimaryIP("test").Save(context.Background())
	_, err = client.Target.Create().SetHostname("test2").SetName("test2").SetPrimaryIP("test2").Save(context.Background())
	if err != nil {
		t.Errorf("failed to create target: %w", err)
	}
	resolver := &resolve.Resolver{Graph: client}
	query := resolver.Query()
	targets, err := query.Targets(context.Background(), &models.Filter{})
	if err != nil {
		t.Errorf("Targets query failed with %w", err)
	}
	if len(targets) != 2 {
		t.Errorf("Targets query failed to get the targets \n expected: %d \n given: %d", 2, len(targets))
	}
	t1, t2 := false, false
	for _, t := range targets {
		if t.Name == "test" {
			t1 = true
		} else if t.Name == "test2" {
			t2 = true
		}
	}
	if !(t1 && t2) {
		t.Errorf("Did not get the correct names for the targets")
	}
}

func TestTargetQuery(t *testing.T) {
	client := testClient()
	defer client.Close()
	target, err := client.Target.Create().SetHostname("test").SetName("test").SetPrimaryIP("test").Save(context.Background())
	if err != nil {
		t.Errorf("failed to create target: %w", err)
	}
	resolver := &resolve.Resolver{Graph: client}
	query := resolver.Query()
	queriedTarget, err := query.Target(context.Background(), target.ID)
	if err != nil {
		t.Errorf("Target query failed with %w", err)
	}
	if queriedTarget.ID != target.ID {
		t.Errorf("Target query failed to get the correct target \n expected: %#v \n given: %#v", target, queriedTarget)
	}
}

func TestJobQuery(t *testing.T) {
	client := testClient()
	defer client.Close()
	u, err := client.User.Create().SetName("joe").SetOAuthID("who").SetPhotoURL("uneccsarybutfine").Save(context.Background())
	if err != nil {
		t.Errorf("failed to create user: %w", err)
	}
	job, err := client.Job.Create().SetName("test").SetOwner(u).SetContent("wat").Save(context.Background())
	if err != nil {
		t.Errorf("failed to create job: %w", err)
	}
	task, err := client.Task.Create().SetContent("test").SetLastChangedTime(time.Now()).SetJob(job).Save(context.Background())
	if err != nil {
		t.Errorf("failed to create task: %w", err)
	}
	resolver := &resolve.Resolver{Graph: client}
	query := resolver.Query()
	queriedTask, err := query.Task(context.Background(), task.ID)
	if err != nil {
		t.Errorf("Task query failed with %w", err)
	}
	if queriedTask.ID != task.ID {
		t.Errorf("Task query failed to get the correct target \n expected: %#v \n given: %#v", task, queriedTask)
	}
}

func TestJobsQuery(t *testing.T) {
	client := testClient()
	defer client.Close()
	u, err := client.User.Create().SetName("joe").SetOAuthID("who").SetPhotoURL("uneccsarybutfine").Save(context.Background())
	if err != nil {
		t.Errorf("failed to create user: %w", err)
	}
	job, err := client.Job.Create().SetName("test").SetOwner(u).SetContent("wat").Save(context.Background())
	if err != nil {
		t.Errorf("failed to create job: %w", err)
	}
	_, err = client.Task.Create().SetContent("test").SetLastChangedTime(time.Now()).SetJob(job).Save(context.Background())
	if err != nil {
		t.Errorf("failed to create task: %w", err)
	}
	resolver := &resolve.Resolver{Graph: client}
	query := resolver.Query()
	queriedTasks, err := query.Tasks(context.Background(), &models.Filter{})
	if err != nil {
		t.Errorf("Tasks query failed with %w", err)
	}
	if len(queriedTasks) != 1 {
		t.Errorf("Tasks query failed to get the correct len \n expected: %#v \n given: %#v", 1, len(queriedTasks))
	}
}
