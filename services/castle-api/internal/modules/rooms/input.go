package rooms

import (
	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
)

type UpdateRoomInput struct {
	AddVersion    int64
	LastMessageID *pulid.ID
}

func (i *UpdateRoomInput) Mutate(m *ent.RoomMutation) {
	m.AddVersion(i.AddVersion)
	if v := i.LastMessageID; v != nil {
		m.SetLastMessageID(*v)
	}
}

type CreateRoomInput struct {
	LastMessageID *pulid.ID
}

func (i *CreateRoomInput) Mutate(m *ent.RoomMutation) {
	if v := i.LastMessageID; v != nil {
		m.SetLastMessageID(*v)
	}
}
