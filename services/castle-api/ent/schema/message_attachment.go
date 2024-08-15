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
		field.Enum("type").
			Values("Image", "Video", "Document").
			Annotations(
				entgql.OrderField("TYPE"),
			),
		field.Uint64("order").
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
		edge.From("message", Message.Type).
			Ref("attachments").
			Required().
			Unique(),
		edge.To("file", File.Type).
			Required().
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}
