package jobs

import (
	"context"
	"errors"
	"time"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/job"
	"github.com/kcarretto/paragon/ent/predicate"
	"github.com/kcarretto/paragon/ent/tag"
	"github.com/kcarretto/paragon/ent/target"
)

// Service is the generic struct used to handle the Ent service
type Service struct {
	EntClient *ent.Client
}

func castIntArrayToInt64Array(ints []int) []int64 {
	var temp []int64
	for _, val := range ints {
		temp = append(temp, int64(val))
	}
	return temp
}

// Fetch is used for getting a single Ent given an ID
// TODO: @cictrone
func (svc *Service) Fetch(ctx context.Context, req *FetchRequest) (*FetchResponse, error) {
	id := req.GetId()
	if id == 0 {
		return nil, errors.New("expected id but none given")
	}
	job := svc.EntClient.Job.GetX(ctx, int(id))
	return &FetchResponse{
		Name:    job.Name,
		Content: job.Content,
		Tasks:   castIntArrayToInt64Array(job.QueryTasks().IDsX(ctx)),
		Tags:    castIntArrayToInt64Array(job.QueryTags().IDsX(ctx)),
	}, nil
}

// Find is used for getting a list of Ent IDs given a string
// TODO: @cictrone
func (svc *Service) Find(ctx context.Context, req *FindRequest) (*FindResponse, error) {
	offset := req.GetOffset()
	limit := req.GetLimit()
	if limit == 0 {
		limit = 10
	}

	ids := svc.EntClient.Job.Query().
		Where(job.Or(
			job.NameContains(req.GetFilter()),
			job.ContentContains(req.GetFilter()),
			job.HasTagsWith(tag.NameContains(req.GetFilter())),
		)).
		Offset(int(offset)).
		Limit(int(limit)).
		IDsX(ctx)

	newOffset := int(offset) + len(ids)
	return &FindResponse{Ids: castIntArrayToInt64Array(ids), NewOffset: int64(newOffset)}, nil
}

// Create is used for creating an Ent given proper parameters
// TODO: @cictrone
func (svc *Service) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	name := req.GetName()
	if name == "" {
		return nil, errors.New("expected name but none given")
	}

	content := req.GetContent()
	if content == "" {
		return nil, errors.New("expected content but none given")
	}

	tags := req.GetTags()
	var tagEnts []*ent.Tag
	for _, tagName := range tags {
		tagEnts = append(tagEnts, svc.EntClient.Tag.Query().
			Where(tag.NameEQ(tagName)).
			OnlyX(ctx))
	}
	job := svc.EntClient.Job.Create().
		SetName(name).
		SetContent(content).
		AddTags(tagEnts...).
		SaveX(ctx)
	return &CreateResponse{Id: int64(job.ID)}, nil
}

// Queue is used for enqueueing a job
// TODO: @cictrone
func (svc *Service) Queue(ctx context.Context, req *QueueRequest) (*QueueResponse, error) {
	id := req.GetId()
	if id == 0 {
		return nil, errors.New("expected id but none given")
	}

	job := svc.EntClient.Job.GetX(ctx, int(id))

	tags := job.QueryTags().IDsX(ctx)

	var tagPredicates []predicate.Target
	for _, tagID := range tags {
		tagPredicates = append(tagPredicates, target.HasTagsWith(tag.ID(tagID)))
	}

	targets := svc.EntClient.Target.Query().
		Where(target.And(tagPredicates...)).
		AllX(ctx)

	currentTime := time.Now()

	var taskIDs []int64

	for _, target := range targets {
		task := svc.EntClient.Task.Create().
			SetQueueTime(currentTime).
			SetContent(job.Content).
			AddTagIDs(tags...).
			SetJobID(job.ID).
			SaveX(ctx)
		target.Update().
			AddTaskIDs(task.ID).
			SaveX(ctx)
		taskIDs = append(taskIDs, int64(task.ID))
	}

	// TODO: publish task queued event here

	return &QueueResponse{Ids: taskIDs}, nil
}
