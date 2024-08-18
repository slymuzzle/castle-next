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

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Mixin returns User mixed-in schema.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		pulid.MixinWithPrefix("UR"),
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("first_name").
			Annotations(
				entgql.OrderField("FIRST_NAME"),
			),
		field.String("last_name").
			Annotations(
				entgql.OrderField("LAST_NAME"),
			),
		field.String("nickname").
			Unique().
			Annotations(
				entgql.OrderField("NICKNAME"),
			),
		field.String("email").
			Unique().
			Optional().
			Annotations(
				entgql.OrderField("EMAIL"),
			),
		// field.String("contact_pin").
		// 	Unique(),
		field.String("password").
			Sensitive().
			Annotations(
				entgql.Skip(entgql.SkipAll),
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

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("contacts", User.Type).
			Through("user_contacts", UserContact.Type).
			Annotations(
				entgql.RelayConnection(),
			),
		edge.From("rooms", Room.Type).
			Through("memberships", RoomMember.Type).
			Ref("users").
			Annotations(
				entgql.RelayConnection(),
			),
		edge.To("messages", Message.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entgql.RelayConnection(),
			),
	}
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
	}
}
