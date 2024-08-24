package rooms

import (
	"context"
	"fmt"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/graph/model"
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

	IncrementRoomVersion(
		ctx context.Context,
		ID pulid.ID,
		lastMessage *ent.Message,
	) (*ent.Room, error)

	IncrementUnreadMessagesCount(
		ctx context.Context,
		ID pulid.ID,
	) error

	SubscribeToRoomsUpdatedEvent(
		ctx context.Context,
	) (<-chan *model.RoomUpdatedEvent, error)

	DeleteRoomMember(
		ctx context.Context,
		roomMemberID pulid.ID,
	) (*ent.RoomMember, error)

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

func (s *service) IncrementRoomVersion(
	ctx context.Context,
	ID pulid.ID,
	lastMessage *ent.Message,
) (*ent.Room, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	room, err := s.roomsRepository.IncrementVersion(ctx, ID, &lastMessage.ID)
	if err != nil {
		return nil, err
	}

	event := model.RoomUpdatedEvent{
		ID:      room.ID,
		Name:    room.Name,
		Version: room.Version,
		Type:    room.Type,
		LastMessage: &model.LastMessageUpdatedEvent{
			ID:        lastMessage.ID,
			Content:   lastMessage.Content,
			CreatedAt: lastMessage.CreatedAt,
			UpdatedAt: lastMessage.UpdatedAt,
		},
		CreatedAt: room.CreatedAt,
		UpdatedAt: room.UpdatedAt,
	}

	natsClient := s.natsService.Client()

	subject := fmt.Sprintf("users.%s.rooms.updated", currentUserID)
	if err := natsClient.Publish(subject, event); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *service) IncrementUnreadMessagesCount(
	ctx context.Context,
	ID pulid.ID,
) error {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return err
	}

	return s.roomsRepository.IncrementUnreadMessagesCount(ctx, ID, currentUserID)
}

func (s *service) SubscribeToRoomsUpdatedEvent(
	ctx context.Context,
) (<-chan *model.RoomUpdatedEvent, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	subject := fmt.Sprintf("users.%s.rooms.updated", currentUserID)

	return s.subscribe(ctx, subject)
}

func (s *service) DeleteRoomMember(
	ctx context.Context,
	roomMemberID pulid.ID,
) (*ent.RoomMember, error) {
	_, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return s.roomsRepository.DeleteRoomMember(ctx, roomMemberID)
}

func (s *service) DeleteRoom(
	ctx context.Context,
	ID pulid.ID,
) (*ent.Room, error) {
	_, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return s.roomsRepository.Delete(ctx, ID)
}

func (s *service) subscribe(
	ctx context.Context,
	subject string,
) (<-chan *model.RoomUpdatedEvent, error) {
	natsClient := s.natsService.Client()

	ch := make(chan *model.RoomUpdatedEvent, 1)

	sub, err := natsClient.Subscribe(subject, func(room *model.RoomUpdatedEvent) {
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
