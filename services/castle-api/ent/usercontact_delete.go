// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"journeyhub/ent/predicate"
	"journeyhub/ent/usercontact"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// UserContactDelete is the builder for deleting a UserContact entity.
type UserContactDelete struct {
	config
	hooks    []Hook
	mutation *UserContactMutation
}

// Where appends a list predicates to the UserContactDelete builder.
func (ucd *UserContactDelete) Where(ps ...predicate.UserContact) *UserContactDelete {
	ucd.mutation.Where(ps...)
	return ucd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ucd *UserContactDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, ucd.sqlExec, ucd.mutation, ucd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (ucd *UserContactDelete) ExecX(ctx context.Context) int {
	n, err := ucd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ucd *UserContactDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(usercontact.Table, sqlgraph.NewFieldSpec(usercontact.FieldID, field.TypeString))
	if ps := ucd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, ucd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	ucd.mutation.done = true
	return affected, err
}

// UserContactDeleteOne is the builder for deleting a single UserContact entity.
type UserContactDeleteOne struct {
	ucd *UserContactDelete
}

// Where appends a list predicates to the UserContactDelete builder.
func (ucdo *UserContactDeleteOne) Where(ps ...predicate.UserContact) *UserContactDeleteOne {
	ucdo.ucd.mutation.Where(ps...)
	return ucdo
}

// Exec executes the deletion query.
func (ucdo *UserContactDeleteOne) Exec(ctx context.Context) error {
	n, err := ucdo.ucd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{usercontact.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ucdo *UserContactDeleteOne) ExecX(ctx context.Context) {
	if err := ucdo.Exec(ctx); err != nil {
		panic(err)
	}
}
