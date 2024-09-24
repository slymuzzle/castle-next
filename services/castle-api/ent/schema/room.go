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

// Room holds the schema definition for the Room entity.
type Room struct {
	ent.Schema
}

// Mixin returns Room mixed-in schema.
func (Room) Mixin() []ent.Mixin {
	return []ent.Mixin{
		pulid.MixinWithPrefix("RO"),
		mixin.SoftDeleteMixin{},
	}
}

// Fields of the Room.
func (Room) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Optional().
			Annotations(
				entgql.OrderField("NAME"),
			),
		field.String("description").
			Optional().
			Annotations(
				entgql.OrderField("DESCRIPTION"),
			),
		field.Uint64("version").
			Positive().
			Default(1).
			Annotations(
				entgql.Type("Uint64"),
				entgql.OrderField("VERSION"),
			),
		field.Enum("type").
			Values("Personal", "Group").
			Annotations(
				entgql.OrderField("TYPE"),
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

// Edges of the Room.
func (Room) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user_contacts", UserContact.Type).
			Ref("room"),
		edge.To("users", User.Type).
			Through("room_members", RoomMember.Type).
			Annotations(
				entgql.RelayConnection(),
			),
		edge.To("last_message", Message.Type).
			Unique().
			Annotations(
				entgql.OrderField("LAST_MESSAGE_CREATED_AT"),
			),
		edge.To("messages", Message.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entgql.RelayConnection(),
			),
		edge.To("message_voices", MessageVoice.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entgql.RelayConnection(),
			),
		edge.To("message_attachments", MessageAttachment.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entgql.RelayConnection(),
			),
		edge.To("message_links", MessageLink.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entgql.RelayConnection(),
			),
	}
}

// Annotations of the Room.
func (Room) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.MultiOrder(),
		entgql.RelayConnection(),
	}
}
