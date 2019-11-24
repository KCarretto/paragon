package targets

import (
	"context"
	"errors"

	"github.com/kcarretto/paragon/ent"
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
	target := svc.EntClient.Target.GetX(ctx, int(id))

	tasks := target.QueryTasks().IDsX(ctx)

	tags := target.QueryTags().IDsX(ctx)

	credentials := target.QueryCredentials().IDsX(ctx)

	return &FetchResponse{
		Name:        target.Name,
		MachineUUID: target.MachineUUID,
		PrimaryIP:   target.PrimaryIP,
		PublicIP:    target.PublicIP,
		PrimaryMAC:  target.PrimaryMAC,
		Hostname:    target.Hostname,
		LastSeen:    target.LastSeen.Unix(),
		Tasks:       castIntArrayToInt64Array(tasks),
		Tags:        castIntArrayToInt64Array(tags),
		Credentials: castIntArrayToInt64Array(credentials),
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

	ids := svc.EntClient.Target.Query().
		Where(target.Or(
			target.NameContains(req.GetFilter()),
			target.PrimaryIPContains(req.GetFilter()),
			target.PrimaryMACContains(req.GetFilter()),
			target.HostnameContains(req.GetFilter()),
			target.PublicIPContains(req.GetFilter()),
			target.MachineUUIDContains(req.GetFilter()),
			target.HasTagsWith(tag.NameContains(req.GetFilter())),
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

	primaryIP := req.GetPrimaryIP()
	if primaryIP == "" {
		return nil, errors.New("expected primaryIP but none given")
	}

	tags := req.GetTags()
	var tagEnts []*ent.Tag
	for _, tagName := range tags {
		tagEnts = append(tagEnts, svc.EntClient.Tag.Query().
			Where(tag.NameEQ(tagName)).
			FirstX(ctx),
		)
	}
	target := svc.EntClient.Target.Create().
		SetName(name).
		SetPrimaryIP(primaryIP).
		AddTags(tagEnts...).
		SaveX(ctx)
	return &CreateResponse{Id: int64(target.ID)}, nil
}

// AddCredential is used for adding a credential to a target
// TODO: @cictrone
func (svc *Service) AddCredential(ctx context.Context, req *AddCredentialRequest) (*AddCredentialResponse, error) {
	targetID := req.GetTargetID()
	if targetID == 0 {
		return nil, errors.New("expected targetID but none was given")
	}

	credential := svc.EntClient.Credential.Create().
		SetPrincipal(req.GetPrincipal()).
		SetSecret(req.GetSecret()).
		SaveX(ctx)

	svc.EntClient.Target.Update().
		AddCredentialIDs(credential.ID).
		Where(target.ID(int(targetID))).
		SaveX(ctx)
	return &AddCredentialResponse{}, nil
}
