package graphql_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/target"
	"github.com/kcarretto/paragon/graphql/models"
	"github.com/kcarretto/paragon/graphql/resolve"
	"github.com/kcarretto/paragon/pkg/auth"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
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
	targetCreater := test.Client.Target.Create().SetHostname("test").SetName("test").SetOS(target.OSLINUX).SetPrimaryIP("test")
	for _, opt := range options {
		opt(targetCreater)
	}

	target, err := targetCreater.Save(context.Background())
	require.NoError(t, err, "failed to create target")

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
	require.NoError(t, err, "failed to create user")

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
	require.NoError(t, err, "failed to create job")

	return job
}

func (test *testResolver) newCredential(t *testing.T, options ...func(*ent.CredentialCreate)) *ent.Credential {
	credentialCreater := test.Client.Credential.Create().SetPrincipal("testdata").SetSecret("testdata").SetKind("password").SetFails(1)
	for _, opt := range options {
		opt(credentialCreater)
	}
	credential, err := credentialCreater.Save(context.Background())
	require.NoError(t, err, "failed to create credential")

	return credential
}

func (test *testResolver) newTask(t *testing.T, options ...func(*ent.TaskCreate)) *ent.Task {
	job := test.newJob(t)
	taskCreater := test.Client.Task.Create().SetContent("test").SetLastChangedTime(time.Now()).SetJob(job)
	for _, opt := range options {
		opt(taskCreater)
	}

	task, err := taskCreater.Save(context.Background())
	require.NoError(t, err, "failed to create task")

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
	require.NoError(t, err, "Failed to create link")

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
	require.NoError(t, err, "Failed to create file")

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
	require.NoError(t, err, "failed to create tag")

	return tag
}

func (test *testResolver) newService(t *testing.T, options ...func(*ent.ServiceCreate)) *ent.Service {
	tag := test.newTag(t)
	uniqueNumber++
	uniqueData := fmt.Sprintf("test%d", uniqueNumber)
	serviceCreater := test.Client.Service.Create().SetName("test").SetPubKey(uniqueData).SetTag(tag)

	for _, opt := range options {
		opt(serviceCreater)
	}

	service, err := serviceCreater.Save(context.Background())
	require.NoError(t, err, "failed to create service")

	return service
}

func (test *testResolver) newEvent(t *testing.T, options ...func(*ent.EventCreate)) *ent.Event {
	eventCreater := test.Client.Event.Create().SetKind("OTHER")

	for _, opt := range options {
		opt(eventCreater)
	}

	service, err := eventCreater.Save(context.Background())
	require.NoError(t, err, "failed to create event")

	return service
}

func NewTestClient(t *testing.T) *testResolver {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	require.NoError(t, err, "failed opening connection to sqlite")

	err = client.Schema.Create(context.Background())
	require.NoError(t, err, "failed creating schema resources")

	return &testResolver{Client: client, Resolver: &resolve.Resolver{Graph: client}}
}

// test @master:/graphql/generated/generated.go -> MutationResolver and QueryResolver are the interfaces

func TestTargetsQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()

	client.newTarget(t, func(tar *ent.TargetCreate) { tar.SetName("test").SetOS(target.OSLINUX) })
	client.newTarget(t, func(tar *ent.TargetCreate) { tar.SetName("test2").SetOS(target.OSLINUX) })

	query := client.Resolver.Query()
	targets, err := query.Targets(context.Background(), &models.Filter{})
	require.NoError(t, err, "Targets query failed")
	require.Len(t, targets, 2, "Targets query failed to get the targets")

	t1, t2 := false, false
	for _, t := range targets {
		if t.Name == "test" {
			t1 = true
		} else if t.Name == "test2" {
			t2 = true
		}
	}

	require.True(t, t1, "Did not get correct name for first target")
	require.True(t, t2, "Did not get correct name for second target")
}

func TestTargetQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()

	target := client.newTarget(t)

	query := client.Resolver.Query()
	queriedTarget, err := query.Target(context.Background(), target.ID)
	require.NoError(t, err, "Target query failed with")
	require.Equal(t, target.ID, queriedTarget.ID)
}

func TestTaskQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()

	task := client.newTask(t)

	query := client.Resolver.Query()
	queriedTask, err := query.Task(context.Background(), task.ID)
	require.NoError(t, err, "Task query failed with")
	require.Equal(t, task.ID, queriedTask.ID)
}

func TestTasksQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()

	client.newTask(t)

	query := client.Resolver.Query()
	queriedTasks, err := query.Tasks(context.Background(), &models.Filter{})
	require.NoError(t, err, "Tasks query failed with")
	require.Len(t, queriedTasks, 1)
}

func TestCredentialQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()
	c := client.newCredential(t)
	query := client.Resolver.Query()
	queriedCredential, err := query.Credential(context.Background(), c.ID)
	require.NoError(t, err, "Credential query failed with")
	require.Equal(t, c.ID, queriedCredential.ID)
}

func TestCredentialsQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()
	c1 := client.newCredential(t)
	c2 := client.newCredential(t)
	query := client.Resolver.Query()
	queriedCredentials, err := query.Credentials(context.Background(), &models.Filter{})
	require.NoError(t, err, "Credentials query failed with")
	require.Len(t, queriedCredentials, 2)

	ids := map[int]int{c1.ID: 0, c2.ID: 0}
	for _, c := range queriedCredentials {
		if _, ok := ids[c.ID]; ok {
			delete(ids, c.ID)
		}
	}

	require.Len(t, ids, 0)
}

func TestLinkQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()

	l := client.newLink(t)

	query := client.Resolver.Query()
	queriedLink, err := query.Link(context.Background(), l.ID)
	require.NoError(t, err)
	require.Equal(t, l.ID, queriedLink.ID)
}

func TestLinksQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()

	numOfLinks := 2
	links := map[int]*ent.Link{}
	for i := 0; i < numOfLinks; i++ {
		l := client.newLink(t)
		links[l.ID] = l
	}

	query := client.Resolver.Query()
	queriedLinks, err := query.Links(context.Background(), &models.Filter{})
	require.NoError(t, err)

	for _, l := range queriedLinks {
		if _, ok := links[l.ID]; ok {
			delete(links, l.ID)
		}
	}
	require.Len(t, links, 0)
}

func TestFileQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()

	f := client.newFile(t)

	query := client.Resolver.Query()
	queriedFile, err := query.Link(context.Background(), f.ID)
	require.NoError(t, err)
	require.Equal(t, f.ID, queriedFile.ID)
}

func TestFilesQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()

	numOfFiles := 2
	files := map[int]*ent.File{}
	for i := 0; i < numOfFiles; i++ {
		f := client.newFile(t)
		files[f.ID] = f
	}

	query := client.Resolver.Query()
	queriedFiles, err := query.Links(context.Background(), &models.Filter{})
	require.NoError(t, err)

	for _, f := range queriedFiles {
		if _, ok := files[f.ID]; ok {
			delete(files, f.ID)
		}
	}
	require.Len(t, files, 0)
}

func TestJobQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()
	j := client.newJob(t)
	query := client.Resolver.Query()
	queriedJob, err := query.Job(context.Background(), j.ID)
	require.NoError(t, err)
	require.Equal(t, j.ID, queriedJob.ID)
}

func TestJobsQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()

	numOfJobs := 2
	jobs := map[int]*ent.Job{}
	for i := 0; i < numOfJobs; i++ {
		j := client.newJob(t)
		jobs[j.ID] = j
	}

	query := client.Resolver.Query()
	queriedJobs, err := query.Jobs(context.Background(), &models.Filter{})
	require.NoError(t, err)

	for _, j := range queriedJobs {
		if _, ok := jobs[j.ID]; ok {
			delete(jobs, j.ID)
		}
	}

	require.Len(t, jobs, 0)
}

func TestTagQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()
	tag := client.newTag(t)
	query := client.Resolver.Query()
	queriedTag, err := query.Tag(context.Background(), tag.ID)
	require.NoError(t, err)
	require.Equal(t, tag.ID, queriedTag.ID)
}

func TestTagsQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()
	numOfTags := 2
	tags := map[int]*ent.Tag{}
	for i := 0; i < numOfTags; i++ {
		tag := client.newTag(t)
		tags[tag.ID] = tag
	}
	query := client.Resolver.Query()
	queriedTags, err := query.Tags(context.Background(), &models.Filter{})
	require.NoError(t, err)

	for _, tag := range queriedTags {
		if _, ok := tags[tag.ID]; ok {
			delete(tags, tag.ID)
		}
	}

	require.Len(t, tags, 0)
}

func TestUserQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()
	u := client.newUser(t)
	query := client.Resolver.Query()
	queriedUser, err := query.User(context.Background(), u.ID)
	require.NoError(t, err)
	require.Equal(t, u.ID, queriedUser.ID)
}

func TestUsersQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()
	numOfUsers := 2
	users := map[int]*ent.User{}
	for i := 0; i < numOfUsers; i++ {
		u := client.newUser(t)
		users[u.ID] = u
	}
	query := client.Resolver.Query()
	queriedUsers, err := query.Users(context.Background(), &models.Filter{})
	require.NoError(t, err)

	for _, u := range queriedUsers {
		if _, ok := users[u.ID]; ok {
			delete(users, u.ID)
		}
	}

	require.Len(t, users, 0)
}

func TestServiceQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()
	s := client.newService(t)
	query := client.Resolver.Query()
	queriedService, err := query.Service(context.Background(), s.ID)
	require.NoError(t, err)
	require.Equal(t, s.ID, queriedService.ID)
}

func TestServicesQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()
	numOfServices := 2
	services := map[int]*ent.Service{}
	for i := 0; i < numOfServices; i++ {
		s := client.newService(t)
		services[s.ID] = s
	}
	query := client.Resolver.Query()
	queriedServices, err := query.Services(context.Background(), &models.Filter{})
	require.NoError(t, err)

	for _, s := range queriedServices {
		if _, ok := services[s.ID]; ok {
			delete(services, s.ID)
		}
	}

	require.Len(t, services, 0)
}

func TestEventQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()

	e := client.newEvent(t)

	query := client.Resolver.Query()
	queriedEvent, err := query.Event(context.Background(), e.ID)
	require.NoError(t, err)
	require.Equal(t, e.ID, queriedEvent.ID)
}

func TestEventsQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()

	numOfEvents := 2
	events := map[int]*ent.Event{}
	for i := 0; i < numOfEvents; i++ {
		e := client.newEvent(t)
		events[e.ID] = e
	}

	query := client.Resolver.Query()
	queriedEvents, err := query.Events(context.Background(), &models.Filter{})
	require.NoError(t, err)

	for _, e := range queriedEvents {
		if _, ok := events[e.ID]; ok {
			delete(events, e.ID)
		}
	}

	require.Len(t, events, 0)
}

func TestMeQuery(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()
	u := client.newUser(t, func(user *ent.UserCreate) { user.SetName("testuser101") })
	query := client.Resolver.Query()

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/graphql", nil)
	authedReq := auth.CreateUserSession(rr, req, u)
	_, err := query.Me(authedReq.Context())
	require.NoError(t, err)
}

// mutations
func TestFailCredentialMutation(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()
	mutation := client.Resolver.Mutation()
	c := client.newCredential(t)
	c, err := mutation.FailCredential(context.Background(), &models.FailCredentialRequest{ID: c.ID})
	require.NoError(t, err)
}
