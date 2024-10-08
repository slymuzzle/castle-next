// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"journeyhub/ent/messagevoice"
	"journeyhub/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// MessageVoiceDelete is the builder for deleting a MessageVoice entity.
type MessageVoiceDelete struct {
	config
	hooks    []Hook
	mutation *MessageVoiceMutation
}

// Where appends a list predicates to the MessageVoiceDelete builder.
func (mvd *MessageVoiceDelete) Where(ps ...predicate.MessageVoice) *MessageVoiceDelete {
	mvd.mutation.Where(ps...)
	return mvd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (mvd *MessageVoiceDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, mvd.sqlExec, mvd.mutation, mvd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (mvd *MessageVoiceDelete) ExecX(ctx context.Context) int {
	n, err := mvd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (mvd *MessageVoiceDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(messagevoice.Table, sqlgraph.NewFieldSpec(messagevoice.FieldID, field.TypeString))
	if ps := mvd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, mvd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	mvd.mutation.done = true
	return affected, err
}

// MessageVoiceDeleteOne is the builder for deleting a single MessageVoice entity.
type MessageVoiceDeleteOne struct {
	mvd *MessageVoiceDelete
}

// Where appends a list predicates to the MessageVoiceDelete builder.
func (mvdo *MessageVoiceDeleteOne) Where(ps ...predicate.MessageVoice) *MessageVoiceDeleteOne {
	mvdo.mvd.mutation.Where(ps...)
	return mvdo
}

// Exec executes the deletion query.
func (mvdo *MessageVoiceDeleteOne) Exec(ctx context.Context) error {
	n, err := mvdo.mvd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{messagevoice.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (mvdo *MessageVoiceDeleteOne) ExecX(ctx context.Context) {
	if err := mvdo.Exec(ctx); err != nil {
		panic(err)
	}
}
