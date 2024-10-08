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

// File holds the schema definition for the MessageLink entity.
type MessageLink struct {
	ent.Schema
}

// Mixin returns MessageLink mixed-in schema.
func (MessageLink) Mixin() []ent.Mixin {
	return []ent.Mixin{
		pulid.MixinWithPrefix("ML"),
	}
}

// Fields of the MessageLink.
func (MessageLink) Fields() []ent.Field {
	return []ent.Field{
		field.String("link").
			Annotations(
				entgql.OrderField("LINK"),
			),
		field.String("title").
			Optional().
			Annotations(
				entgql.OrderField("TITLE"),
			),
		field.String("description").
			Optional().
			Annotations(
				entgql.OrderField("DESCRIPTION"),
			),
		field.String("image_url").
			Optional().
			Annotations(
				entgql.OrderField("IMAGE_URL"),
			),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Annotations(
				entgql.OrderField("CREATED_AT"),
			),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(
				entgql.OrderField("UPDATED_AT"),
			),
	}
}

// Edges of the MessageLink.
func (MessageLink) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("room", Room.Type).
			Ref("message_links").
			Unique().
			Required(),
		edge.From("message", Message.Type).
			Ref("links").
			Required().
			Unique(),
	}
}

// Annotations of the MessageLink.
func (MessageLink) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.MultiOrder(),
		entgql.RelayConnection(),
	}
}
