package room

import (
	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
)

type (
	RoomOrderInput = ent.RoomOrder
	RoomWhereInput = ent.RoomWhereInput
)

type UpdateRoomInput struct {
	AddVersion    *int64
	LastMessageID *pulid.ID
}

func (i *UpdateRoomInput) Mutate(m *ent.RoomMutation) {
	if v := i.AddVersion; v != nil {
		m.AddVersion(*v)
	}
	if v := i.LastMessageID; v != nil {
		m.SetLastMessageID(*v)
	}
}

type CreateRoomInput struct {
	AddVersion    *int64
	LastMessageID *pulid.ID
}

func (i *CreateRoomInput) Mutate(m *ent.RoomMutation) {
	if v := i.AddVersion; v != nil {
		m.AddVersion(*v)
	}
	if v := i.LastMessageID; v != nil {
		m.SetLastMessageID(*v)
	}
}
