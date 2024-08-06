package schema

import (
	"journeyhub/ent/schema/pulid"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// RoomMember holds the edge schema definition of the Friendship relationship.
type RoomMember struct {
	ent.Schema
}

// Mixin returns RoomMember mixed-in schema.
func (RoomMember) Mixin() []ent.Mixin {
	return []ent.Mixin{
		pulid.MixinWithPrefix("RM"),
	}
}

// Fields of the RoomMember.
func (RoomMember) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id").
			GoType(pulid.ID("")),
		field.String("room_id").
			GoType(pulid.ID("")),
		field.Time("joined_at").
			Immutable().
			Default(time.Now),
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
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Indexes of the Friendship.
func (RoomMember) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "room_id").
			Unique(),
	}
}
