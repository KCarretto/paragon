package credentials

import (
	"context"
	"errors"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/credential"
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
	credential := svc.EntClient.Credential.GetX(ctx, int(id))
	return &FetchResponse{
		Principal: credential.Principal,
		Secret:    credential.Secret,
		Fails:     int64(credential.Fails),
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

	ids := svc.EntClient.Credential.Query().
		Where(credential.Or(
			credential.PrincipalContains(req.GetFilter()),
			credential.SecretContains(req.GetFilter()),
		)).
		Offset(int(offset)).
		Limit(int(limit)).
		IDsX(ctx)

	newOffset := int(offset) + len(ids)
	return &FindResponse{Ids: castIntArrayToInt64Array(ids), NewOffset: int64(newOffset)}, nil
}
