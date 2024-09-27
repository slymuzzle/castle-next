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

// Device holds the schema definition for the Device entity.
type Device struct {
	ent.Schema
}

// Mixin returns Device mixed-in schema.
func (Device) Mixin() []ent.Mixin {
	return []ent.Mixin{
		pulid.MixinWithPrefix("DE"),
	}
}

// Fields of the Device.
func (Device) Fields() []ent.Field {
	return []ent.Field{
		field.String("device_id").
			Unique().
			Annotations(
				entgql.OrderField("DEVICE_ID"),
			),
		field.String("fcm_token").
			Unique().
			Annotations(
				entgql.OrderField("FCM_TOKEN"),
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

// Edges of the Device.
func (Device) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("device").
			Unique().
			Required(),
	}
}

// Annotations of the Device.
func (Device) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.MultiOrder(),
		entgql.RelayConnection(),
	}
}
