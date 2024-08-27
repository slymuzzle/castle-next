package contacts

import (
	"context"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/ent/user"
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

type Repository interface {
	UpdateUserPinCode(
		ctx context.Context,
		userID pulid.ID,
		pincode string,
	) (*ent.User, error)

	FindUserByPinCode(
		ctx context.Context,
		pincode string,
	) (*ent.User, error)

	Create(
		ctx context.Context,
		input CreateUserContactInput,
	) (*ent.UserContact, error)

	Update(
		ctx context.Context,
		ID pulid.ID,
		input UpdateUserContactInput,
	) (*ent.UserContact, error)

	Delete(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.UserContact, error)
}

type repository struct {
	entClient *ent.Client
}

func NewRepository(entClient *ent.Client) Repository {
	return &repository{
		entClient: entClient,
	}
}

func (r *repository) FindUserByPinCode(
	ctx context.Context,
	pincode string,
) (*ent.User, error) {
	return r.entClient.User.
		Query().
		Where(
			user.ContactPin(pincode),
		).
		Only(ctx)
}

func (r *repository) UpdateUserPinCode(
	ctx context.Context,
	userID pulid.ID,
	pincode string,
) (*ent.User, error) {
	return r.entClient.User.
		UpdateOneID(userID).
		SetContactPin(pincode).
		Save(ctx)
}

func (r *repository) Create(
	ctx context.Context,
	input CreateUserContactInput,
) (*ent.UserContact, error) {
	u := r.entClient.UserContact.Create()
	input.Mutate(u.Mutation())
	return u.Save(ctx)
}

func (r *repository) Update(
	ctx context.Context,
	ID pulid.ID,
	input UpdateUserContactInput,
) (*ent.UserContact, error) {
	u := r.entClient.UserContact.UpdateOneID(ID)
	input.Mutate(u.Mutation())
	return u.Save(ctx)
}

func (r *repository) Delete(
	ctx context.Context,
	ID pulid.ID,
) (*ent.UserContact, error) {
	uc, err := r.entClient.UserContact.Get(ctx, ID)
	if err != nil {
		return nil, err
	}

	err = r.entClient.UserContact.
		DeleteOneID(ID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return uc, nil
}
