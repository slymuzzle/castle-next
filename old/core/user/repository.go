package user

import (
	"context"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
)

type Repository interface {
	FindAll(
		ctx context.Context,
		limit int,
		offset int,
		filter UserWhereInput,
	) ([]*ent.User, error)

	FindOne(
		ctx context.Context,
		filter UserWhereInput,
	) (*ent.User, error)

	Create(
		ctx context.Context,
		input CreateUserInput,
	) (*ent.User, error)

	Update(
		ctx context.Context,
		ID pulid.ID,
		input UpdateUserInput,
	) (*ent.User, error)

	Delete(
		ctx context.Context,
		ID pulid.ID,
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

func (r *repository) FindAll(
	ctx context.Context,
	limit int,
	offset int,
	where UserWhereInput,
) ([]*ent.User, error) {
	q := r.entClient.User.
		Query().
		Limit(limit).
		Offset(offset)
	where.Filter(q)
	return q.All(ctx)
}

func (r *repository) FindOne(
	ctx context.Context,
	where UserWhereInput,
) (*ent.User, error) {
	q := r.entClient.User.Query()
	where.Filter(q)
	return q.Only(ctx)
}

func (r *repository) Create(
	ctx context.Context,
	input CreateUserInput,
) (*ent.User, error) {
	u := r.entClient.User.Create()
	input.Mutate(u.Mutation())
	return u.Save(ctx)
}

func (r *repository) Update(
	ctx context.Context,
	ID pulid.ID,
	input UpdateUserInput,
) (*ent.User, error) {
	u := r.entClient.User.UpdateOneID(ID)
	input.Mutate(u.Mutation())
	return u.Save(ctx)
}

func (r *repository) Delete(
	ctx context.Context,
	ID pulid.ID,
) (*ent.User, error) {
	uc, err := r.entClient.User.Get(ctx, ID)
	if err != nil {
		return nil, err
	}

	err = r.entClient.User.
		DeleteOneID(ID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return uc, nil
}
