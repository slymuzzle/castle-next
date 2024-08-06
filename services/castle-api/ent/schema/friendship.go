package schema

import (
	"journeyhub/ent/schema/pulid"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Friendship holds the edge schema definition of the Friendship relationship.
type Friendship struct {
	ent.Schema
}

// Mixin returns Friendship mixed-in schema.
func (Friendship) Mixin() []ent.Mixin {
	return []ent.Mixin{
		pulid.MixinWithPrefix("FS"),
	}
}

// Fields of the Friendship.
func (Friendship) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id").
			GoType(pulid.ID("")),
		field.String("friend_id").
			GoType(pulid.ID("")),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Friendship.
func (Friendship) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Required().
			Unique().
			Field("user_id"),
		edge.To("friend", User.Type).
			Required().
			Unique().
			Field("friend_id"),
	}
}

// Indexes of the Friendship.
func (Friendship) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "friend_id").
			Unique(),
	}
}
