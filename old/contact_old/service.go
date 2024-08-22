package contacts

import (
	"context"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
)

type Service interface {
	GenerateUserPinCode(
		ctx context.Context,
		userID pulid.ID,
	) (*pulid.ID, error)
	AddUserContact(
		ctx context.Context,
		userID pulid.ID,
		pincode pulid.ID,
	) (*ent.User, error)
}

type service struct {
	contactsRepository Repository
}

func NewService(contactsRepository Repository) Service {
	return &service{
		contactsRepository: contactsRepository,
	}
}

func (s *service) GenerateUserPinCode(
	ctx context.Context,
	userID pulid.ID,
) (*pulid.ID, error) {
	pincode := pulid.MustNew("PIN")

	err := s.contactsRepository.SetUserPinCode(ctx, userID, pincode)
	if err != nil {
		return nil, err
	}

	return &pincode, nil
}

func (s *service) AddUserContact(ctx context.Context,
	userID pulid.ID,
	pincode pulid.ID,
) (*ent.User, error) {
	user, err := s.contactsRepository.
		CreateContactByPinCode(ctx, userID, pincode)
	if err != nil {
		return nil, err
	}

	return user, nil
}
