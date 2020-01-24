package resolve

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/credential"
	"github.com/kcarretto/paragon/ent/target"
	"github.com/kcarretto/paragon/ent/task"
	"github.com/kcarretto/paragon/graphql/generated"
	"github.com/kcarretto/paragon/graphql/models"
	"github.com/kcarretto/paragon/pkg/event"
	"go.uber.org/zap"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

// Resolver is the root struct for handling all resolves
type Resolver struct {
	Log    *zap.Logger
	Graph  *ent.Client
	Events event.Publisher
}

// Credential is the Resolver for the Credential Ent
func (r *Resolver) Credential() generated.CredentialResolver {
	return &credentialResolver{r}
}

// File is the Resolver for the File Ent
func (r *Resolver) File() generated.FileResolver {
	return &fileResolver{r}
}

// Job is the Resolver for the Job Ent
func (r *Resolver) Job() generated.JobResolver {
	return &jobResolver{r}
}

// Link is the Resolver for the Link Ent
func (r *Resolver) Link() generated.LinkResolver {
	return &linkResolver{r}
}

// Mutation is the Resolver for the Mutations
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

// Query is the Resolver for the Queries
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

// Tag is the Resolver for the Tag Ent
func (r *Resolver) Tag() generated.TagResolver {
	return &tagResolver{r}
}

// Target is the Resolver for the Target Ent
func (r *Resolver) Target() generated.TargetResolver {
	return &targetResolver{r}
}

// Task is the Resolver for the Task Ent
func (r *Resolver) Task() generated.TaskResolver {
	return &taskResolver{r}
}

type credentialResolver struct{ *Resolver }

func (r *credentialResolver) Kind(ctx context.Context, obj *ent.Credential) (*string, error) {
	kind := obj.Kind.String()
	return &kind, nil
}

type fileResolver struct{ *Resolver }

func (r *fileResolver) Links(ctx context.Context, obj *ent.File, input *models.Filter) ([]*ent.Link, error) {
	q := obj.QueryLinks()
	if input != nil {
		if input.Offset != nil {
			q.Offset(*input.Offset)
		}
		if input.Limit != nil {
			q.Limit(*input.Limit)
		}
	}
	return q.All(ctx)
}

type jobResolver struct{ *Resolver }

func (r *jobResolver) Tasks(ctx context.Context, obj *ent.Job, input *models.Filter) ([]*ent.Task, error) {
	q := obj.QueryTasks()
	if input != nil {
		if input.Offset != nil {
			q.Offset(*input.Offset)
		}
		if input.Limit != nil {
			q.Limit(*input.Limit)
		}
	}
	return q.Order(ent.Desc(task.FieldLastChangedTime)).All(ctx)
}
func (r *jobResolver) Tags(ctx context.Context, obj *ent.Job, input *models.Filter) ([]*ent.Tag, error) {
	q := obj.QueryTags()
	if input != nil {
		if input.Offset != nil {
			q.Offset(*input.Offset)
		}
		if input.Limit != nil {
			q.Limit(*input.Limit)
		}
	}
	return q.All(ctx)
}
func (r *jobResolver) Next(ctx context.Context, obj *ent.Job) (*ent.Job, error) {
	return obj.QueryNext().Only(ctx)
}
func (r *jobResolver) Prev(ctx context.Context, obj *ent.Job) (*ent.Job, error) {
	return obj.QueryPrev().Only(ctx)
}

type linkResolver struct{ *Resolver }

