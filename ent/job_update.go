// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/kcarretto/paragon/ent/job"
	"github.com/kcarretto/paragon/ent/predicate"
	"github.com/kcarretto/paragon/ent/tag"
	"github.com/kcarretto/paragon/ent/task"
	"github.com/kcarretto/paragon/ent/user"
)

// JobUpdate is the builder for updating Job entities.
type JobUpdate struct {
	config
	Name         *string
	CreationTime *time.Time
	Content      *string
	Staged       *bool
	tasks        map[int]struct{}
	tags         map[int]struct{}
	prev         map[int]struct{}
	next         map[int]struct{}
	owner        map[int]struct{}
	removedTasks map[int]struct{}
	removedTags  map[int]struct{}
	clearedPrev  bool
	clearedNext  bool
	clearedOwner bool
	predicates   []predicate.Job
}

// Where adds a new predicate for the builder.
func (ju *JobUpdate) Where(ps ...predicate.Job) *JobUpdate {
	ju.predicates = append(ju.predicates, ps...)
	return ju
}

// SetName sets the Name field.
func (ju *JobUpdate) SetName(s string) *JobUpdate {
	ju.Name = &s
	return ju
}

// SetCreationTime sets the CreationTime field.
func (ju *JobUpdate) SetCreationTime(t time.Time) *JobUpdate {
	ju.CreationTime = &t
	return ju
}

// SetNillableCreationTime sets the CreationTime field if the given value is not nil.
func (ju *JobUpdate) SetNillableCreationTime(t *time.Time) *JobUpdate {
	if t != nil {
		ju.SetCreationTime(*t)
	}
	return ju
}

// SetContent sets the Content field.
func (ju *JobUpdate) SetContent(s string) *JobUpdate {
	ju.Content = &s
	return ju
}

// SetStaged sets the Staged field.
func (ju *JobUpdate) SetStaged(b bool) *JobUpdate {
	ju.Staged = &b
	return ju
}

// AddTaskIDs adds the tasks edge to Task by ids.
func (ju *JobUpdate) AddTaskIDs(ids ...int) *JobUpdate {
	if ju.tasks == nil {
		ju.tasks = make(map[int]struct{})
	}
	for i := range ids {
		ju.tasks[ids[i]] = struct{}{}
	}
	return ju
}

// AddTasks adds the tasks edges to Task.
func (ju *JobUpdate) AddTasks(t ...*Task) *JobUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return ju.AddTaskIDs(ids...)
}

// AddTagIDs adds the tags edge to Tag by ids.
func (ju *JobUpdate) AddTagIDs(ids ...int) *JobUpdate {
	if ju.tags == nil {
		ju.tags = make(map[int]struct{})
	}
	for i := range ids {
		ju.tags[ids[i]] = struct{}{}
	}
	return ju
}

// AddTags adds the tags edges to Tag.
func (ju *JobUpdate) AddTags(t ...*Tag) *JobUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return ju.AddTagIDs(ids...)
}

// SetPrevID sets the prev edge to Job by id.
func (ju *JobUpdate) SetPrevID(id int) *JobUpdate {
	if ju.prev == nil {
		ju.prev = make(map[int]struct{})
	}
	ju.prev[id] = struct{}{}
	return ju
}

// SetNillablePrevID sets the prev edge to Job by id if the given value is not nil.
func (ju *JobUpdate) SetNillablePrevID(id *int) *JobUpdate {
	if id != nil {
		ju = ju.SetPrevID(*id)
	}
	return ju
}

// SetPrev sets the prev edge to Job.
func (ju *JobUpdate) SetPrev(j *Job) *JobUpdate {
	return ju.SetPrevID(j.ID)
}

// SetNextID sets the next edge to Job by id.
func (ju *JobUpdate) SetNextID(id int) *JobUpdate {
	if ju.next == nil {
		ju.next = make(map[int]struct{})
	}
	ju.next[id] = struct{}{}
	return ju
}

// SetNillableNextID sets the next edge to Job by id if the given value is not nil.
func (ju *JobUpdate) SetNillableNextID(id *int) *JobUpdate {
	if id != nil {
		ju = ju.SetNextID(*id)
	}
	return ju
}

// SetNext sets the next edge to Job.
func (ju *JobUpdate) SetNext(j *Job) *JobUpdate {
	return ju.SetNextID(j.ID)
}

