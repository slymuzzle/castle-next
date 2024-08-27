package roommembers

import (
	"context"
	"fmt"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/modules/auth"
	"journeyhub/internal/platform/nats"
)

type Service interface {
	CreateRoomMembers(
		ctx context.Context,
		inputs []CreateRoomMemberInput,
	) ([]*ent.RoomMember, error)

	SubscribeToRoomMemberCreatedEvent(
		ctx context.Context,
	) (<-chan *ent.RoomMemberEdge, error)

	IncrementUnreadMessagesCount(
		ctx context.Context,
		roomID pulid.ID,
	) ([]*ent.RoomMember, error)

	SubscribeToRoomMemberUpdatedEvent(
		ctx context.Context,
	) (<-chan *ent.RoomMemberEdge, error)

	DeleteRoomMember(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.RoomMember, error)

	SubscribeToRoomMemberDeletedEvent(
		ctx context.Context,
	) (<-chan pulid.ID, error)

	MarkRoomMemberAsSeen(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.RoomMember, error)

	RestoreRoomMembersByRoom(
		ctx context.Context,
		roomID pulid.ID,
	) ([]*ent.RoomMember, error)
}

type service struct {
	roomMembersRepository Repository
	authService           auth.Service
	natsService           nats.Service
}

func NewService(
	roomMembersRepository Repository,
	authService auth.Service,
	natsService nats.Service,
) Service {
	return &service{
		roomMembersRepository: roomMembersRepository,
		authService:           authService,
		natsService:           natsService,
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

	natsClient := s.natsService.Client()
	for _, roomMember := range roomMembersToNotify {
		subject := fmt.Sprintf("users.%s.roommembers.created", roomMember.UserID)
		if err := natsClient.Publish(subject, roomMember.ID); err != nil {
			return nil, err
		}
	}

	return roomMembersToNotify, nil
}

func (s *service) SubscribeToRoomMemberCreatedEvent(
	ctx context.Context,
) (<-chan *ent.RoomMemberEdge, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	subject := fmt.Sprintf("users.%s.roommembers.created", currentUserID)

	return s.subscribe(ctx, subject)
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

	natsClient := s.natsService.Client()
	for _, roomMember := range roomMembersToNotify {
		subject := fmt.Sprintf("users.%s.roommembers.updated", roomMember.UserID)
		if err := natsClient.Publish(subject, roomMember.ID); err != nil {
			return nil, err
		}
	}

	return roomMembersToNotify, nil
}

func (s *service) SubscribeToRoomMemberUpdatedEvent(
	ctx context.Context,
) (<-chan *ent.RoomMemberEdge, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	subject := fmt.Sprintf("users.%s.roommembers.updated", currentUserID)

	return s.subscribe(ctx, subject)
}

func (s *service) DeleteRoomMember(
	ctx context.Context,
	ID pulid.ID,
) (*ent.RoomMember, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	natsClient := s.natsService.Client()
	subject := fmt.Sprintf("users.%s.roommembers.deleted", currentUserID)
	if err := natsClient.Publish(subject, ID); err != nil {
		return nil, err
	}

	return s.roomMembersRepository.Delete(ctx, ID)
}

func (s *service) SubscribeToRoomMemberDeletedEvent(
	ctx context.Context,
) (<-chan pulid.ID, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	subject := fmt.Sprintf("users.%s.roommembers.deleted", currentUserID)

	natsClient := s.natsService.Client()

	ch := make(chan pulid.ID, 1)

	sub, err := natsClient.Subscribe(subject, func(messageID pulid.ID) {
		ch <- messageID
	})
	if err != nil {
		return ch, err
	}

	go func() {
		<-ctx.Done()
		sub.Unsubscribe()
	}()

	return ch, nil
}

func (s *service) MarkRoomMemberAsSeen(
	ctx context.Context,
	ID pulid.ID,
) (*ent.RoomMember, error) {
	roomMember, err := s.roomMembersRepository.MarkAsSeen(ctx, ID)
	if err != nil {
		return nil, err
	}

	natsClient := s.natsService.Client()
	subject := fmt.Sprintf("users.%s.roommembers.updated", roomMember.UserID)
	if err := natsClient.Publish(subject, roomMember.ID); err != nil {
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

	natsClient := s.natsService.Client()
	for _, roomMember := range roomMembersToNotify {
		subject := fmt.Sprintf("users.%s.roommembers.created", roomMember.UserID)
		if err := natsClient.Publish(subject, roomMember.ID); err != nil {
			return nil, err
		}
	}

	return roomMembersToNotify, nil
}

func (s *service) subscribe(
	ctx context.Context,
	subject string,
) (<-chan *ent.RoomMemberEdge, error) {
	natsClient := s.natsService.Client()

	ch := make(chan *ent.RoomMemberEdge, 1)

	sub, err := natsClient.Subscribe(subject, func(roomMemberID pulid.ID) {
		rm, err := s.roomMembersRepository.FindByID(ctx, roomMemberID)
		if err != nil {
			return
		}
		ch <- rm.ToEdge(ent.DefaultRoomMemberOrder)
	})
	if err != nil {
		return ch, err
	}

	go func() {
		<-ctx.Done()
		sub.Unsubscribe()
	}()

	return ch, nil
}
