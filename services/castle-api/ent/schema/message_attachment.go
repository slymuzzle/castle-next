package schema

import (
	"journeyhub/ent/schema/pulid"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// File holds the schema definition for the MessageAttachment entity.
type MessageAttachment struct {
	ent.Schema
}

// Mixin returns MessageAttachment mixed-in schema.
func (MessageAttachment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		pulid.MixinWithPrefix("MA"),
	}
}

// Fields of the MessageAttachment.
func (MessageAttachment) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("order").
			Annotations(
				entgql.Type("Uint"),
				entgql.OrderField("ORDER"),
			),
		field.Time("attached_at").
			Immutable().
			Default(time.Now).
			Annotations(
				entgql.OrderField("ATTACHED_AT"),
			),
	}
}

// Edges of the MessageAttachment.
func (MessageAttachment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("room", Room.Type).
			Ref("message_attachments").
			Unique().
			Required(),
		edge.From("message", Message.Type).
			Ref("attachments").
			Unique().
			Required(),
		edge.From("file", File.Type).
			Ref("message_attachment").
			Unique().
			Required(),
	}
}

// Annotations of the MessageAttachment.
func (MessageAttachment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
	}
}
