package roommembers

import (
	"context"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/modules/auth"
)

type Service interface {
	CreateRoomMembers(
		ctx context.Context,
		inputs []CreateRoomMemberInput,
	) ([]*ent.RoomMember, error)

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
	subscriptions         Subscriptions
	roomMembersRepository Repository
	authService           auth.Service
}

func NewService(
	subscriptions Subscriptions,
	roomMembersRepository Repository,
	authService auth.Service,
) Service {
	return &service{
		subscriptions:         subscriptions,
		roomMembersRepository: roomMembersRepository,
		authService:           authService,
	}
}

func (s *service) CreateRoomMembers(
	ctx context.Context,
	inputs []CreateRoomMemberInput,
) ([]*ent.RoomMember, error) {
	roomMembersToNotify, err := s.roomMembersRepository.CreateBulk(ctx, inputs)
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

func (s *service) IncrementUnreadMessagesCount(
	ctx context.Context,
	roomID pulid.ID,
) ([]*ent.RoomMember, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	roomMembersToNotify, err := s.roomMembersRepository.
		IncrementUnreadMessagesCount(ctx, currentUserID, roomID)
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
	roomMember, err := s.roomMembersRepository.Delete(ctx, ID)
	if err != nil {
		return nil, err
	}

	_, err = s.subscriptions.PublishRoomMemberDeletedEvent(ctx, roomMember.ID)
	if err != nil {
		return nil, err
	}

	return s.roomMembersRepository.Delete(ctx, ID)
}

func (s *service) MarkRoomMemberAsSeen(
	ctx context.Context,
	ID pulid.ID,
) (*ent.RoomMember, error) {
	roomMember, err := s.roomMembersRepository.MarkAsSeen(ctx, ID)
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
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	roomMembersToNotify, err := s.roomMembersRepository.RestoreByRoom(ctx, currentUserID, roomID)
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