// SetOwnerID sets the owner edge to User by id.
func (ju *JobUpdate) SetOwnerID(id int) *JobUpdate {
	if ju.owner == nil {
		ju.owner = make(map[int]struct{})
	}
	ju.owner[id] = struct{}{}
	return ju
}

// SetOwner sets the owner edge to User.
func (ju *JobUpdate) SetOwner(u *User) *JobUpdate {
	return ju.SetOwnerID(u.ID)
}

// RemoveTaskIDs removes the tasks edge to Task by ids.
func (ju *JobUpdate) RemoveTaskIDs(ids ...int) *JobUpdate {
	if ju.removedTasks == nil {
		ju.removedTasks = make(map[int]struct{})
	}
	for i := range ids {
		ju.removedTasks[ids[i]] = struct{}{}
	}
	return ju
}

// RemoveTasks removes tasks edges to Task.
func (ju *JobUpdate) RemoveTasks(t ...*Task) *JobUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return ju.RemoveTaskIDs(ids...)
}

// RemoveTagIDs removes the tags edge to Tag by ids.
func (ju *JobUpdate) RemoveTagIDs(ids ...int) *JobUpdate {
	if ju.removedTags == nil {
		ju.removedTags = make(map[int]struct{})
	}
	for i := range ids {
		ju.removedTags[ids[i]] = struct{}{}
	}
	return ju
}

// RemoveTags removes tags edges to Tag.
func (ju *JobUpdate) RemoveTags(t ...*Tag) *JobUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return ju.RemoveTagIDs(ids...)
}

// ClearPrev clears the prev edge to Job.
func (ju *JobUpdate) ClearPrev() *JobUpdate {
	ju.clearedPrev = true
	return ju
}

// ClearNext clears the next edge to Job.
func (ju *JobUpdate) ClearNext() *JobUpdate {
	ju.clearedNext = true
	return ju
}

