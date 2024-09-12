package roommembers

import (
	"context"

	"journeyhub/ent"
	"journeyhub/ent/roommember"
	"journeyhub/ent/schema/mixin"
	"journeyhub/ent/schema/pulid"
)

type CreateRoomMemberInput struct {
	UserID pulid.ID
	RoomID pulid.ID
}

func (c *CreateRoomMemberInput) Mutate(m *ent.RoomMemberMutation) {
	m.SetUserID(c.UserID)
	m.SetRoomID(c.RoomID)
}

type Repository interface {
	CreateBulk(
		ctx context.Context,
		inputs []CreateRoomMemberInput,
	) ([]*ent.RoomMember, error)

	FindByID(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.RoomMember, error)

	IncrementUnreadMessagesCount(
		ctx context.Context,
		currentUserID pulid.ID,
		ID pulid.ID,
	) ([]*ent.RoomMember, error)

	Delete(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.RoomMember, error)

	MarkAsSeen(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.RoomMember, error)

	RestoreByRoom(
		ctx context.Context,
		currentUserID pulid.ID,
		roomID pulid.ID,
	) ([]*ent.RoomMember, error)
}

type repository struct {
	entClient *ent.Client
}

func NewRepository(entClient *ent.Client) Repository {
	return &repository{
		entClient: entClient,
	}
}

func (r *repository) CreateBulk(
	ctx context.Context,
	inputs []CreateRoomMemberInput,
) ([]*ent.RoomMember, error) {
	return r.getClient(ctx).RoomMember.MapCreateBulk(inputs, func(c *ent.RoomMemberCreate, i int) {
		c.SetUserID(inputs[i].UserID).SetRoomID(inputs[i].RoomID)
	}).Save(ctx)
}

func (r *repository) FindByID(
	ctx context.Context,
	ID pulid.ID,
) (*ent.RoomMember, error) {
	return r.getClient(ctx).RoomMember.Get(ctx, ID)
}

func (r *repository) IncrementUnreadMessagesCount(
	ctx context.Context,
	currentUserID pulid.ID,
	roomID pulid.ID,
) ([]*ent.RoomMember, error) {
	client := r.getClient(ctx)

	err := client.RoomMember.
		Update().
		Where(
			roommember.And(
				roommember.UserIDNEQ(currentUserID),
				roommember.RoomID(roomID),
			),
		).
		AddUnreadMessagesCount(1).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return client.RoomMember.
		Query().
		Where(
			roommember.And(
				roommember.UserIDNEQ(currentUserID),
				roommember.RoomID(roomID),
			),
		).
		WithRoom().
		WithUser().
		All(ctx)
}

func (r *repository) Delete(
	ctx context.Context,
	ID pulid.ID,
) (*ent.RoomMember, error) {
	client := r.getClient(ctx)

	roomMember, err := client.RoomMember.Get(ctx, ID)
	if err != nil {
		return nil, err
	}

	err = client.RoomMember.
		DeleteOneID(ID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return roomMember, nil
}

func (r *repository) MarkAsSeen(
	ctx context.Context,
	ID pulid.ID,
) (*ent.RoomMember, error) {
	return r.getClient(ctx).RoomMember.
		UpdateOneID(ID).
		SetUnreadMessagesCount(0).
		Save(ctx)
}

func (r *repository) RestoreByRoom(
	ctx context.Context,
	currentUserID pulid.ID,
	roomID pulid.ID,
) ([]*ent.RoomMember, error) {
	client := r.getClient(ctx)

	roomMembers, err := client.RoomMember.
		Query().
		Where(
			roommember.RoomID(roomID),
			roommember.DeletedAtNotNil(),
		).All(mixin.SkipSoftDelete(ctx))
	if err != nil {
		return nil, err
	}

	err = client.RoomMember.
		Update().
		Where(
			roommember.RoomID(roomID),
			roommember.DeletedAtNotNil(),
		).
		ClearDeletedAt().
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return roomMembers, nil
}

func (r *repository) getClient(ctx context.Context) *ent.Client {
	var client *ent.Client
	if clientFromCtx := ent.FromContext(ctx); clientFromCtx != nil {
		client = clientFromCtx
	} else {
		client = r.entClient
	}
	return client
}
