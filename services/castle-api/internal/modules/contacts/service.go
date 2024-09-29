package contacts

import (
	"context"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/ent/user"
	"journeyhub/internal/modules/auth"
)

type Service interface {
	GenerateUserPinCode(
		ctx context.Context,
	) (*string, error)

	AddUserContact(
		ctx context.Context,
		pincode string,
	) (*ent.UserContact, error)

	DeleteUserContact(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.UserContact, error)
}

type service struct {
	entClient   *ent.Client
	authService auth.Service
}

func NewService(entClient *ent.Client, authService auth.Service) Service {
	return &service{
		entClient:   entClient,
		authService: authService,
	}
}

func (s *service) GenerateUserPinCode(
	ctx context.Context,
) (*string, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	pincode := string(pulid.MustNew("PIN"))

	repository := s.entClient

	err = repository.User.
		UpdateOneID(currentUserID).
		SetContactPin(pincode).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &pincode, nil
}

func (s *service) AddUserContact(
	ctx context.Context,
	pincode string,
) (*ent.UserContact, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	repository := s.entClient

	pincodeUser, err := repository.User.
		Query().
		Where(
			user.ContactPin(pincode),
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	userContact, err := repository.UserContact.
		Create().
		SetUserID(currentUserID).
		SetContactID(pincodeUser.ID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return userContact, nil
}

func (s *service) DeleteUserContact(
	ctx context.Context,
	ID pulid.ID,
) (*ent.UserContact, error) {
	_, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	repository := s.entClient

	userContact, err := repository.UserContact.Get(ctx, ID)
	if err != nil {
		return nil, err
	}

	err = repository.UserContact.
		DeleteOneID(ID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return userContact, nil
}
