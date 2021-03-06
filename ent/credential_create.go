// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/kcarretto/paragon/ent/credential"
	"github.com/kcarretto/paragon/ent/target"
)

// CredentialCreate is the builder for creating a Credential entity.
type CredentialCreate struct {
	config
	principal *string
	secret    *string
	kind      *credential.Kind
	fails     *int
	target    map[int]struct{}
}

// SetPrincipal sets the principal field.
func (cc *CredentialCreate) SetPrincipal(s string) *CredentialCreate {
	cc.principal = &s
	return cc
}

// SetSecret sets the secret field.
func (cc *CredentialCreate) SetSecret(s string) *CredentialCreate {
	cc.secret = &s
	return cc
}

// SetKind sets the kind field.
func (cc *CredentialCreate) SetKind(c credential.Kind) *CredentialCreate {
	cc.kind = &c
	return cc
}

// SetFails sets the fails field.
func (cc *CredentialCreate) SetFails(i int) *CredentialCreate {
	cc.fails = &i
	return cc
}

// SetNillableFails sets the fails field if the given value is not nil.
func (cc *CredentialCreate) SetNillableFails(i *int) *CredentialCreate {
	if i != nil {
		cc.SetFails(*i)
	}
	return cc
}

// SetTargetID sets the target edge to Target by id.
func (cc *CredentialCreate) SetTargetID(id int) *CredentialCreate {
	if cc.target == nil {
		cc.target = make(map[int]struct{})
	}
	cc.target[id] = struct{}{}
	return cc
}

// SetNillableTargetID sets the target edge to Target by id if the given value is not nil.
func (cc *CredentialCreate) SetNillableTargetID(id *int) *CredentialCreate {
	if id != nil {
		cc = cc.SetTargetID(*id)
	}
	return cc
}

// SetTarget sets the target edge to Target.
func (cc *CredentialCreate) SetTarget(t *Target) *CredentialCreate {
	return cc.SetTargetID(t.ID)
}

// Save creates the Credential in the database.
func (cc *CredentialCreate) Save(ctx context.Context) (*Credential, error) {
	if cc.principal == nil {
		return nil, errors.New("ent: missing required field \"principal\"")
	}
	if cc.secret == nil {
		return nil, errors.New("ent: missing required field \"secret\"")
	}
	if err := credential.SecretValidator(*cc.secret); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"secret\": %v", err)
	}
	if cc.kind == nil {
		return nil, errors.New("ent: missing required field \"kind\"")
	}
	if err := credential.KindValidator(*cc.kind); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"kind\": %v", err)
	}
	if cc.fails == nil {
		v := credential.DefaultFails
		cc.fails = &v
	}
	if err := credential.FailsValidator(*cc.fails); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"fails\": %v", err)
	}
	if len(cc.target) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"target\"")
	}
	return cc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (cc *CredentialCreate) SaveX(ctx context.Context) *Credential {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (cc *CredentialCreate) sqlSave(ctx context.Context) (*Credential, error) {
	var (
		c     = &Credential{config: cc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: credential.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: credential.FieldID,
			},
		}
	)
	if value := cc.principal; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: credential.FieldPrincipal,
		})
		c.Principal = *value
	}
	if value := cc.secret; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: credential.FieldSecret,
		})
		c.Secret = *value
	}
	if value := cc.kind; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  *value,
			Column: credential.FieldKind,
		})
		c.Kind = *value
	}
	if value := cc.fails; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: credential.FieldFails,
		})
		c.Fails = *value
	}
	if nodes := cc.target; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   credential.TargetTable,
			Columns: []string{credential.TargetColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: target.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, cc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	c.ID = int(id)
	return c, nil
}
