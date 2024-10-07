package schema

import (
	"time"

	"journeyhub/ent/schema/pulid"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// File holds the schema definition for the MessageVoice entity.
type MessageVoice struct {
	ent.Schema
}

// Mixin returns MessageVoice mixed-in schema.
func (MessageVoice) Mixin() []ent.Mixin {
	return []ent.Mixin{
		pulid.MixinWithPrefix("MV"),
	}
}

// Fields of the MessageAttachment.
func (MessageVoice) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("length").
			Annotations(
				entgql.Type("Uint64"),
				entgql.OrderField("LENGTH"),
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
func (MessageVoice) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("room", Room.Type).
			Ref("message_voices").
			Unique().
			Required(),
		edge.From("message", Message.Type).
			Ref("voice").
			Unique().
			Required(),
		edge.From("file", File.Type).
			Ref("message_voice").
			Unique().
			Required(),
	}
}

// Annotations of the MessageVoice.
func (MessageVoice) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.MultiOrder(),
		entgql.RelayConnection(),
	}
}
