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

type testResolver struct {
	*ent.Client
	Resolver *resolve.Resolver
}

// potentially useful, but a lot of work by hand to set up so why?
// type testTarget struct {
// 	*ent.TargetClient
// }

// func (t *testTarget) Create() *ent.TargetCreate {
// 	return t.TargetClient.Create().SetHostname("test").SetName("test").SetPrimaryIP("test")
// }

func newTarget(client *testResolver) *ent.TargetCreate {
	return client.Target.Create().SetHostname("test").SetName("test").SetPrimaryIP("test")
}

func NewTestClient() *testResolver {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return &testResolver{Client: client, Resolver: &resolve.Resolver{Graph: client}}
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

// test @master:/graphql/generated/generated.go -> MutationResolver and QueryResolver are the interfaces

func TestTargetsQuery(t *testing.T) {
	client := NewTestClient()
	defer client.Close()
	createOrFail := func(target *ent.TargetCreate) {
		_, err := target.Save(context.Background())
		if err != nil {
			t.Errorf("failed to create target: %w", err)
		}
	}
	createOrFail(newTarget(client).SetName("test"))
	createOrFail(newTarget(client).SetName("test2"))

	query := client.Resolver.Query()
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
	client := NewTestClient()
	defer client.Close()
	target, err := client.Target.Create().SetHostname("test").SetName("test").SetPrimaryIP("test").Save(context.Background())
	if err != nil {
		t.Errorf("failed to create target: %w", err)
	}
	query := client.Resolver.Query()
	queriedTarget, err := query.Target(context.Background(), target.ID)
	if err != nil {
		t.Errorf("Target query failed with %w", err)
	}
	if queriedTarget.ID != target.ID {
		t.Errorf("Target query failed to get the correct target \n expected: %#v \n given: %#v", target, queriedTarget)
	}
}

func TestJobQuery(t *testing.T) {
	client := NewTestClient()
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
	query := client.Resolver.Query()
	queriedTask, err := query.Task(context.Background(), task.ID)
	if err != nil {
		t.Errorf("Task query failed with %w", err)
	}
	if queriedTask.ID != task.ID {
		t.Errorf("Task query failed to get the correct target \n expected: %#v \n given: %#v", task, queriedTask)
	}
}

func TestJobsQuery(t *testing.T) {
	client := NewTestClient()
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
	query := client.Resolver.Query()
	queriedTasks, err := query.Tasks(context.Background(), &models.Filter{})
	if err != nil {
		t.Errorf("Tasks query failed with %w", err)
	}
	if len(queriedTasks) != 1 {
		t.Errorf("Tasks query failed to get the correct len \n expected: %#v \n given: %#v", 1, len(queriedTasks))
	}
}
