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

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Mixin returns Room mixed-in schema.
func (Message) Mixin() []ent.Mixin {
	return []ent.Mixin{
		pulid.MixinWithPrefix("ME"),
	}
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.Text("content"),
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

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("voice", MessageVoice.Type).
			Unique(),
		edge.To("reply_to", Message.Type).
			Unique(),
		edge.To("attachments", MessageAttachment.Type),
		edge.To("links", MessageLink.Type),
		edge.From("user", User.Type).
			Ref("messages").
			Unique(),
		edge.From("room", Room.Type).
			Ref("messages").
			Unique(),
	}
}

// Annotations of the Message.
func (Message) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.MultiOrder(),
		entgql.RelayConnection(),
	}
}
