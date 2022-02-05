// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/kcarretto/paragon/ent/credential"
	"github.com/kcarretto/paragon/ent/tag"
	"github.com/kcarretto/paragon/ent/target"
	"github.com/kcarretto/paragon/ent/task"
)

// TargetCreate is the builder for creating a Target entity.
type TargetCreate struct {
	config
	mutation *TargetMutation
	hooks    []Hook
}

// SetName sets the "Name" field.
func (tc *TargetCreate) SetName(s string) *TargetCreate {
	tc.mutation.SetName(s)
	return tc
}

// SetOS sets the "OS" field.
func (tc *TargetCreate) SetOS(t target.OS) *TargetCreate {
	tc.mutation.SetOS(t)
	return tc
}

// SetPrimaryIP sets the "PrimaryIP" field.
func (tc *TargetCreate) SetPrimaryIP(s string) *TargetCreate {
	tc.mutation.SetPrimaryIP(s)
	return tc
}

// SetMachineUUID sets the "MachineUUID" field.
func (tc *TargetCreate) SetMachineUUID(s string) *TargetCreate {
	tc.mutation.SetMachineUUID(s)
	return tc
}

// SetNillableMachineUUID sets the "MachineUUID" field if the given value is not nil.
func (tc *TargetCreate) SetNillableMachineUUID(s *string) *TargetCreate {
	if s != nil {
		tc.SetMachineUUID(*s)
	}
	return tc
}

// SetPublicIP sets the "PublicIP" field.
func (tc *TargetCreate) SetPublicIP(s string) *TargetCreate {
	tc.mutation.SetPublicIP(s)
	return tc
}

// SetNillablePublicIP sets the "PublicIP" field if the given value is not nil.
func (tc *TargetCreate) SetNillablePublicIP(s *string) *TargetCreate {
	if s != nil {
		tc.SetPublicIP(*s)
	}
	return tc
}

// SetPrimaryMAC sets the "PrimaryMAC" field.
func (tc *TargetCreate) SetPrimaryMAC(s string) *TargetCreate {
	tc.mutation.SetPrimaryMAC(s)
	return tc
}

// SetNillablePrimaryMAC sets the "PrimaryMAC" field if the given value is not nil.
func (tc *TargetCreate) SetNillablePrimaryMAC(s *string) *TargetCreate {
	if s != nil {
		tc.SetPrimaryMAC(*s)
	}
	return tc
}

// SetHostname sets the "Hostname" field.
func (tc *TargetCreate) SetHostname(s string) *TargetCreate {
	tc.mutation.SetHostname(s)
	return tc
}

// SetNillableHostname sets the "Hostname" field if the given value is not nil.
func (tc *TargetCreate) SetNillableHostname(s *string) *TargetCreate {
	if s != nil {
		tc.SetHostname(*s)
	}
	return tc
}

// SetLastSeen sets the "LastSeen" field.
func (tc *TargetCreate) SetLastSeen(t time.Time) *TargetCreate {
	tc.mutation.SetLastSeen(t)
	return tc
}

// SetNillableLastSeen sets the "LastSeen" field if the given value is not nil.
func (tc *TargetCreate) SetNillableLastSeen(t *time.Time) *TargetCreate {
	if t != nil {
		tc.SetLastSeen(*t)
	}
	return tc
}

// AddTaskIDs adds the "tasks" edge to the Task entity by IDs.
func (tc *TargetCreate) AddTaskIDs(ids ...int) *TargetCreate {
	tc.mutation.AddTaskIDs(ids...)
	return tc
}

// AddTasks adds the "tasks" edges to the Task entity.
func (tc *TargetCreate) AddTasks(t ...*Task) *TargetCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tc.AddTaskIDs(ids...)
}

// AddTagIDs adds the "tags" edge to the Tag entity by IDs.
func (tc *TargetCreate) AddTagIDs(ids ...int) *TargetCreate {
	tc.mutation.AddTagIDs(ids...)
	return tc
}

// AddTags adds the "tags" edges to the Tag entity.
func (tc *TargetCreate) AddTags(t ...*Tag) *TargetCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tc.AddTagIDs(ids...)
}

// AddCredentialIDs adds the "credentials" edge to the Credential entity by IDs.
func (tc *TargetCreate) AddCredentialIDs(ids ...int) *TargetCreate {
	tc.mutation.AddCredentialIDs(ids...)
	return tc
}

// AddCredentials adds the "credentials" edges to the Credential entity.
func (tc *TargetCreate) AddCredentials(c ...*Credential) *TargetCreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return tc.AddCredentialIDs(ids...)
}

// Mutation returns the TargetMutation object of the builder.
func (tc *TargetCreate) Mutation() *TargetMutation {
	return tc.mutation
}

