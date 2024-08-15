package schema

import (
	"journeyhub/ent/schema/pulid"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// File holds the schema definition for the File entity.
type File struct {
	ent.Schema
}

// Mixin returns File mixed-in schema.
func (File) Mixin() []ent.Mixin {
	return []ent.Mixin{
		pulid.MixinWithPrefix("FI"),
	}
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Annotations(
				entgql.OrderField("NAME"),
			),
		field.String("mime_type").
			Annotations(
				entgql.OrderField("MIME_TYPE"),
			),
		field.String("disk").
			Annotations(
				entgql.OrderField("DISK"),
			),
		field.Uint64("size").
			Annotations(
				entgql.Type("Uint64"),
				entgql.OrderField("SIZE"),
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

// Edges of the File.
func (File) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("message_attachment", MessageAttachment.Type).
			Ref("file").
			Unique(),
		edge.From("message_voice", MessageVoice.Type).
			Ref("file").
			Unique(),
	}
}
