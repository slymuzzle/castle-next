package message

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
		filter MessageWhereInput,
	) ([]*ent.Message, error)

	FindOne(
		ctx context.Context,
		filter MessageWhereInput,
	) (*ent.Message, error)

	Create(
		ctx context.Context,
		input CreateMessageInput,
	) (*ent.Message, error)

	Update(
		ctx context.Context,
		ID pulid.ID,
		input UpdateMessageInput,
	) (*ent.Message, error)

	Delete(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.Message, error)
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
	where MessageWhereInput,
) ([]*ent.Message, error) {
	q := r.entClient.Message.
		Query().
		Limit(limit).
		Offset(offset)
	where.Filter(q)
	return q.All(ctx)
}

func (r *repository) FindOne(
	ctx context.Context,
	where MessageWhereInput,
) (*ent.Message, error) {
	q := r.entClient.Message.Query()
	where.Filter(q)
	return q.Only(ctx)
}

func (r *repository) Create(
	ctx context.Context,
	input CreateMessageInput,
) (*ent.Message, error) {
	u := r.entClient.Message.Create()
	input.Mutate(u.Mutation())
	return u.Save(ctx)
}

func (r *repository) Update(
	ctx context.Context,
	ID pulid.ID,
	input UpdateMessageInput,
) (*ent.Message, error) {
	u := r.entClient.Message.UpdateOneID(ID)
	input.Mutate(u.Mutation())
	return u.Save(ctx)
}

func (r *repository) Delete(
	ctx context.Context,
	ID pulid.ID,
) (*ent.Message, error) {
	rm, err := r.entClient.Message.Get(ctx, ID)
	if err != nil {
		return nil, err
	}

	err = r.entClient.Message.
		DeleteOneID(ID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return rm, nil
}