// ClearOwner clears the owner edge to User.
func (ju *JobUpdate) ClearOwner() *JobUpdate {
	ju.clearedOwner = true
	return ju
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (ju *JobUpdate) Save(ctx context.Context) (int, error) {
	if ju.Name != nil {
		if err := job.NameValidator(*ju.Name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"Name\": %v", err)
		}
	}
	if ju.Content != nil {
		if err := job.ContentValidator(*ju.Content); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"Content\": %v", err)
		}
	}
	if len(ju.prev) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"prev\"")
	}
	if len(ju.next) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"next\"")
	}
	if len(ju.owner) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"owner\"")
	}
	if ju.clearedOwner && ju.owner == nil {
		return 0, errors.New("ent: clearing a unique edge \"owner\"")
	}
	return ju.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (ju *JobUpdate) SaveX(ctx context.Context) int {
	affected, err := ju.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ju *JobUpdate) Exec(ctx context.Context) error {
	_, err := ju.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ju *JobUpdate) ExecX(ctx context.Context) {
	if err := ju.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ju *JobUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   job.Table,
			Columns: job.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: job.FieldID,
			},
		},
	}
	if ps := ju.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := ju.Name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: job.FieldName,
		})
	}
	if value := ju.CreationTime; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: job.FieldCreationTime,
		})
	}
	if value := ju.Content; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: job.FieldContent,
		})
	}
	if value := ju.Staged; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: job.FieldStaged,
		})
	}
	if nodes := ju.removedTasks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   job.TasksTable,
			Columns: []string{job.TasksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: task.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ju.tasks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   job.TasksTable,
			Columns: []string{job.TasksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: task.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := ju.removedTags; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   job.TagsTable,
			Columns: job.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: tag.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ju.tags; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   job.TagsTable,
			Columns: job.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: tag.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ju.clearedPrev {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   job.PrevTable,
			Columns: []string{job.PrevColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: job.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ju.prev; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   job.PrevTable,
			Columns: []string{job.PrevColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: job.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ju.clearedNext {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   job.NextTable,
			Columns: []string{job.NextColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: job.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ju.next; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   job.NextTable,
			Columns: []string{job.NextColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: job.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ju.clearedOwner {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   job.OwnerTable,
			Columns: []string{job.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ju.owner; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   job.OwnerTable,
			Columns: []string{job.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ju.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// JobUpdateOne is the builder for updating a single Job entity.
type JobUpdateOne struct {
	config
	id           int
	Name         *string
	CreationTime *time.Time
	Content      *string
	Staged       *bool
	tasks        map[int]struct{}
	tags         map[int]struct{}
	prev         map[int]struct{}
	next         map[int]struct{}
	owner        map[int]struct{}
	removedTasks map[int]struct{}
	removedTags  map[int]struct{}
	clearedPrev  bool
	clearedNext  bool
	clearedOwner bool
}

// SetName sets the Name field.
func (juo *JobUpdateOne) SetName(s string) *JobUpdateOne {
	juo.Name = &s
	return juo
}

// SetCreationTime sets the CreationTime field.
func (juo *JobUpdateOne) SetCreationTime(t time.Time) *JobUpdateOne {
	juo.CreationTime = &t
	return juo
}

// SetNillableCreationTime sets the CreationTime field if the given value is not nil.
func (juo *JobUpdateOne) SetNillableCreationTime(t *time.Time) *JobUpdateOne {
	if t != nil {
		juo.SetCreationTime(*t)
	}
	return juo
}

// SetContent sets the Content field.
func (juo *JobUpdateOne) SetContent(s string) *JobUpdateOne {
	juo.Content = &s
	return juo
}

// SetStaged sets the Staged field.
func (juo *JobUpdateOne) SetStaged(b bool) *JobUpdateOne {
	juo.Staged = &b
	return juo
}

// AddTaskIDs adds the tasks edge to Task by ids.
func (juo *JobUpdateOne) AddTaskIDs(ids ...int) *JobUpdateOne {
	if juo.tasks == nil {
		juo.tasks = make(map[int]struct{})
	}
	for i := range ids {
		juo.tasks[ids[i]] = struct{}{}
	}
	return juo
}

// AddTasks adds the tasks edges to Task.
func (juo *JobUpdateOne) AddTasks(t ...*Task) *JobUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return juo.AddTaskIDs(ids...)
}

// AddTagIDs adds the tags edge to Tag by ids.
func (juo *JobUpdateOne) AddTagIDs(ids ...int) *JobUpdateOne {
	if juo.tags == nil {
		juo.tags = make(map[int]struct{})
	}
	for i := range ids {
		juo.tags[ids[i]] = struct{}{}
	}
	return juo
}

// AddTags adds the tags edges to Tag.
func (juo *JobUpdateOne) AddTags(t ...*Tag) *JobUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return juo.AddTagIDs(ids...)
}

// SetPrevID sets the prev edge to Job by id.
func (juo *JobUpdateOne) SetPrevID(id int) *JobUpdateOne {
	if juo.prev == nil {
		juo.prev = make(map[int]struct{})
	}
	juo.prev[id] = struct{}{}
	return juo
}

// SetNillablePrevID sets the prev edge to Job by id if the given value is not nil.
func (juo *JobUpdateOne) SetNillablePrevID(id *int) *JobUpdateOne {
	if id != nil {
		juo = juo.SetPrevID(*id)
	}
	return juo
}

// SetPrev sets the prev edge to Job.
func (juo *JobUpdateOne) SetPrev(j *Job) *JobUpdateOne {
	return juo.SetPrevID(j.ID)
}

// SetNextID sets the next edge to Job by id.
func (juo *JobUpdateOne) SetNextID(id int) *JobUpdateOne {
	if juo.next == nil {
		juo.next = make(map[int]struct{})
	}
	juo.next[id] = struct{}{}
	return juo
}

// SetNillableNextID sets the next edge to Job by id if the given value is not nil.
func (juo *JobUpdateOne) SetNillableNextID(id *int) *JobUpdateOne {
	if id != nil {
		juo = juo.SetNextID(*id)
	}
	return juo
}

// SetNext sets the next edge to Job.
func (juo *JobUpdateOne) SetNext(j *Job) *JobUpdateOne {
	return juo.SetNextID(j.ID)
}

// SetOwnerID sets the owner edge to User by id.
func (juo *JobUpdateOne) SetOwnerID(id int) *JobUpdateOne {
	if juo.owner == nil {
		juo.owner = make(map[int]struct{})
	}
	juo.owner[id] = struct{}{}
	return juo
}

// SetOwner sets the owner edge to User.
func (juo *JobUpdateOne) SetOwner(u *User) *JobUpdateOne {
	return juo.SetOwnerID(u.ID)
}

// RemoveTaskIDs removes the tasks edge to Task by ids.
func (juo *JobUpdateOne) RemoveTaskIDs(ids ...int) *JobUpdateOne {
	if juo.removedTasks == nil {
		juo.removedTasks = make(map[int]struct{})
	}
	for i := range ids {
		juo.removedTasks[ids[i]] = struct{}{}
	}
	return juo
}

// RemoveTasks removes tasks edges to Task.
func (juo *JobUpdateOne) RemoveTasks(t ...*Task) *JobUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return juo.RemoveTaskIDs(ids...)
}

// RemoveTagIDs removes the tags edge to Tag by ids.
func (juo *JobUpdateOne) RemoveTagIDs(ids ...int) *JobUpdateOne {
	if juo.removedTags == nil {
		juo.removedTags = make(map[int]struct{})
	}
	for i := range ids {
		juo.removedTags[ids[i]] = struct{}{}
	}
	return juo
}

// RemoveTags removes tags edges to Tag.
func (juo *JobUpdateOne) RemoveTags(t ...*Tag) *JobUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return juo.RemoveTagIDs(ids...)
}

// ClearPrev clears the prev edge to Job.
func (juo *JobUpdateOne) ClearPrev() *JobUpdateOne {
	juo.clearedPrev = true
	return juo
}

// ClearNext clears the next edge to Job.
func (juo *JobUpdateOne) ClearNext() *JobUpdateOne {
	juo.clearedNext = true
	return juo
}

// ClearOwner clears the owner edge to User.
func (juo *JobUpdateOne) ClearOwner() *JobUpdateOne {
	juo.clearedOwner = true
	return juo
}

// Save executes the query and returns the updated entity.
func (juo *JobUpdateOne) Save(ctx context.Context) (*Job, error) {
	if juo.Name != nil {
		if err := job.NameValidator(*juo.Name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"Name\": %v", err)
		}
	}
	if juo.Content != nil {
		if err := job.ContentValidator(*juo.Content); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"Content\": %v", err)
		}
	}
	if len(juo.prev) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"prev\"")
	}
	if len(juo.next) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"next\"")
	}
	if len(juo.owner) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"owner\"")
	}
	if juo.clearedOwner && juo.owner == nil {
		return nil, errors.New("ent: clearing a unique edge \"owner\"")
	}
	return juo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (juo *JobUpdateOne) SaveX(ctx context.Context) *Job {
	j, err := juo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return j
}

