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

func (test *testResolver) newTarget(t *testing.T, options ...func(*ent.TargetCreate)) *ent.Target {
	targetCreater := test.Client.Target.Create().SetHostname("test").SetName("test").SetPrimaryIP("test")
	for _, opt := range options {
		opt(targetCreater)
	}
	target, err := targetCreater.Save(context.Background())
	if err != nil {
		t.Errorf("failed to create target: %w", err)
	}
	return target
}

func (test *testResolver) newUser(t *testing.T, options ...func(*ent.UserCreate)) *ent.User {
	userCreater := test.Client.User.Create().SetName("joe").SetOAuthID("who").SetPhotoURL("uneccsarybutfine")
	for _, opt := range options {
		opt(userCreater)
	}
	user, err := userCreater.Save(context.Background())
	if err != nil {
		t.Errorf("failed to create user: %w", err)
	}
	return user
}

func (test *testResolver) newJob(t *testing.T, options ...func(*ent.JobCreate)) *ent.Job {
	u := test.newUser(t)
	jobCreater := test.Client.Job.Create().
		SetName("test").
		SetOwner(u).
		SetContent("wat").
		SetStaged(false)

	for _, opt := range options {
		opt(jobCreater)
	}
	job, err := jobCreater.Save(context.Background())
	if err != nil {
		t.Errorf("failed to create job: %w", err)
	}
	return job
}

func (test *testResolver) newTask(t *testing.T, options ...func(*ent.TaskCreate)) *ent.Task {
	job := test.newJob(t)
	taskCreater := test.Client.Task.Create().SetContent("test").SetLastChangedTime(time.Now()).SetJob(job)
	for _, opt := range options {
		opt(taskCreater)
	}
	task, err := taskCreater.Save(context.Background())
	if err != nil {
		t.Errorf("failed to create task: %w", err)
	}
	return task
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
	client.newTarget(t, func(target *ent.TargetCreate) { target.SetName("test") })
	client.newTarget(t, func(target *ent.TargetCreate) { target.SetName("test2") })

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
	target := client.newTarget(t)
	query := client.Resolver.Query()
	queriedTarget, err := query.Target(context.Background(), target.ID)
	if err != nil {
		t.Errorf("Target query failed with %w", err)
	}
	if queriedTarget.ID != target.ID {
		t.Errorf("Target query failed to get the correct target \n expected: %#v \n given: %#v", target, queriedTarget)
	}
}

func TestTaskQuery(t *testing.T) {
	client := NewTestClient()
	defer client.Close()
	task := client.newTask(t)
	query := client.Resolver.Query()
	queriedTask, err := query.Task(context.Background(), task.ID)
	if err != nil {
		t.Errorf("Task query failed with %w", err)
	}
	if queriedTask.ID != task.ID {
		t.Errorf("Task query failed to get the correct target \n expected: %#v \n given: %#v", task, queriedTask)
	}
}

func TestTasksQuery(t *testing.T) {
	client := NewTestClient()
	defer client.Close()
	client.newTask(t)
	query := client.Resolver.Query()
	queriedTasks, err := query.Tasks(context.Background(), &models.Filter{})
	if err != nil {
		t.Errorf("Tasks query failed with %w", err)
	}
	if len(queriedTasks) != 1 {
		t.Errorf("Tasks query failed to get the correct len \n expected: %#v \n given: %#v", 1, len(queriedTasks))
	}
}
