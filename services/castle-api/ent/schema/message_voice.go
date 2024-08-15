package schema

import (
	"journeyhub/ent/schema/pulid"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
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
		field.Int("length").
			Annotations(
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
		edge.To("file", File.Type).
			Required().
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}
