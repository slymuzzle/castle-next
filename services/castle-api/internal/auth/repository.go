package auth

import (
	"context"
	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/ent/user"
)

type Repository interface {
	FindUserByID(
		ctx context.Context,
		userID pulid.ID,
	) (*ent.User, error)
	FindUserByNickname(
		ctx context.Context,
		nickname string,
	) (*ent.User, error)
	CreateUser(
		ctx context.Context,
		firstName string,
		lastName string,
		nickname string,
		password string,
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

func (r *repository) FindUserByID(ctx context.Context, userID pulid.ID) (*ent.User, error) {
	return r.entClient.User.Get(ctx, userID)
}

func (r *repository) FindUserByNickname(
	ctx context.Context,
	nickname string,
) (*ent.User, error) {
	return r.entClient.User.
		Query().
		Where(
			user.Nickname(nickname),
		).Only(ctx)
}

func (r *repository) CreateUser(
	ctx context.Context,
	firstName string,
	lastName string,
	nickname string,
	password string,
) (*ent.User, error) {
	return r.entClient.User.
		Create().
		SetFirstName(firstName).
		SetLastName(lastName).
		SetNickname(nickname).
		SetPassword(password).
		Save(ctx)
}
