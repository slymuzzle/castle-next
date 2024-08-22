package contacts

import (
	"context"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
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
}

type service struct {
	contactsRepository Repository
	authService        auth.Service
}

func NewService(userContactsRepository Repository, authService auth.Service) Service {
	return &service{
		contactsRepository: userContactsRepository,
		authService:        authService,
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

	_, err = s.contactsRepository.UpdateUserPinCode(
		ctx,
		currentUserID,
		pincode,
	)
	if err != nil {
		return nil, err
	}

	return &pincode, nil
}

func (s *service) AddUserContact(ctx context.Context,
	pincode string,
) (*ent.UserContact, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	pincodeUser, err := s.contactsRepository.FindUserByPinCode(ctx, pincode)
	if err != nil {
		return nil, err
	}

	userContact, err := s.contactsRepository.Create(ctx, CreateUserContactInput{
		UserID:    currentUserID,
		ContactID: pincodeUser.ID,
	})
	if err != nil {
		return nil, err
	}

	return userContact, nil
}
