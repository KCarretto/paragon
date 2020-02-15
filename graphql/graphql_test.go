package graphql_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/graphql/models"
	"github.com/kcarretto/paragon/graphql/resolve"

	_ "github.com/mattn/go-sqlite3"
)

var (
	uniqueNumber int
)

// TODO (@rwhittier) need to change to require syntax
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
	jobCreater := test.Client.Job.Create().SetName("test").SetOwner(u).SetContent("wat")
	for _, opt := range options {
		opt(jobCreater)
	}
	job, err := jobCreater.Save(context.Background())
	if err != nil {
		t.Errorf("failed to create job: %w", err)
	}
	return job
}

func (test *testResolver) newCredential(t *testing.T, options ...func(*ent.CredentialCreate)) *ent.Credential {
	credentialCreater := test.Client.Credential.Create().SetPrincipal("testdata").SetSecret("testdata").SetKind("password").SetFails(1)
	for _, opt := range options {
		opt(credentialCreater)
	}
	credential, err := credentialCreater.Save(context.Background())
	if err != nil {
		t.Errorf("failed to create credential: %w", err)
	}
	return credential
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

func (test *testResolver) newLink(t *testing.T, options ...func(*ent.LinkCreate)) *ent.Link {
	uniqueNumber++
	uniqueData := fmt.Sprintf("test%d", uniqueNumber)
	linkCreater := test.Client.Link.Create().SetAlias(uniqueData)
	for _, opt := range options {
		opt(linkCreater)
	}
	link, err := linkCreater.Save(context.Background())
	if err != nil {
		t.Errorf("Failed to create link %w", err)
	}
	return link
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

func TestCredentialQuery(t *testing.T) {
	client := NewTestClient()
	defer client.Close()
	c := client.newCredential(t)
	query := client.Resolver.Query()
	queriedCredential, err := query.Credential(context.Background(), c.ID)
	if err != nil {
		t.Errorf("Credential query failed with %w", err)
	}
	if c.ID != queriedCredential.ID {
		t.Errorf("Credential query failed to get the correct target \n expected: %#v \n given: %#v", c, queriedCredential)
	}
}

func TestCredentialsQuery(t *testing.T) {
	client := NewTestClient()
	defer client.Close()
	c1 := client.newCredential(t)
	c2 := client.newCredential(t)
	query := client.Resolver.Query()
	queriedCredentials, err := query.Credentials(context.Background(), &models.Filter{})
	if err != nil {
		t.Errorf("Credentials query failed with %w", err)
	}
	if len(queriedCredentials) != 2 {
		t.Errorf("Credentials query returned wrong length expected: %d, got: %d", 2, len(queriedCredentials))
	}

	ids := map[int]int{c1.ID: 0, c2.ID: 0}
	for _, c := range queriedCredentials {
		if _, ok := ids[c.ID]; ok {
			delete(ids, c.ID)
		}
	}

	if len(ids) != 0 {
		idsLeft := []int{}
		for id := range ids {
			idsLeft = append(idsLeft, id)
		}
		t.Errorf("Credentials query returned missing expected credential(s) %#v", idsLeft)
	}

}

func TestLinkQuery(t *testing.T) {
	client := NewTestClient()
	defer client.Close()
	l := client.newLink(t)
	query := client.Resolver.Query()
	queriedLink, err := query.Link(context.Background(), l.ID)
	if err != nil {
		t.Errorf("Link query errored with %w", err)
	}
	if l.ID != queriedLink.ID {
		t.Errorf("Link query returned wrong id expected: %d, got: %d", l.ID, queriedLink.ID)
	}

}

func TestLinksQuery(t *testing.T) {
	client := NewTestClient()
	defer client.Close()
	numOfLinks := 2
	links := map[int]*ent.Link{}
	for i := 0; i < numOfLinks; i++ {
		l := client.newLink(t)
		links[l.ID] = l
	}
	query := client.Resolver.Query()
	queriedLinks, err := query.Links(context.Background(), &models.Filter{})
	if err != nil {
		t.Errorf("Links query errored with %w", err)
	}

	for _, l := range queriedLinks {
		if _, ok := links[l.ID]; ok {
			delete(links, l.ID)
		}
	}

	if len(links) != 0 {
		idsLeft := []int{}
		for id := range links {
			idsLeft = append(idsLeft, id)
		}
		t.Errorf("Links query returned missing expected links(s) %#v", idsLeft)
	}

}
