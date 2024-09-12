package roommembers

import (
	"context"
	"fmt"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/modules/auth"
	"journeyhub/internal/platform/nats"
)

type Subscriptions interface {
	PublishRoomMemberCreatedEvent(
		ctx context.Context,
		userID pulid.ID,
		roomMemberID pulid.ID,
	) (string, error)

	SubscribeToRoomMemberCreatedEvent(
		ctx context.Context,
	) (<-chan *ent.RoomMemberEdge, error)

	PublishRoomMemberUpdatedEvent(
		ctx context.Context,
		userID pulid.ID,
		roomMemberID pulid.ID,
	) (string, error)

	SubscribeToRoomMemberUpdatedEvent(
		ctx context.Context,
	) (<-chan *ent.RoomMemberEdge, error)

	PublishRoomMemberDeletedEvent(
		ctx context.Context,
		roomMemberID pulid.ID,
	) (string, error)

	SubscribeToRoomMemberDeletedEvent(
		ctx context.Context,
	) (<-chan pulid.ID, error)
}

type subscriptions struct {
	roomMembersRepository Repository
	authService           auth.Service
	natsService           nats.Service
}

func NewSubscriptions(
	roomMembersRepository Repository,
	authService auth.Service,
	natsService nats.Service,
) Subscriptions {
	return &subscriptions{
		roomMembersRepository: roomMembersRepository,
		authService:           authService,
		natsService:           natsService,
	}
}

func (s *subscriptions) PublishRoomMemberCreatedEvent(
	ctx context.Context,
	userID pulid.ID,
	roomMemberID pulid.ID,
) (string, error) {
	natsClient := s.natsService.Client()

	subject := fmt.Sprintf("users.%s.roommembers.created", userID)
	if err := natsClient.Publish(subject, roomMemberID); err != nil {
		return "", err
	}

	return subject, nil
}

func (s *subscriptions) SubscribeToRoomMemberCreatedEvent(
	ctx context.Context,
) (<-chan *ent.RoomMemberEdge, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	subject := fmt.Sprintf("users.%s.roommembers.created", currentUserID)

	return s.subscribe(ctx, subject)
}

func (s *subscriptions) PublishRoomMemberUpdatedEvent(
	ctx context.Context,
	userID pulid.ID,
	roomMemberID pulid.ID,
) (string, error) {
	natsClient := s.natsService.Client()

	subject := fmt.Sprintf("users.%s.roommembers.updated", userID)
	if err := natsClient.Publish(subject, roomMemberID); err != nil {
		return "", err
	}

	return subject, nil
}

func (s *subscriptions) SubscribeToRoomMemberUpdatedEvent(
	ctx context.Context,
) (<-chan *ent.RoomMemberEdge, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	subject := fmt.Sprintf("users.%s.roommembers.updated", currentUserID)

	return s.subscribe(ctx, subject)
}

func (s *subscriptions) PublishRoomMemberDeletedEvent(
	ctx context.Context,
	roomMemberID pulid.ID,
) (string, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return "", err
	}

	natsClient := s.natsService.Client()

	subject := fmt.Sprintf("users.%s.roommembers.deleted", currentUserID)
	if err := natsClient.Publish(subject, roomMemberID); err != nil {
		return "", err
	}

	return subject, nil
}

func (s *subscriptions) SubscribeToRoomMemberDeletedEvent(
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

func (s *subscriptions) subscribe(
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
