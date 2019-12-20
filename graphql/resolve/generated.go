package resolve

import (
	"context"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/graphql/generated"
	"github.com/kcarretto/paragon/graphql/models"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) ClaimTask(ctx context.Context, input *models.ClaimTaskRequest) (*bool, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Task(ctx context.Context, id *int) (*ent.Task, error) {
	panic("not implemented")
}
func (r *queryResolver) Tasks(ctx context.Context) ([]*int, error) {
	panic("not implemented")
}
