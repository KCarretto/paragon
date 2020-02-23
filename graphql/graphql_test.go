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
	uniqueNumber++
	uniqueData := fmt.Sprintf("test%d", uniqueNumber)
	userCreater := test.Client.User.Create().SetName("joe").SetOAuthID(uniqueData).SetPhotoURL("uneccsarybutfine")
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

func (test *testResolver) newFile(t *testing.T, options ...func(*ent.FileCreate)) *ent.File {
	l := test.newLink(t)
	uniqueNumber++
	uniqueData := fmt.Sprintf("test%d", uniqueNumber)
	fileCreater := test.Client.File.Create().SetName(uniqueData).SetLastModifiedTime(time.Now()).SetContent([]byte("test")).SetHash("test").SetContentType("test").AddLinks(l)
	for _, opt := range options {
		opt(fileCreater)
	}
	file, err := fileCreater.Save(context.Background())
	if err != nil {
		t.Errorf("Failed to create file %w", err)
	}
	return file
}

func (test *testResolver) newTag(t *testing.T, options ...func(*ent.TagCreate)) *ent.Tag {
	uniqueNumber++
	uniqueData := fmt.Sprintf("test%d", uniqueNumber)
	tagCreater := test.Client.Tag.Create().SetName(uniqueData)
	for _, opt := range options {
		opt(tagCreater)
	}
	tag, err := tagCreater.Save(context.Background())
	if err != nil {
		t.Errorf("failed to create tag: %w", err)
	}
	return tag
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

func TestFileQuery(t *testing.T) {
	client := NewTestClient()
	defer client.Close()
	f := client.newFile(t)
	query := client.Resolver.Query()
	queriedFile, err := query.Link(context.Background(), f.ID)
	if err != nil {
		t.Errorf("Link query errored with %w", err)
	}
	if f.ID != queriedFile.ID {
		t.Errorf("Link query returned wrong id expected: %d, got: %d", f.ID, queriedFile.ID)
	}

}

func TestFilesQuery(t *testing.T) {
	client := NewTestClient()
	defer client.Close()
	numOfFiles := 2
	files := map[int]*ent.File{}
	for i := 0; i < numOfFiles; i++ {
		f := client.newFile(t)
		files[f.ID] = f
	}
	query := client.Resolver.Query()
	queriedFiles, err := query.Links(context.Background(), &models.Filter{})
	if err != nil {
		t.Errorf("Links query errored with %w", err)
	}

	for _, f := range queriedFiles {
		if _, ok := files[f.ID]; ok {
			delete(files, f.ID)
		}
	}

	if len(files) != 0 {
		idsLeft := []int{}
		for id := range files {
			idsLeft = append(idsLeft, id)
		}
		t.Errorf("Files query returned missing expected file(s) %#v", idsLeft)
	}

}

func TestJobQuery(t *testing.T) {
	client := NewTestClient()
	defer client.Close()
	j := client.newJob(t)
	query := client.Resolver.Query()
	queriedJob, err := query.Job(context.Background(), j.ID)
	if err != nil {
		t.Errorf("Job query errored with %w", err)
	}
	if j.ID != queriedJob.ID {
		t.Errorf("Job query returned wrong id expected: %d, got: %d", j.ID, queriedJob.ID)
	}

}

func TestJobsQuery(t *testing.T) {
	client := NewTestClient()
	defer client.Close()
	numOfJobs := 2
	jobs := map[int]*ent.Job{}
	for i := 0; i < numOfJobs; i++ {
		j := client.newJob(t)
		jobs[j.ID] = j
	}
	query := client.Resolver.Query()
	queriedJobs, err := query.Jobs(context.Background(), &models.Filter{})
	if err != nil {
		t.Errorf("Jobs query errored with %w", err)
	}

	for _, j := range queriedJobs {
		if _, ok := jobs[j.ID]; ok {
			delete(jobs, j.ID)
		}
	}

	if len(jobs) != 0 {
		idsLeft := []int{}
		for id := range jobs {
			idsLeft = append(idsLeft, id)
		}
		t.Errorf("Jobs query returned missing expected job(s) %#v", idsLeft)
	}

}

func TestTagQuery(t *testing.T) {
	client := NewTestClient()
	defer client.Close()
	tag := client.newTag(t)
	query := client.Resolver.Query()
	queriedTag, err := query.Tag(context.Background(), tag.ID)
	if err != nil {
		t.Errorf("Tag query errored with %w", err)
	}
	if tag.ID != queriedTag.ID {
		t.Errorf("Tag query returned wrong id expected: %d, got: %d", tag.ID, queriedTag.ID)
	}
}

func TestTagsQuery(t *testing.T) {
	client := NewTestClient()
	defer client.Close()
	numOfTags := 2
	tags := map[int]*ent.Tag{}
	for i := 0; i < numOfTags; i++ {
		tag := client.newTag(t)
		tags[tag.ID] = tag
	}
	query := client.Resolver.Query()
	queriedTags, err := query.Tags(context.Background(), &models.Filter{})
	if err != nil {
		t.Errorf("Jobs query errored with %w", err)
	}

	for _, tag := range queriedTags {
		if _, ok := tags[tag.ID]; ok {
			delete(tags, tag.ID)
		}
	}

	if len(tags) != 0 {
		idsLeft := []int{}
		for id := range tags {
			idsLeft = append(idsLeft, id)
		}
		t.Errorf("Tags query returned missing expected tag(s) %#v", idsLeft)
	}

}
