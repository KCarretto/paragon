package tasks

import (
	"context"
	"errors"
	"time"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/job"
	"github.com/kcarretto/paragon/ent/tag"
	"github.com/kcarretto/paragon/ent/task"
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

func castInt64ArrayToIntArray(ints []int64) []int {
	var temp []int
	for _, val := range ints {
		temp = append(temp, int(val))
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
	task := svc.EntClient.Task.GetX(ctx, int(id))
	return &FetchResponse{
		QueueTime:     task.QueueTime.Unix(),
		ClaimTime:     task.ClaimTime.Unix(),
		ExecStartTime: task.ExecStartTime.Unix(),
		ExecStopTime:  task.ExecStopTime.Unix(),
		Content:       task.Content,
		Output:        task.Output,
		Error:         task.Error,
		SessionID:     task.SessionID,
		Tags:          castIntArrayToInt64Array(task.QueryTags().IDsX(ctx)),
		Job:           int64(task.QueryJob().OnlyXID(ctx)),
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

	ids := svc.EntClient.Task.Query().
		Where(task.Or(
			task.SessionIDContains(req.GetFilter()),
			task.HasJobWith(job.NameContains(req.GetFilter())),
			task.ContentContains(req.GetFilter()),
			task.Error(req.GetFilter()),
			task.HasTagsWith(tag.NameContains(req.GetFilter())),
		)).
		Order(ent.Desc(task.FieldQueueTime)).
		Offset(int(offset)).
		Limit(int(limit)).
		IDsX(ctx)

	newOffset := int(offset) + len(ids)
	return &FindResponse{Ids: castIntArrayToInt64Array(ids), NewOffset: int64(newOffset)}, nil
}

// Claim is used for claiming a task
// TODO: @cictrone
func (svc *Service) Claim(ctx context.Context, req *ClaimRequest) (*ClaimResponse, error) {
	id := req.GetId()
	if id == 0 {
		return nil, errors.New("expected id but none given")
	}

	task := svc.EntClient.Task.GetX(ctx, int(id))

	if task.ClaimTime.Unix() > 0 {
		return nil, errors.New("task is already claimed")
	}

	task.Update().
		SetClaimTime(time.Now()).
		SaveX(ctx)

	// TODO: publish event
	return &ClaimResponse{}, nil
}
