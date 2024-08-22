package message

import (
	"context"

	"journeyhub/ent"
	"journeyhub/internal/platform/media"
	"journeyhub/internal/platform/nats"
)

var (
	MessageCreatedSubject = "room.%s.message.created"
	MessageUpdatedSubject = "room.%s.message.updated"
	MessageDeletedSubject = "room.%s.message.deleted"
)

type Service interface {
	CreateMessage(
		ctx context.Context,
		input CreateMessageInput,
	) (*ent.Message, error)
}

type service struct {
	messageRepository Repository
	mediaService      media.Service
}

func NewService(
	messageRepository Repository,
	natsService nats.Service,
	mediaService media.Service,
) Service {
	return &service{
		messageRepository: messageRepository,
		mediaService:      mediaService,
	}
}

func (s *service) CreateMessage(
	ctx context.Context,
	input CreateMessageInput,
) (*ent.Message, error) {
	return s.messageRepository.Create(ctx, input)
}
