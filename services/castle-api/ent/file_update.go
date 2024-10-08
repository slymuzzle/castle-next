// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"journeyhub/ent/file"
	"journeyhub/ent/messageattachment"
	"journeyhub/ent/messagevoice"
	"journeyhub/ent/predicate"
	"journeyhub/ent/schema/pulid"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// FileUpdate is the builder for updating File entities.
type FileUpdate struct {
	config
	hooks    []Hook
	mutation *FileMutation
}

// Where appends a list predicates to the FileUpdate builder.
func (fu *FileUpdate) Where(ps ...predicate.File) *FileUpdate {
	fu.mutation.Where(ps...)
	return fu
}

// SetName sets the "name" field.
func (fu *FileUpdate) SetName(s string) *FileUpdate {
	fu.mutation.SetName(s)
	return fu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (fu *FileUpdate) SetNillableName(s *string) *FileUpdate {
	if s != nil {
		fu.SetName(*s)
	}
	return fu
}

// SetContentType sets the "content_type" field.
func (fu *FileUpdate) SetContentType(s string) *FileUpdate {
	fu.mutation.SetContentType(s)
	return fu
}

// SetNillableContentType sets the "content_type" field if the given value is not nil.
func (fu *FileUpdate) SetNillableContentType(s *string) *FileUpdate {
	if s != nil {
		fu.SetContentType(*s)
	}
	return fu
}

// SetSize sets the "size" field.
func (fu *FileUpdate) SetSize(u uint64) *FileUpdate {
	fu.mutation.ResetSize()
	fu.mutation.SetSize(u)
	return fu
}

// SetNillableSize sets the "size" field if the given value is not nil.
func (fu *FileUpdate) SetNillableSize(u *uint64) *FileUpdate {
	if u != nil {
		fu.SetSize(*u)
	}
	return fu
}

// AddSize adds u to the "size" field.
func (fu *FileUpdate) AddSize(u int64) *FileUpdate {
	fu.mutation.AddSize(u)
	return fu
}

// SetLocation sets the "location" field.
func (fu *FileUpdate) SetLocation(s string) *FileUpdate {
	fu.mutation.SetLocation(s)
	return fu
}

// SetNillableLocation sets the "location" field if the given value is not nil.
func (fu *FileUpdate) SetNillableLocation(s *string) *FileUpdate {
	if s != nil {
		fu.SetLocation(*s)
	}
	return fu
}

// ClearLocation clears the value of the "location" field.
func (fu *FileUpdate) ClearLocation() *FileUpdate {
	fu.mutation.ClearLocation()
	return fu
}

// SetBucket sets the "bucket" field.
func (fu *FileUpdate) SetBucket(s string) *FileUpdate {
	fu.mutation.SetBucket(s)
	return fu
}

// SetNillableBucket sets the "bucket" field if the given value is not nil.
func (fu *FileUpdate) SetNillableBucket(s *string) *FileUpdate {
	if s != nil {
		fu.SetBucket(*s)
	}
	return fu
}

// SetPath sets the "path" field.
func (fu *FileUpdate) SetPath(s string) *FileUpdate {
	fu.mutation.SetPath(s)
	return fu
}

// SetNillablePath sets the "path" field if the given value is not nil.
func (fu *FileUpdate) SetNillablePath(s *string) *FileUpdate {
	if s != nil {
		fu.SetPath(*s)
	}
	return fu
}

// SetUpdatedAt sets the "updated_at" field.
func (fu *FileUpdate) SetUpdatedAt(t time.Time) *FileUpdate {
	fu.mutation.SetUpdatedAt(t)
	return fu
}

// SetMessageAttachmentID sets the "message_attachment" edge to the MessageAttachment entity by ID.
func (fu *FileUpdate) SetMessageAttachmentID(id pulid.ID) *FileUpdate {
	fu.mutation.SetMessageAttachmentID(id)
	return fu
}

// SetNillableMessageAttachmentID sets the "message_attachment" edge to the MessageAttachment entity by ID if the given value is not nil.
func (fu *FileUpdate) SetNillableMessageAttachmentID(id *pulid.ID) *FileUpdate {
	if id != nil {
		fu = fu.SetMessageAttachmentID(*id)
	}
	return fu
}

// SetMessageAttachment sets the "message_attachment" edge to the MessageAttachment entity.
func (fu *FileUpdate) SetMessageAttachment(m *MessageAttachment) *FileUpdate {
	return fu.SetMessageAttachmentID(m.ID)
}

// SetMessageVoiceID sets the "message_voice" edge to the MessageVoice entity by ID.
func (fu *FileUpdate) SetMessageVoiceID(id pulid.ID) *FileUpdate {
	fu.mutation.SetMessageVoiceID(id)
	return fu
}

// SetNillableMessageVoiceID sets the "message_voice" edge to the MessageVoice entity by ID if the given value is not nil.
func (fu *FileUpdate) SetNillableMessageVoiceID(id *pulid.ID) *FileUpdate {
	if id != nil {
		fu = fu.SetMessageVoiceID(*id)
	}
	return fu
}

// SetMessageVoice sets the "message_voice" edge to the MessageVoice entity.
func (fu *FileUpdate) SetMessageVoice(m *MessageVoice) *FileUpdate {
	return fu.SetMessageVoiceID(m.ID)
}

// Mutation returns the FileMutation object of the builder.
func (fu *FileUpdate) Mutation() *FileMutation {
	return fu.mutation
}

// ClearMessageAttachment clears the "message_attachment" edge to the MessageAttachment entity.
func (fu *FileUpdate) ClearMessageAttachment() *FileUpdate {
	fu.mutation.ClearMessageAttachment()
	return fu
}

// ClearMessageVoice clears the "message_voice" edge to the MessageVoice entity.
func (fu *FileUpdate) ClearMessageVoice() *FileUpdate {
	fu.mutation.ClearMessageVoice()
	return fu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (fu *FileUpdate) Save(ctx context.Context) (int, error) {
	fu.defaults()
	return withHooks(ctx, fu.sqlSave, fu.mutation, fu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (fu *FileUpdate) SaveX(ctx context.Context) int {
	affected, err := fu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fu *FileUpdate) Exec(ctx context.Context) error {
	_, err := fu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fu *FileUpdate) ExecX(ctx context.Context) {
	if err := fu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (fu *FileUpdate) defaults() {
	if _, ok := fu.mutation.UpdatedAt(); !ok {
		v := file.UpdateDefaultUpdatedAt()
		fu.mutation.SetUpdatedAt(v)
	}
}

func (fu *FileUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(file.Table, file.Columns, sqlgraph.NewFieldSpec(file.FieldID, field.TypeString))
	if ps := fu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fu.mutation.Name(); ok {
		_spec.SetField(file.FieldName, field.TypeString, value)
	}
	if value, ok := fu.mutation.ContentType(); ok {
		_spec.SetField(file.FieldContentType, field.TypeString, value)
	}
	if value, ok := fu.mutation.Size(); ok {
		_spec.SetField(file.FieldSize, field.TypeUint64, value)
	}
	if value, ok := fu.mutation.AddedSize(); ok {
		_spec.AddField(file.FieldSize, field.TypeUint64, value)
	}
	if value, ok := fu.mutation.Location(); ok {
		_spec.SetField(file.FieldLocation, field.TypeString, value)
	}
	if fu.mutation.LocationCleared() {
		_spec.ClearField(file.FieldLocation, field.TypeString)
	}
	if value, ok := fu.mutation.Bucket(); ok {
		_spec.SetField(file.FieldBucket, field.TypeString, value)
	}
	if value, ok := fu.mutation.Path(); ok {
		_spec.SetField(file.FieldPath, field.TypeString, value)
	}
	if value, ok := fu.mutation.UpdatedAt(); ok {
		_spec.SetField(file.FieldUpdatedAt, field.TypeTime, value)
	}
	if fu.mutation.MessageAttachmentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   file.MessageAttachmentTable,
			Columns: []string{file.MessageAttachmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(messageattachment.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.MessageAttachmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   file.MessageAttachmentTable,
			Columns: []string{file.MessageAttachmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(messageattachment.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fu.mutation.MessageVoiceCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   file.MessageVoiceTable,
			Columns: []string{file.MessageVoiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(messagevoice.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.MessageVoiceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   file.MessageVoiceTable,
			Columns: []string{file.MessageVoiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(messagevoice.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{file.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	fu.mutation.done = true
	return n, nil
}

// FileUpdateOne is the builder for updating a single File entity.
type FileUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *FileMutation
}

// SetName sets the "name" field.
func (fuo *FileUpdateOne) SetName(s string) *FileUpdateOne {
	fuo.mutation.SetName(s)
	return fuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableName(s *string) *FileUpdateOne {
	if s != nil {
		fuo.SetName(*s)
	}
	return fuo
}

// SetContentType sets the "content_type" field.
func (fuo *FileUpdateOne) SetContentType(s string) *FileUpdateOne {
	fuo.mutation.SetContentType(s)
	return fuo
}

// SetNillableContentType sets the "content_type" field if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableContentType(s *string) *FileUpdateOne {
	if s != nil {
		fuo.SetContentType(*s)
	}
	return fuo
}

// SetSize sets the "size" field.
func (fuo *FileUpdateOne) SetSize(u uint64) *FileUpdateOne {
	fuo.mutation.ResetSize()
	fuo.mutation.SetSize(u)
	return fuo
}

// SetNillableSize sets the "size" field if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableSize(u *uint64) *FileUpdateOne {
	if u != nil {
		fuo.SetSize(*u)
	}
	return fuo
}

// AddSize adds u to the "size" field.
func (fuo *FileUpdateOne) AddSize(u int64) *FileUpdateOne {
	fuo.mutation.AddSize(u)
	return fuo
}

// SetLocation sets the "location" field.
func (fuo *FileUpdateOne) SetLocation(s string) *FileUpdateOne {
	fuo.mutation.SetLocation(s)
	return fuo
}

// SetNillableLocation sets the "location" field if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableLocation(s *string) *FileUpdateOne {
	if s != nil {
		fuo.SetLocation(*s)
	}
	return fuo
}

// ClearLocation clears the value of the "location" field.
func (fuo *FileUpdateOne) ClearLocation() *FileUpdateOne {
	fuo.mutation.ClearLocation()
	return fuo
}

// SetBucket sets the "bucket" field.
func (fuo *FileUpdateOne) SetBucket(s string) *FileUpdateOne {
	fuo.mutation.SetBucket(s)
	return fuo
}

// SetNillableBucket sets the "bucket" field if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableBucket(s *string) *FileUpdateOne {
	if s != nil {
		fuo.SetBucket(*s)
	}
	return fuo
}

// SetPath sets the "path" field.
func (fuo *FileUpdateOne) SetPath(s string) *FileUpdateOne {
	fuo.mutation.SetPath(s)
	return fuo
}

// SetNillablePath sets the "path" field if the given value is not nil.
func (fuo *FileUpdateOne) SetNillablePath(s *string) *FileUpdateOne {
	if s != nil {
		fuo.SetPath(*s)
	}
	return fuo
}

// SetUpdatedAt sets the "updated_at" field.
func (fuo *FileUpdateOne) SetUpdatedAt(t time.Time) *FileUpdateOne {
	fuo.mutation.SetUpdatedAt(t)
	return fuo
}

// SetMessageAttachmentID sets the "message_attachment" edge to the MessageAttachment entity by ID.
func (fuo *FileUpdateOne) SetMessageAttachmentID(id pulid.ID) *FileUpdateOne {
	fuo.mutation.SetMessageAttachmentID(id)
	return fuo
}

// SetNillableMessageAttachmentID sets the "message_attachment" edge to the MessageAttachment entity by ID if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableMessageAttachmentID(id *pulid.ID) *FileUpdateOne {
	if id != nil {
		fuo = fuo.SetMessageAttachmentID(*id)
	}
	return fuo
}

// SetMessageAttachment sets the "message_attachment" edge to the MessageAttachment entity.
func (fuo *FileUpdateOne) SetMessageAttachment(m *MessageAttachment) *FileUpdateOne {
	return fuo.SetMessageAttachmentID(m.ID)
}

// SetMessageVoiceID sets the "message_voice" edge to the MessageVoice entity by ID.
func (fuo *FileUpdateOne) SetMessageVoiceID(id pulid.ID) *FileUpdateOne {
	fuo.mutation.SetMessageVoiceID(id)
	return fuo
}

// SetNillableMessageVoiceID sets the "message_voice" edge to the MessageVoice entity by ID if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableMessageVoiceID(id *pulid.ID) *FileUpdateOne {
	if id != nil {
		fuo = fuo.SetMessageVoiceID(*id)
	}
	return fuo
}

// SetMessageVoice sets the "message_voice" edge to the MessageVoice entity.
func (fuo *FileUpdateOne) SetMessageVoice(m *MessageVoice) *FileUpdateOne {
	return fuo.SetMessageVoiceID(m.ID)
}

// Mutation returns the FileMutation object of the builder.
func (fuo *FileUpdateOne) Mutation() *FileMutation {
	return fuo.mutation
}

// ClearMessageAttachment clears the "message_attachment" edge to the MessageAttachment entity.
func (fuo *FileUpdateOne) ClearMessageAttachment() *FileUpdateOne {
	fuo.mutation.ClearMessageAttachment()
	return fuo
}

// ClearMessageVoice clears the "message_voice" edge to the MessageVoice entity.
func (fuo *FileUpdateOne) ClearMessageVoice() *FileUpdateOne {
	fuo.mutation.ClearMessageVoice()
	return fuo
}

// Where appends a list predicates to the FileUpdate builder.
func (fuo *FileUpdateOne) Where(ps ...predicate.File) *FileUpdateOne {
	fuo.mutation.Where(ps...)
	return fuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (fuo *FileUpdateOne) Select(field string, fields ...string) *FileUpdateOne {
	fuo.fields = append([]string{field}, fields...)
	return fuo
}

// Save executes the query and returns the updated File entity.
func (fuo *FileUpdateOne) Save(ctx context.Context) (*File, error) {
	fuo.defaults()
	return withHooks(ctx, fuo.sqlSave, fuo.mutation, fuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (fuo *FileUpdateOne) SaveX(ctx context.Context) *File {
	node, err := fuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (fuo *FileUpdateOne) Exec(ctx context.Context) error {
	_, err := fuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fuo *FileUpdateOne) ExecX(ctx context.Context) {
	if err := fuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (fuo *FileUpdateOne) defaults() {
	if _, ok := fuo.mutation.UpdatedAt(); !ok {
		v := file.UpdateDefaultUpdatedAt()
		fuo.mutation.SetUpdatedAt(v)
	}
}

func (fuo *FileUpdateOne) sqlSave(ctx context.Context) (_node *File, err error) {
	_spec := sqlgraph.NewUpdateSpec(file.Table, file.Columns, sqlgraph.NewFieldSpec(file.FieldID, field.TypeString))
	id, ok := fuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "File.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := fuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, file.FieldID)
		for _, f := range fields {
			if !file.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != file.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := fuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fuo.mutation.Name(); ok {
		_spec.SetField(file.FieldName, field.TypeString, value)
	}
	if value, ok := fuo.mutation.ContentType(); ok {
		_spec.SetField(file.FieldContentType, field.TypeString, value)
	}
	if value, ok := fuo.mutation.Size(); ok {
		_spec.SetField(file.FieldSize, field.TypeUint64, value)
	}
	if value, ok := fuo.mutation.AddedSize(); ok {
		_spec.AddField(file.FieldSize, field.TypeUint64, value)
	}
	if value, ok := fuo.mutation.Location(); ok {
		_spec.SetField(file.FieldLocation, field.TypeString, value)
	}
	if fuo.mutation.LocationCleared() {
		_spec.ClearField(file.FieldLocation, field.TypeString)
	}
	if value, ok := fuo.mutation.Bucket(); ok {
		_spec.SetField(file.FieldBucket, field.TypeString, value)
	}
	if value, ok := fuo.mutation.Path(); ok {
		_spec.SetField(file.FieldPath, field.TypeString, value)
	}
	if value, ok := fuo.mutation.UpdatedAt(); ok {
		_spec.SetField(file.FieldUpdatedAt, field.TypeTime, value)
	}
	if fuo.mutation.MessageAttachmentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   file.MessageAttachmentTable,
			Columns: []string{file.MessageAttachmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(messageattachment.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.MessageAttachmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   file.MessageAttachmentTable,
			Columns: []string{file.MessageAttachmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(messageattachment.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fuo.mutation.MessageVoiceCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   file.MessageVoiceTable,
			Columns: []string{file.MessageVoiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(messagevoice.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.MessageVoiceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   file.MessageVoiceTable,
			Columns: []string{file.MessageVoiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(messagevoice.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &File{config: fuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, fuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{file.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	fuo.mutation.done = true
	return _node, nil
}
