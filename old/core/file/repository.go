package file

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
		filter FileWhereInput,
	) ([]*ent.File, error)

	FindOne(
		ctx context.Context,
		filter FileWhereInput,
	) (*ent.File, error)

	Create(
		ctx context.Context,
		input CreateFileInput,
	) (*ent.File, error)

	Update(
		ctx context.Context,
		ID pulid.ID,
		input UpdateFileInput,
	) (*ent.File, error)

	Delete(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.File, error)
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
	where FileWhereInput,
) ([]*ent.File, error) {
	q := r.entClient.File.
		Query().
		Limit(limit).
		Offset(offset)
	where.Filter(q)
	return q.All(ctx)
}

func (r *repository) FindOne(
	ctx context.Context,
	where FileWhereInput,
) (*ent.File, error) {
	q := r.entClient.File.Query()
	where.Filter(q)
	return q.Only(ctx)
}

func (r *repository) Create(
	ctx context.Context,
	input CreateFileInput,
) (*ent.File, error) {
	u := r.entClient.File.Create()
	input.Mutate(u.Mutation())
	return u.Save(ctx)
}

func (r *repository) Update(
	ctx context.Context,
	ID pulid.ID,
	input UpdateFileInput,
) (*ent.File, error) {
	u := r.entClient.File.UpdateOneID(ID)
	input.Mutate(u.Mutation())
	return u.Save(ctx)
}

func (r *repository) Delete(
	ctx context.Context,
	ID pulid.ID,
) (*ent.File, error) {
	uc, err := r.entClient.File.Get(ctx, ID)
	if err != nil {
		return nil, err
	}

	err = r.entClient.File.
		DeleteOneID(ID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return uc, nil
}
