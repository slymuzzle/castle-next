package schema

import (
	"time"

	"journeyhub/ent/schema/mixin"
	"journeyhub/ent/schema/pulid"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// RoomMember holds the edge schema definition of the Friendship relationship.
type RoomMember struct {
	ent.Schema
}

// Mixin returns RoomMember mixed-in schema.
func (RoomMember) Mixin() []ent.Mixin {
	return []ent.Mixin{
		pulid.MixinWithPrefix("RM"),
		mixin.SoftDeleteMixin{},
	}
}

// Fields of the RoomMember.
func (RoomMember) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Optional().
			Annotations(
				entgql.OrderField("NAME"),
			),
		field.Int("unread_messages_count").
			Default(0).
			Annotations(
				entgql.OrderField("UNREAD_MESSAGES_COUNT"),
			),
		field.String("user_id").
			GoType(pulid.ID("")),
		field.String("room_id").
			GoType(pulid.ID("")),
		field.Time("joined_at").
			Immutable().
			Default(time.Now).
			Annotations(
				entgql.OrderField("JOINED_AT"),
			),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(
				entgql.OrderField("UPDATED_AT"),
			),
	}
}

// Edges of the RoomMember.
func (RoomMember) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Required().
			Unique().
			Field("user_id").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
		edge.To("room", Room.Type).
			Required().
			Unique().
			Field("room_id").
			Annotations(
				entgql.OrderField("ROOM_UPDATED_AT"),
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Indexes of the RoomMember.
func (RoomMember) Indexes() []ent.Index {
	return []ent.Index{}
}

// Annotations of the RoomMember.
func (RoomMember) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.MultiOrder(),
		entgql.RelayConnection(),
	}
}
