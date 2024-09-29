package roommembers

import (
	"context"

	"journeyhub/ent"
	"journeyhub/ent/roommember"
	"journeyhub/ent/schema/mixin"
	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/modules/auth"
)

type Service interface {
	IncrementUnreadMessagesCount(
		ctx context.Context,
		roomID pulid.ID,
	) ([]*ent.RoomMember, error)

	DeleteRoomMember(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.RoomMember, error)

	MarkRoomMemberAsSeen(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.RoomMember, error)

	RestoreRoomMembersByRoom(
		ctx context.Context,
		roomID pulid.ID,
	) ([]*ent.RoomMember, error)

	Subscriptions() Subscriptions
}

type service struct {
	entClient     *ent.Client
	subscriptions Subscriptions
	authService   auth.Service
}

func NewService(
	entClient *ent.Client,
	subscriptions Subscriptions,
	authService auth.Service,
) Service {
	return &service{
		entClient:     entClient,
		subscriptions: subscriptions,
		authService:   authService,
	}
}

func (s *service) IncrementUnreadMessagesCount(
	ctx context.Context,
	roomID pulid.ID,
) ([]*ent.RoomMember, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	repository := s.entClient

	err = repository.RoomMember.
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

	roomMembersToNotify, err := repository.RoomMember.
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
	if err != nil {
		return nil, err
	}

	for _, roomMember := range roomMembersToNotify {
		_, err := s.subscriptions.PublishRoomMemberUpdatedEvent(ctx, roomMember.UserID, roomMember.ID)
		if err != nil {
			return nil, err
		}
	}

	return roomMembersToNotify, nil
}

func (s *service) DeleteRoomMember(
	ctx context.Context,
	ID pulid.ID,
) (*ent.RoomMember, error) {
	repository := s.entClient

	roomMember, err := repository.RoomMember.Get(ctx, ID)
	if err != nil {
		return nil, err
	}

	err = repository.RoomMember.
		DeleteOneID(ID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	_, err = s.subscriptions.PublishRoomMemberDeletedEvent(ctx, roomMember.ID)
	if err != nil {
		return nil, err
	}

	return roomMember, nil
}

func (s *service) MarkRoomMemberAsSeen(
	ctx context.Context,
	ID pulid.ID,
) (*ent.RoomMember, error) {
	repository := s.entClient

	roomMember, err := repository.RoomMember.
		UpdateOneID(ID).
		SetUnreadMessagesCount(0).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	_, err = s.subscriptions.PublishRoomMemberUpdatedEvent(ctx, roomMember.UserID, roomMember.ID)
	if err != nil {
		return nil, err
	}

	return roomMember, nil
}

func (s *service) RestoreRoomMembersByRoom(
	ctx context.Context,
	roomID pulid.ID,
) ([]*ent.RoomMember, error) {
	_, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	repository := s.entClient

	roomMembersToNotify, err := repository.RoomMember.
		Query().
		Where(
			roommember.RoomID(roomID),
			roommember.DeletedAtNotNil(),
		).
		All(mixin.SkipSoftDelete(ctx))
	if err != nil {
		return nil, err
	}

	err = repository.RoomMember.
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

	for _, roomMember := range roomMembersToNotify {
		_, err := s.subscriptions.PublishRoomMemberCreatedEvent(ctx, roomMember.UserID, roomMember.ID)
		if err != nil {
			return nil, err
		}
	}

	return roomMembersToNotify, nil
}

func (s *service) Subscriptions() Subscriptions {
	return s.subscriptions
}
