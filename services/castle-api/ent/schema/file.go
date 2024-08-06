package schema

import (
	"journeyhub/ent/schema/pulid"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// File holds the schema definition for the File entity.
type File struct {
	ent.Schema
}

// Mixin returns File mixed-in schema.
func (File) Mixin() []ent.Mixin {
	return []ent.Mixin{
		pulid.MixinWithPrefix("FE"),
	}
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("file_name"),
		field.String("mime_type"),
		field.String("disk"),
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
	return nil
}