// Exec executes the query on the entity.
func (juo *JobUpdateOne) Exec(ctx context.Context) error {
	_, err := juo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (juo *JobUpdateOne) ExecX(ctx context.Context) {
	if err := juo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (juo *JobUpdateOne) sqlSave(ctx context.Context) (j *Job, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   job.Table,
			Columns: job.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  juo.id,
				Type:   field.TypeInt,
				Column: job.FieldID,
			},
		},
	}
	if value := juo.Name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: job.FieldName,
		})
	}
	if value := juo.CreationTime; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: job.FieldCreationTime,
		})
	}
	if value := juo.Content; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: job.FieldContent,
		})
	}
	if value := juo.Staged; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: job.FieldStaged,
		})
	}
	if nodes := juo.removedTasks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   job.TasksTable,
			Columns: []string{job.TasksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: task.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := juo.tasks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   job.TasksTable,
			Columns: []string{job.TasksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: task.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := juo.removedTags; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   job.TagsTable,
			Columns: job.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: tag.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := juo.tags; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   job.TagsTable,
			Columns: job.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: tag.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if juo.clearedPrev {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   job.PrevTable,
			Columns: []string{job.PrevColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: job.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := juo.prev; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   job.PrevTable,
			Columns: []string{job.PrevColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: job.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if juo.clearedNext {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   job.NextTable,
			Columns: []string{job.NextColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: job.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := juo.next; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   job.NextTable,
			Columns: []string{job.NextColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: job.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if juo.clearedOwner {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   job.OwnerTable,
			Columns: []string{job.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := juo.owner; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   job.OwnerTable,
			Columns: []string{job.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	j = &Job{config: juo.config}
	_spec.Assign = j.assignValues
	_spec.ScanValues = j.scanValues()
	if err = sqlgraph.UpdateNode(ctx, juo.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return j, nil
}
