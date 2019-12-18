package jobs

import (
	"context"
	"errors"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/job"
	"github.com/kcarretto/paragon/ent/jobtemplate"
	"github.com/kcarretto/paragon/ent/tag"
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
		Name:       job.Name,
		Parameters: job.Parameters,
		Tasks:      castIntArrayToInt64Array(job.QueryTasks().IDsX(ctx)),
		Tags:       castIntArrayToInt64Array(job.QueryTags().IDsX(ctx)),
		Template:   int64(job.QueryTemplate().OnlyXID(ctx)),
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
			job.HasTagsWith(tag.NameContains(req.GetFilter())),
			job.HasTemplateWith(jobtemplate.ContentContains(req.GetFilter())),
		)).
		Offset(int(offset)).
		Limit(int(limit)).
		IDsX(ctx)

	newOffset := int(offset) + len(ids)
	return &FindResponse{Ids: castIntArrayToInt64Array(ids), NewOffset: int64(newOffset)}, nil
}
