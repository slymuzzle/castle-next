package rooms

import (
	"context"
	"fmt"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/modules/auth"
	"journeyhub/internal/platform/nats"
)

type Service interface {
	FindOrCreatePersonalRoom(
		ctx context.Context,
		targetUserID pulid.ID,
	) (*ent.Room, error)

	FindRoomByMessage(
		ctx context.Context,
		messageID pulid.ID,
	) (*ent.Room, error)

	CreateRoom(
		ctx context.Context,
		input CreateRoomInput,
	) (*ent.Room, error)

	UpdateRoom(
		ctx context.Context,
		ID pulid.ID,
		input UpdateRoomInput,
	) (*ent.RoomEdge, error)

	SubscribeToRoomsUpdatedEvent(
		ctx context.Context,
	) (<-chan *ent.RoomEdge, error)

	DeleteRoom(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.Room, error)
}

type service struct {
	roomsRepository Repository
	authService     auth.Service
	natsService     nats.Service
}

func NewService(
	roomsRepository Repository,
	authService auth.Service,
	natsService nats.Service,
) Service {
	return &service{
		roomsRepository: roomsRepository,
		authService:     authService,
		natsService:     natsService,
	}
}

func (s *service) FindOrCreatePersonalRoom(
	ctx context.Context,
	targetUserID pulid.ID,
) (*ent.Room, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return s.roomsRepository.FindOrCreatePersonal(ctx, currentUserID, targetUserID)
}

func (s *service) FindRoomByMessage(
	ctx context.Context,
	messageID pulid.ID,
) (*ent.Room, error) {
	return s.roomsRepository.FindByMessage(ctx, messageID)
}

func (s *service) CreateRoom(
	ctx context.Context,
	input CreateRoomInput,
) (*ent.Room, error) {
	return s.roomsRepository.Create(ctx, input)
}

func (s *service) UpdateRoom(
	ctx context.Context,
	ID pulid.ID,
	input UpdateRoomInput,
) (*ent.RoomEdge, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	room, err := s.roomsRepository.Update(ctx, ID, input)
	if err != nil {
		return nil, err
	}

	roomEdge := room.ToEdge(ent.DefaultRoomOrder)

	natsClient := s.natsService.Client()

	subject := fmt.Sprintf("users.%s.rooms.updated", currentUserID)
	if err := natsClient.Publish(subject, roomEdge); err != nil {
		return roomEdge, err
	}

	return roomEdge, nil
}

func (s *service) SubscribeToRoomsUpdatedEvent(
	ctx context.Context,
) (<-chan *ent.RoomEdge, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	subject := fmt.Sprintf("users.%s.rooms.updated", currentUserID)

	return s.subscribe(ctx, subject)
}

func (s *service) DeleteRoom(
	ctx context.Context,
	ID pulid.ID,
) (*ent.Room, error) {
	return s.roomsRepository.Delete(ctx, ID)
}

func (s *service) subscribe(
	ctx context.Context,
	subject string,
) (<-chan *ent.RoomEdge, error) {
	natsClient := s.natsService.Client()

	ch := make(chan *ent.RoomEdge, 1)

	sub, err := natsClient.Subscribe(subject, func(room *ent.RoomEdge) {
		ch <- room
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