// Save creates the Target in the database.
func (tc *TargetCreate) Save(ctx context.Context) (*Target, error) {
	var (
		err  error
		node *Target
	)
	if len(tc.hooks) == 0 {
		if err = tc.check(); err != nil {
			return nil, err
		}
		node, err = tc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TargetMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = tc.check(); err != nil {
				return nil, err
			}
			tc.mutation = mutation
			if node, err = tc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(tc.hooks) - 1; i >= 0; i-- {
			if tc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = tc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, tc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TargetCreate) SaveX(ctx context.Context) *Target {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tc *TargetCreate) Exec(ctx context.Context) error {
	_, err := tc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tc *TargetCreate) ExecX(ctx context.Context) {
	if err := tc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tc *TargetCreate) check() error {
	if _, ok := tc.mutation.Name(); !ok {
		return &ValidationError{Name: "Name", err: errors.New(`ent: missing required field "Target.Name"`)}
	}
	if v, ok := tc.mutation.Name(); ok {
		if err := target.NameValidator(v); err != nil {
			return &ValidationError{Name: "Name", err: fmt.Errorf(`ent: validator failed for field "Target.Name": %w`, err)}
		}
	}
	if _, ok := tc.mutation.OS(); !ok {
		return &ValidationError{Name: "OS", err: errors.New(`ent: missing required field "Target.OS"`)}
	}
	if v, ok := tc.mutation.OS(); ok {
		if err := target.OSValidator(v); err != nil {
			return &ValidationError{Name: "OS", err: fmt.Errorf(`ent: validator failed for field "Target.OS": %w`, err)}
		}
	}
	if _, ok := tc.mutation.PrimaryIP(); !ok {
		return &ValidationError{Name: "PrimaryIP", err: errors.New(`ent: missing required field "Target.PrimaryIP"`)}
	}
	if v, ok := tc.mutation.MachineUUID(); ok {
		if err := target.MachineUUIDValidator(v); err != nil {
			return &ValidationError{Name: "MachineUUID", err: fmt.Errorf(`ent: validator failed for field "Target.MachineUUID": %w`, err)}
		}
	}
	return nil
}

func (tc *TargetCreate) sqlSave(ctx context.Context) (*Target, error) {
	_node, _spec := tc.createSpec()
	if err := sqlgraph.CreateNode(ctx, tc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (tc *TargetCreate) createSpec() (*Target, *sqlgraph.CreateSpec) {
	var (
		_node = &Target{config: tc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: target.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: target.FieldID,
			},
		}
	)
	if value, ok := tc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: target.FieldName,
		})
		_node.Name = value
	}
	if value, ok := tc.mutation.OS(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: target.FieldOS,
		})
		_node.OS = value
	}
	if value, ok := tc.mutation.PrimaryIP(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: target.FieldPrimaryIP,
		})
		_node.PrimaryIP = value
	}
	if value, ok := tc.mutation.MachineUUID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: target.FieldMachineUUID,
		})
		_node.MachineUUID = value
	}
	if value, ok := tc.mutation.PublicIP(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: target.FieldPublicIP,
		})
		_node.PublicIP = value
	}
	if value, ok := tc.mutation.PrimaryMAC(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: target.FieldPrimaryMAC,
		})
		_node.PrimaryMAC = value
	}
	if value, ok := tc.mutation.Hostname(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: target.FieldHostname,
		})
		_node.Hostname = value
	}
	if value, ok := tc.mutation.LastSeen(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: target.FieldLastSeen,
		})
		_node.LastSeen = value
	}
	if nodes := tc.mutation.TasksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   target.TasksTable,
			Columns: []string{target.TasksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: task.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.mutation.TagsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   target.TagsTable,
			Columns: target.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: tag.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.mutation.CredentialsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   target.CredentialsTable,
			Columns: []string{target.CredentialsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: credential.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// TargetCreateBulk is the builder for creating many Target entities in bulk.
type TargetCreateBulk struct {
	config
	builders []*TargetCreate
}

// Save creates the Target entities in the database.
func (tcb *TargetCreateBulk) Save(ctx context.Context) ([]*Target, error) {
	specs := make([]*sqlgraph.CreateSpec, len(tcb.builders))
	nodes := make([]*Target, len(tcb.builders))
	mutators := make([]Mutator, len(tcb.builders))
	for i := range tcb.builders {
		func(i int, root context.Context) {
			builder := tcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TargetMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, tcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, tcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tcb *TargetCreateBulk) SaveX(ctx context.Context) []*Target {
	v, err := tcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcb *TargetCreateBulk) Exec(ctx context.Context) error {
	_, err := tcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcb *TargetCreateBulk) ExecX(ctx context.Context) {
	if err := tcb.Exec(ctx); err != nil {
		panic(err)
	}
}
