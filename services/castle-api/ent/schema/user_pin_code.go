package schema

import (
	"journeyhub/ent/schema/pulid"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// UserPinCode holds the edge schema definition of the Friendship relationship.
type UserPinCode struct {
	ent.Schema
}

// Mixin returns UserPinCode mixed-in schema.
func (UserPinCode) Mixin() []ent.Mixin {
	return []ent.Mixin{
		pulid.MixinWithPrefix("UC"),
	}
}

// Fields of the UserPinCode.
func (UserPinCode) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id").
			GoType(pulid.ID("")),
		field.String("contact_id").
			GoType(pulid.ID("")),
		field.String("room_id").
			Optional().
			GoType(pulid.ID("")),
		field.Time("created_at").
			Default(time.Now).
			Annotations(
				entgql.OrderField("CREATED_AT"),
			),
	}
}

// Edges of the UserPinCode.
func (UserPinCode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Required().
			Unique().
			Field("user_id"),
		edge.To("contact", User.Type).
			Required().
			Unique().
			Field("contact_id"),
		edge.To("room", Room.Type).
			Unique().
			Field("room_id"),
	}
}
