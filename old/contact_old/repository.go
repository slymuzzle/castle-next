package contacts

import (
	"context"
	"journeyhub/ent"
	"journeyhub/ent/user"
)

type Repository interface {
	SetUserPinCode(
		ctx context.Context,
		userID ID,
		pincode ID,
	) error
	CreateContactByPinCode(
		ctx context.Context,
		userID ID,
		pincode ID,
	) (*ent.User, error)
}

type repository struct {
	entClient *ent.Client
}

func NewRepository(entClient *ent.Client) Repository {
	return &repository{
		entClient: entClient,
	}
}

func (r *repository) SetUserPinCode(
	ctx context.Context,
	userID ID,
	pincode ID,
) error {
	return r.entClient.User.
		UpdateOneID(userID).
		SetContactPin(string(pincode)).
		Exec(ctx)
}

func (r *repository) CreateContactByPinCode(
	ctx context.Context,
	userID ID,
	pincode ID,
) (*ent.User, error) {
	pinUser, err := r.entClient.User.
		Query().
		Where(
			user.ContactPin(string(pincode)),
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	err = r.entClient.User.
		UpdateOneID(userID).
		AddContacts(pinUser).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return pinUser, nil
}
