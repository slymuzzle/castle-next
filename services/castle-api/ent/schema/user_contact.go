package schema

import (
	"journeyhub/ent/schema/pulid"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// UserContact holds the edge schema definition of the Friendship relationship.
type UserContact struct {
	ent.Schema
}

// Mixin returns UserContacts mixed-in schema.
func (UserContact) Mixin() []ent.Mixin {
	return []ent.Mixin{
		pulid.MixinWithPrefix("UC"),
	}
}

// Fields of the UserContacts.
func (UserContact) Fields() []ent.Field {
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

// Edges of the UserContacts.
func (UserContact) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Required().
			Unique().
			Field("user_id").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
		edge.To("contact", User.Type).
			Required().
			Unique().
			Field("contact_id").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
		edge.To("room", Room.Type).
			Unique().
			Field("room_id").
			Annotations(
				entsql.OnDelete(entsql.SetNull),
			),
	}
}

// Indexes of the UserContact.
func (UserContact) Indexes() []ent.Index {
	return []ent.Index{}
}

// Annotations of the UserContact.
func (UserContact) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.MultiOrder(),
		entgql.RelayConnection(),
	}
}
