package room

import (
	"context"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/platform/nats"
)

type Service interface {
	FindOrCreatePersonalRoom(
		ctx context.Context,
		currentUserID pulid.ID,
		targetUserID pulid.ID,
	) (*ent.Room, error)

	CreateRoom(
		ctx context.Context,
		input CreateRoomInput,
	) (*ent.Room, error)

	UpdateRoom(
		ctx context.Context,
		ID pulid.ID,
		input UpdateRoomInput,
	) (*ent.Room, error)

	DeleteRoom(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.Room, error)
}

type service struct {
	roomsRepository Repository
	natsService     nats.Service
}

func NewService(
	roomsRepository Repository,
	natsService nats.Service,
) Service {
	return &service{
		roomsRepository: roomsRepository,
		natsService:     natsService,
	}
}

func (s *service) FindOrCreatePersonalRoom(
	ctx context.Context,
	currentUserID pulid.ID,
	targetUserID pulid.ID,
) (*ent.Room, error) {
	return nil, nil
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
) (*ent.Room, error) {
	return s.roomsRepository.Update(ctx, ID, input)
}

func (s *service) DeleteRoom(
	ctx context.Context,
	ID pulid.ID,
) (*ent.Room, error) {
	return s.roomsRepository.Delete(ctx, ID)
}
