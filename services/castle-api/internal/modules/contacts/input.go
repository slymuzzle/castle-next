package contacts

import (
	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
)

type CreateUserContactInput struct {
	UserID    pulid.ID
	ContactID pulid.ID
	RoomID    *pulid.ID
}

func (c *CreateUserContactInput) Mutate(m *ent.UserContactMutation) {
	m.SetUserID(c.UserID)
	m.SetContactID(c.ContactID)

	if v := c.RoomID; v != nil {
		m.SetRoomID(*v)
	}
}

type UpdateUserContactInput struct {
	UserID    *pulid.ID
	ContactID *pulid.ID
	RoomID    *pulid.ID
}

func (c *UpdateUserContactInput) Mutate(m *ent.UserContactMutation) {
	if v := c.UserID; v != nil {
		m.SetUserID(*v)
	}
	if v := c.ContactID; v != nil {
		m.SetContactID(*v)
	}
	if v := c.RoomID; v != nil {
		m.SetRoomID(*v)
	}
}
