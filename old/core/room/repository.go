package room

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
		filter RoomWhereInput,
	) ([]*ent.Room, error)

	FindOne(
		ctx context.Context,
		filter RoomWhereInput,
	) (*ent.Room, error)

	Create(
		ctx context.Context,
		input CreateRoomInput,
	) (*ent.Room, error)

	Update(
		ctx context.Context,
		ID pulid.ID,
		input UpdateRoomInput,
	) (*ent.Room, error)

	Delete(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.Room, error)
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
	where RoomWhereInput,
) ([]*ent.Room, error) {
	q := r.entClient.Room.
		Query().
		Limit(limit).
		Offset(offset)
	where.Filter(q)
	return q.All(ctx)
}

func (r *repository) FindOne(
	ctx context.Context,
	where RoomWhereInput,
) (*ent.Room, error) {
	q := r.entClient.Room.Query()
	where.Filter(q)
	return q.Only(ctx)
}

func (r *repository) Create(
	ctx context.Context,
	input CreateRoomInput,
) (*ent.Room, error) {
	u := r.entClient.Room.Create()
	input.Mutate(u.Mutation())
	return u.Save(ctx)
}

func (r *repository) Update(
	ctx context.Context,
	ID pulid.ID,
	input UpdateRoomInput,
) (*ent.Room, error) {
	u := r.entClient.Room.UpdateOneID(ID)
	input.Mutate(u.Mutation())
	return u.Save(ctx)
}

func (r *repository) Delete(
	ctx context.Context,
	ID pulid.ID,
) (*ent.Room, error) {
	rm, err := r.entClient.Room.Get(ctx, ID)
	if err != nil {
		return nil, err
	}

	err = r.entClient.Room.
		DeleteOneID(ID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return rm, nil
}