func (r *linkResolver) File(ctx context.Context, obj *ent.Link) (*ent.File, error) {
	return obj.QueryFile().Only(ctx)
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) FailCredential(ctx context.Context, input *models.FailCredentialRequest) (*ent.Credential, error) {
	cred, err := r.Graph.Credential.GetX(ctx, input.ID).
		Update().
		AddFails(1).
		Save(ctx)
	return cred, err
}
func (r *mutationResolver) CreateJob(ctx context.Context, input *models.CreateJobRequest) (*ent.Job, error) {
	jobCreator := r.Graph.Job.Create().
		SetName(input.Name).
		SetContent(input.Content).
		AddTagIDs(input.Tags...)
	if input.Prev != nil {
		jobCreator.SetPrevID(*input.Prev)
	}

	var targets []*ent.Target
	for _, t := range input.Targets {
		tar, err := r.Graph.Target.Get(ctx, t)
		if err != nil {
			return nil, err
		}
		targets = append(targets, tar)
	}

	// just so all tasks and job are the same time.
	currentTime := time.Now()

	job, err := jobCreator.
		SetCreationTime(currentTime).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if len(targets) == 0 {
		t, err := r.Graph.Task.Create().
			SetQueueTime(currentTime).
			SetLastChangedTime(currentTime).
			SetContent(input.Content).
			SetNillableSessionID(input.SessionID).
			AddTagIDs(input.Tags...).
			SetJobID(job.ID).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		tags, err := t.QueryTags().All(ctx)
		if err != nil {
			return nil, err
		}
		e := event.TaskQueued{
			Target:      nil,
			Task:        t,
			Credentials: nil,
			Tags:        tags,
		}
		d, err := json.Marshal(e)
		if err != nil {
			return nil, err
		}
		err = r.Events.Publish(ctx, d)
		if err != nil {
			return nil, err
		}
	} else {
		for _, target := range targets {
			task, err := r.Graph.Task.Create().
				SetQueueTime(currentTime).
				SetLastChangedTime(currentTime).
				SetContent(input.Content).
				SetNillableSessionID(input.SessionID).
				AddTagIDs(input.Tags...).
				SetJobID(job.ID).
				Save(ctx)
			if err != nil {
				return nil, err
			}
			t, err := target.Update().
				AddTaskIDs(task.ID).
				Save(ctx)
			if err != nil {
				return nil, err
			}
			tags, err := task.QueryTags().All(ctx)
			if err != nil {
				return nil, err
			}
			creds, err := t.QueryCredentials().All(ctx)
			if err != nil {
				return nil, err
			}
			e := event.TaskQueued{
				Target:      t,
				Task:        task,
				Credentials: creds,
				Tags:        tags,
			}
			d, err := json.Marshal(e)
			if err != nil {
				return nil, err
			}
			err = r.Events.Publish(ctx, d)
			if err != nil {
				return nil, err
			}
		}
	}

	return job, nil
}
func (r *mutationResolver) CreateTag(ctx context.Context, input *models.CreateTagRequest) (*ent.Tag, error) {
	tag, err := r.Graph.Tag.Create().
		SetName(input.Name).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return tag, nil
}
func (r *mutationResolver) ApplyTagToTask(ctx context.Context, input *models.ApplyTagRequest) (*ent.Task, error) {
	task, err := r.Graph.Task.UpdateOneID(input.EntID).
		AddTagIDs(input.TagID).
		Save(ctx)
	return task, err
}
func (r *mutationResolver) ApplyTagToTarget(ctx context.Context, input *models.ApplyTagRequest) (*ent.Target, error) {
	target, err := r.Graph.Target.UpdateOneID(input.EntID).
		AddTagIDs(input.TagID).
		Save(ctx)
	return target, err
}
func (r *mutationResolver) ApplyTagToJob(ctx context.Context, input *models.ApplyTagRequest) (*ent.Job, error) {
	job, err := r.Graph.Job.UpdateOneID(input.EntID).
		AddTagIDs(input.TagID).
		Save(ctx)
	return job, err
}
func (r *mutationResolver) RemoveTagFromTask(ctx context.Context, input *models.RemoveTagRequest) (*ent.Task, error) {
	task, err := r.Graph.Task.UpdateOneID(input.EntID).
		RemoveTagIDs(input.TagID).
		Save(ctx)
	return task, err
}
func (r *mutationResolver) RemoveTagFromTarget(ctx context.Context, input *models.RemoveTagRequest) (*ent.Target, error) {
	target, err := r.Graph.Target.UpdateOneID(input.EntID).
		RemoveTagIDs(input.TagID).
		Save(ctx)
	return target, err
}
func (r *mutationResolver) RemoveTagFromJob(ctx context.Context, input *models.RemoveTagRequest) (*ent.Job, error) {
	job, err := r.Graph.Job.UpdateOneID(input.EntID).
		RemoveTagIDs(input.TagID).
		Save(ctx)
	return job, err
}
func (r *mutationResolver) CreateTarget(ctx context.Context, input *models.CreateTargetRequest) (*ent.Target, error) {
	target, err := r.Graph.Target.Create().
		SetName(input.Name).
		SetPrimaryIP(input.PrimaryIP).
		AddTagIDs(input.Tags...).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return target, nil
}
func (r *mutationResolver) SetTargetFields(ctx context.Context, input *models.SetTargetFieldsRequest) (*ent.Target, error) {
	targetUpdater := r.Graph.Target.UpdateOneID(input.ID)
	if input.Name != nil {
		targetUpdater.SetName(*input.Name)
	}
	if input.PrimaryIP != nil {
		targetUpdater.SetPrimaryIP(*input.PrimaryIP)
	}
	target, err := targetUpdater.
		SetNillableHostname(input.Hostname).
		SetNillableMachineUUID(input.MachineUUID).
		SetNillablePrimaryMAC(input.PrimaryMac).
		SetNillablePublicIP(input.PublicIP).
		Save(ctx)
	return target, err
}
func (r *mutationResolver) DeleteTarget(ctx context.Context, input *models.DeleteTargetRequest) (bool, error) {
	err := r.Graph.Target.DeleteOneID(input.ID).Exec(ctx)
	return err != nil, err
}
func (r *mutationResolver) AddCredentialForTarget(ctx context.Context, input *models.AddCredentialForTargetRequest) (*ent.Target, error) {
	kind := credential.KindPassword
	if input.Kind != nil {
		kind = credential.Kind(*input.Kind)
	}
	credential, err := r.Graph.Credential.Create().
		SetPrincipal(input.Principal).
		SetSecret(input.Secret).
		SetKind(kind).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	target, err := r.Graph.Target.UpdateOneID(input.ID).
		AddCredentialIDs(credential.ID).
		Save(ctx)
	return target, err
}
func (r *mutationResolver) ClaimTasks(ctx context.Context, input *models.ClaimTasksRequest) ([]*ent.Task, error) {
	var targetEnt *ent.Target
	var err error

	// check for valid machineuuid
	if input.MachineUUID != nil && *input.MachineUUID != "" {
		targetEnt, err = r.Graph.Target.Query().
			Where(target.MachineUUID(*input.MachineUUID)).
			Only(ctx)
	}

	// chack for valid primaryIP (if we didnt find a target yet)
	if targetEnt == nil && input.PrimaryIP != nil && *input.PrimaryIP != "" {
		targetEnt, err = r.Graph.Target.Query().
			Where(target.PrimaryIP(*input.PrimaryIP)).
			Only(ctx)
		if err != nil {
			return nil, err
		}
	}

	// if still no target, fail
	if targetEnt == nil {
		return nil, fmt.Errorf("neither valid machineUUID nor primaryIP was given")
	}

	// Prepare Target mutation
	currentTime := time.Now()
	targetMutation := targetEnt.Update().
		SetNillableHostname(input.Hostname).
		SetNillablePrimaryMAC(input.PrimaryMac).
		SetLastSeen(currentTime)

	// Update primary IP if available
	if input.PrimaryIP != nil && *input.PrimaryIP != "" {
		// TODO: Validate IP
		targetMutation.SetPrimaryIP(*input.PrimaryIP)
	}

	// Update MachineUUID if available
	if input.MachineUUID != nil && *input.MachineUUID != "" {
		targetMutation.SetMachineUUID(*input.MachineUUID)
	}

	// set lastSeen on target
	targetEnt, err = targetMutation.Save(ctx)
	if err != nil {
		return nil, err
	}

	// get all tasks for target not claimed
	tasks, err := targetEnt.QueryTasks().
		Where(
			task.ClaimTimeIsNil(),
		).
		All(ctx)

	if err != nil {
		return nil, err
	}

	var sessionID string
	if input.SessionID != nil && *input.SessionID != "" {
		sessionID = *input.SessionID
	}

	// set claimtime on all tasks
	var updatedTasks []*ent.Task
	for _, t := range tasks {
		// filter by session if necessary
		if t.SessionID == "" || t.SessionID == sessionID {
			t, err = t.Update().
				SetLastChangedTime(currentTime).
				SetClaimTime(currentTime).
				Save(ctx)
			if err != nil {
				return nil, err
			}
			updatedTasks = append(updatedTasks, t)
		}
	}
	return updatedTasks, nil
}
func (r *mutationResolver) ClaimTask(ctx context.Context, id int) (*ent.Task, error) {
	task, err := r.Graph.Task.Query().
		Where(
			task.ClaimTimeIsNil(),
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	return task.Update().
		SetClaimTime(currentTime).
		SetLastChangedTime(currentTime).
		Save(ctx)
}
func (r *mutationResolver) SubmitTaskResult(ctx context.Context, input *models.SubmitTaskResultRequest) (*ent.Task, error) {
	taskEnt, err := r.Graph.Task.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	inputOutput := ""
	if input.Output != nil {
		inputOutput = *input.Output
	}
	inputError := ""
	if input.Error != nil {
		inputError = *input.Error
	}
	_, err = r.Graph.Target.Update().
		SetLastSeen(time.Now()).
		Where(target.HasTasksWith(task.ID(taskEnt.ID))).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return taskEnt.Update().
		SetLastChangedTime(time.Now()).
		SetOutput(taskEnt.Output + inputOutput).
		SetError(taskEnt.Error + inputError).
		SetNillableExecStartTime(input.ExecStartTime).
		SetNillableExecStopTime(input.ExecStopTime).
		Save(ctx)
}
func (r *mutationResolver) CreateLink(ctx context.Context, input *models.CreateLinkRequest) (*ent.Link, error) {
	linkCreator := r.Graph.Link.Create().
		SetAlias(input.Alias).
		SetFileID(input.File)
	if input.ExpirationTime != nil {
		linkCreator.SetExpirationTime(*input.ExpirationTime)
	}
	if input.Clicks != nil {
		linkCreator.SetClicks(*input.Clicks)
	}
	return linkCreator.Save(ctx)
}
func (r *mutationResolver) SetLinkFields(ctx context.Context, input *models.SetLinkFieldsRequest) (*ent.Link, error) {
	linkUpdater := r.Graph.Link.UpdateOneID(input.ID)
	if input.Alias != nil {
		linkUpdater.SetAlias(*input.Alias)
	}
	if input.ExpirationTime != nil {
		linkUpdater.SetExpirationTime(*input.ExpirationTime)
	}
	if input.Clicks != nil {
		linkUpdater.SetClicks(*input.Clicks)
	}
	return linkUpdater.Save(ctx)
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Link(ctx context.Context, id int) (*ent.Link, error) {
	return r.Graph.Link.Get(ctx, id)
}
func (r *queryResolver) Links(ctx context.Context, input *models.Filter) ([]*ent.Link, error) {
	q := r.Graph.Link.Query()
	if input != nil {
		if input.Offset != nil {
			q.Offset(*input.Offset)
		}
		if input.Limit != nil {
			q.Limit(*input.Limit)
		}
	}
	return q.All(ctx)
}
func (r *queryResolver) File(ctx context.Context, id int) (*ent.File, error) {
	return r.Graph.File.Get(ctx, id)
}
func (r *queryResolver) Files(ctx context.Context, input *models.Filter) ([]*ent.File, error) {
	q := r.Graph.File.Query()
	if input != nil {
		if input.Offset != nil {
			q.Offset(*input.Offset)
		}
		if input.Limit != nil {
			q.Limit(*input.Limit)
		}
	}
	return q.All(ctx)
}
func (r *queryResolver) Credential(ctx context.Context, id int) (*ent.Credential, error) {
	return r.Graph.Credential.Get(ctx, id)
}
func (r *queryResolver) Credentials(ctx context.Context, input *models.Filter) ([]*ent.Credential, error) {
	q := r.Graph.Credential.Query()
	if input != nil {
		if input.Offset != nil {
			q.Offset(*input.Offset)
		}
		if input.Limit != nil {
			q.Limit(*input.Limit)
		}
	}
	return q.All(ctx)
}
func (r *queryResolver) Job(ctx context.Context, id int) (*ent.Job, error) {
	return r.Graph.Job.Get(ctx, id)
}
func (r *queryResolver) Jobs(ctx context.Context, input *models.Filter) ([]*ent.Job, error) {
	q := r.Graph.Job.Query()
	if input != nil {
		if input.Offset != nil {
			q.Offset(*input.Offset)
		}
		if input.Limit != nil {
			q.Limit(*input.Limit)
		}
	}
	return q.All(ctx)
}
func (r *queryResolver) Tag(ctx context.Context, id int) (*ent.Tag, error) {
	return r.Graph.Tag.Get(ctx, id)
}
func (r *queryResolver) Tags(ctx context.Context, input *models.Filter) ([]*ent.Tag, error) {
	q := r.Graph.Tag.Query()
	if input != nil {
		if input.Offset != nil {
			q.Offset(*input.Offset)
		}
		if input.Limit != nil {
			q.Limit(*input.Limit)
		}
	}
	return q.All(ctx)
}
func (r *queryResolver) Target(ctx context.Context, id int) (*ent.Target, error) {
	return r.Graph.Target.Get(ctx, id)
}
func (r *queryResolver) Targets(ctx context.Context, input *models.Filter) ([]*ent.Target, error) {
	q := r.Graph.Target.Query()
	if input != nil {
		if input.Offset != nil {
			q.Offset(*input.Offset)
		}
		if input.Limit != nil {
			q.Limit(*input.Limit)
		}
	}
	return q.All(ctx)
}
func (r *queryResolver) Task(ctx context.Context, id int) (*ent.Task, error) {
	return r.Graph.Task.Get(ctx, id)
}
func (r *queryResolver) Tasks(ctx context.Context, input *models.Filter) ([]*ent.Task, error) {
	q := r.Graph.Task.Query()
	if input != nil {
		if input.Offset != nil {
			q.Offset(*input.Offset)
		}
		if input.Limit != nil {
			q.Limit(*input.Limit)
		}
	}
	return q.Order(ent.Desc(task.FieldLastChangedTime)).All(ctx)
}

type tagResolver struct{ *Resolver }

func (r *tagResolver) Tasks(ctx context.Context, obj *ent.Tag, input *models.Filter) ([]*ent.Task, error) {
	q := obj.QueryTasks()
	if input != nil {
		if input.Offset != nil {
			q.Offset(*input.Offset)
		}
		if input.Limit != nil {
			q.Limit(*input.Limit)
		}
	}
	return q.Order(ent.Desc(task.FieldLastChangedTime)).All(ctx)
}
func (r *tagResolver) Targets(ctx context.Context, obj *ent.Tag, input *models.Filter) ([]*ent.Target, error) {
	q := obj.QueryTargets()
	if input != nil {
		if input.Offset != nil {
			q.Offset(*input.Offset)
		}
		if input.Limit != nil {
			q.Limit(*input.Limit)
		}
	}
	return q.All(ctx)
}
func (r *tagResolver) Jobs(ctx context.Context, obj *ent.Tag, input *models.Filter) ([]*ent.Job, error) {
	q := obj.QueryJobs()
	if input != nil {
		if input.Offset != nil {
			q.Offset(*input.Offset)
		}
		if input.Limit != nil {
			q.Limit(*input.Limit)
		}
	}
	return q.All(ctx)
}

type targetResolver struct{ *Resolver }

func (r *targetResolver) Tasks(ctx context.Context, obj *ent.Target, input *models.Filter) ([]*ent.Task, error) {
	q := obj.QueryTasks()
	if input != nil {
		if input.Offset != nil {
			q.Offset(*input.Offset)
		}
		if input.Limit != nil {
			q.Limit(*input.Limit)
		}
	}
	return q.Order(ent.Desc(task.FieldLastChangedTime)).All(ctx)
}
func (r *targetResolver) Tags(ctx context.Context, obj *ent.Target, input *models.Filter) ([]*ent.Tag, error) {
	q := obj.QueryTags()
	if input != nil {
		if input.Offset != nil {
			q.Offset(*input.Offset)
		}
		if input.Limit != nil {
			q.Limit(*input.Limit)
		}
	}
	return q.All(ctx)
}
func (r *targetResolver) Credentials(ctx context.Context, obj *ent.Target, input *models.Filter) ([]*ent.Credential, error) {
	q := obj.QueryCredentials()
	if input != nil {
		if input.Offset != nil {
			q.Offset(*input.Offset)
		}
		if input.Limit != nil {
			q.Limit(*input.Limit)
		}
	}
	return q.All(ctx)
}

type taskResolver struct{ *Resolver }

func (r *taskResolver) Job(ctx context.Context, obj *ent.Task) (*ent.Job, error) {
	return obj.QueryJob().Only(ctx)
}
func (r *taskResolver) Target(ctx context.Context, obj *ent.Task) (*ent.Target, error) {
	return obj.QueryTarget().Only(ctx)
}
