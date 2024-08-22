package user

import (
	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
)

type (
	UserWhereInput = ent.UserWhereInput
	UserOrderInput = ent.UserOrder
)

type UpdateUserInput struct {
	PinCode *pulid.ID
}

func (c *UpdateUserInput) Mutate(m *ent.UserMutation) {
}

type CreateUserInput struct {
	AddVersion    *int64
	LastMessageID *pulid.ID
}

func (c *CreateUserInput) Mutate(m *ent.UserMutation) {
}
