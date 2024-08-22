package user

import (
	"context"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
)

type Service interface {
	FindOneUser(
		ctx context.Context,
		where UserWhereInput,
	) (*ent.User, error)

	CreateUser(
		ctx context.Context,
		input CreateUserInput,
	) (*ent.User, error)

	UpdateUser(
		ctx context.Context,
		ID pulid.ID,
		input UpdateUserInput,
	) (*ent.User, error)

	DeleteUser(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.User, error)
}

type service struct {
	usersRepository Repository
}

func NewService(usersRepository Repository) Service {
	return &service{
		usersRepository: usersRepository,
	}
}

func (s *service) FindOneUser(
	ctx context.Context,
	where UserWhereInput,
) (*ent.User, error) {
	return s.usersRepository.FindOne(ctx, where)
}

func (s *service) CreateUser(
	ctx context.Context,
	input CreateUserInput,
) (*ent.User, error) {
	return s.usersRepository.Create(ctx, input)
}

func (s *service) UpdateUser(
	ctx context.Context,
	ID pulid.ID,
	input UpdateUserInput,
) (*ent.User, error) {
	return s.usersRepository.Update(ctx, ID, input)
}

func (s *service) DeleteUser(
	ctx context.Context,
	ID pulid.ID,
) (*ent.User, error) {
	return s.usersRepository.Delete(ctx, ID)
}
