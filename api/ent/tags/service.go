package tags

import (
	"context"
	"errors"

	"github.com/kcarretto/paragon/ent"
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
	tag := svc.EntClient.Tag.GetX(ctx, int(id))
	return &FetchResponse{
		Name: tag.Name,
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

	ids := svc.EntClient.Tag.Query().
		Where(tag.NameContains(req.GetFilter())).
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
	tag := svc.EntClient.Tag.Create().
		SetName(name).
		SaveX(ctx)
	return &CreateResponse{Id: int64(tag.ID)}, nil
}

// ApplyToTask is used for applying a tag to an task
// TODO: @cictrone
func (svc *Service) ApplyToTask(ctx context.Context, req *ApplyRequest) (*ApplyResponse, error) {
	tagID := req.GetTagId()
	if tagID == 0 {
		return nil, errors.New("expected tagID but none given")
	}
	entID := req.GetEntId()
	if entID == 0 {
		return nil, errors.New("expected entID but none given")
	}
	tag := svc.EntClient.Tag.GetX(ctx, int(tagID))
	svc.EntClient.Task.GetX(ctx, int(entID)).
		Update().
		AddTagIDs(tag.ID).
		SaveX(ctx)
	return &ApplyResponse{}, nil
}

// ApplyToTarget is used for applying a tag to an target
// TODO: @cictrone
func (svc *Service) ApplyToTarget(ctx context.Context, req *ApplyRequest) (*ApplyResponse, error) {
	tagID := req.GetTagId()
	if tagID == 0 {
		return nil, errors.New("expected tagID but none given")
	}
	entID := req.GetEntId()
	if entID == 0 {
		return nil, errors.New("expected entID but none given")
	}
	tag := svc.EntClient.Tag.GetX(ctx, int(tagID))
	svc.EntClient.Target.GetX(ctx, int(entID)).
		Update().
		AddTagIDs(tag.ID).
		SaveX(ctx)
	return &ApplyResponse{}, nil
}

// ApplyToJob is used for applying a tag to an job
// TODO: @cictrone
func (svc *Service) ApplyToJob(ctx context.Context, req *ApplyRequest) (*ApplyResponse, error) {
	tagID := req.GetTagId()
	if tagID == 0 {
		return nil, errors.New("expected tagID but none given")
	}
	entID := req.GetEntId()
	if entID == 0 {
		return nil, errors.New("expected entID but none given")
	}
	tag := svc.EntClient.Tag.GetX(ctx, int(tagID))
	svc.EntClient.Job.GetX(ctx, int(entID)).
		Update().
		AddTagIDs(tag.ID).
		SaveX(ctx)
	return &ApplyResponse{}, nil
}

// RemoveFromTask is used for removing a tag from an task
// TODO: @cictrone
func (svc *Service) RemoveFromTask(ctx context.Context, req *RemoveRequest) (*RemoveResponse, error) {
	tagID := req.GetTagId()
	if tagID == 0 {
		return nil, errors.New("expected tagID but none given")
	}
	entID := req.GetEntId()
	if entID == 0 {
		return nil, errors.New("expected entID but none given")
	}
	tag := svc.EntClient.Tag.GetX(ctx, int(tagID))
	svc.EntClient.Task.GetX(ctx, int(entID)).
		Update().
		RemoveTagIDs(tag.ID).
		SaveX(ctx)
	return &RemoveResponse{}, nil
}

// RemoveFromTarget is used for removing a tag from an target
// TODO: @cictrone
func (svc *Service) RemoveFromTarget(ctx context.Context, req *RemoveRequest) (*RemoveResponse, error) {
	tagID := req.GetTagId()
	if tagID == 0 {
		return nil, errors.New("expected tagID but none given")
	}
	entID := req.GetEntId()
	if entID == 0 {
		return nil, errors.New("expected entID but none given")
	}
	tag := svc.EntClient.Tag.GetX(ctx, int(tagID))
	svc.EntClient.Target.GetX(ctx, int(entID)).
		Update().
		RemoveTagIDs(tag.ID).
		SaveX(ctx)
	return &RemoveResponse{}, nil
}

// RemoveFromJob is used for removing a tag from an job
// TODO: @cictrone
func (svc *Service) RemoveFromJob(ctx context.Context, req *RemoveRequest) (*RemoveResponse, error) {
	tagID := req.GetTagId()
	if tagID == 0 {
		return nil, errors.New("expected tagID but none given")
	}
	entID := req.GetEntId()
	if entID == 0 {
		return nil, errors.New("expected entID but none given")
	}
	tag := svc.EntClient.Tag.GetX(ctx, int(tagID))
	svc.EntClient.Job.GetX(ctx, int(entID)).
		Update().
		RemoveTagIDs(tag.ID).
		SaveX(ctx)
	return &RemoveResponse{}, nil
}
