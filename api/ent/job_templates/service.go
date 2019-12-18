package job_templates

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/jobtemplate"
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
	template := svc.EntClient.JobTemplate.GetX(ctx, int(id))
	return &FetchResponse{
		Name:    template.Name,
		Content: template.Content,
		Tags:    castIntArrayToInt64Array(template.QueryTags().IDsX(ctx)),
		Jobs:    castIntArrayToInt64Array(template.QueryJobs().IDsX(ctx)),
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

	ids := svc.EntClient.JobTemplate.Query().
		Where(jobtemplate.Or(
			jobtemplate.NameContains(req.GetFilter()),
			jobtemplate.ContentContains(req.GetFilter()),
			jobtemplate.HasTagsWith(tag.NameContains(req.GetFilter())),
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
	template := svc.EntClient.JobTemplate.Create().
		SetName(name).
		SetContent(content).
		AddTags(tagEnts...).
		SaveX(ctx)
	return &CreateResponse{Id: int64(template.ID)}, nil
}

// Queue is used for enqueueing a job
// TODO: @cictrone
func (svc *Service) Queue(ctx context.Context, req *QueueRequest) (*QueueResponse, error) {
	id := req.GetId()
	if id == 0 {
		return nil, errors.New("expected id but none given")
	}

	name := req.GetJobName()
	if name == "" {
		return nil, errors.New("expected job name but none given")
	}

	parameters := req.GetParameters()
	if parameters == "" {
		return nil, errors.New("expected parameters but none given")
	}

	template := svc.EntClient.JobTemplate.GetX(ctx, int(id))

	tags := template.QueryTags().IDsX(ctx)

	var tagPredicates []predicate.Target
	for _, tagID := range tags {
		tagPredicates = append(tagPredicates, target.HasTagsWith(tag.ID(tagID)))
	}

	targets := svc.EntClient.Target.Query().
		Where(target.And(tagPredicates...)).
		AllX(ctx)

	currentTime := time.Now()

	job := svc.EntClient.Job.Create().
		SetName(name).
		SetParameters(parameters).
		SetTemplate(template).
		AddTagIDs(tags...).
		SaveX(ctx)

	var taskIDs []int

	content := template.Content

	// TODO: this should be fancier...

	paramsMap := make(map[string]string)
	err := json.Unmarshal([]byte(parameters), &paramsMap)
	if err != nil {
		return nil, errors.New("parameters passed failed to parse correctly")
	}
	for k, v := range paramsMap {
		content = strings.ReplaceAll(content, k, v)
	}

	for _, target := range targets {
		task := svc.EntClient.Task.Create().
			SetQueueTime(currentTime).
			SetContent(content).
			AddTagIDs(tags...).
			SetJobID(job.ID).
			SaveX(ctx)
		target.Update().
			AddTaskIDs(task.ID).
			SaveX(ctx)
		taskIDs = append(taskIDs, task.ID)
	}

	// TODO: publish task queued event here

	return &QueueResponse{Id: int64(job.ID)}, nil
}
