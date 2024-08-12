package chat

import (
	"context"
	"errors"
	"fmt"
	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/nats"
)

var (
	ErrCreateRoom         = errors.New("failed to create room")
	ErrTargetUserNotExist = errors.New("target user not exist")
)

type Service interface {
	SendMessage(
		ctx context.Context,
		currentUserID pulid.ID,
		targetUserID pulid.ID,
		content string,
	) (*ent.Message, error)
	Subscribe(roomID pulid.ID) (<-chan *ent.Message, error)
}

type service struct {
	chatRepository Repository
	natsService    nats.Service
}

func NewService(chatRepository Repository, natsConn nats.Service) Service {
	return &service{
		chatRepository: chatRepository,
		natsService:    natsConn,
	}
}

func (s *service) SendMessage(
	ctx context.Context,
	currentUserID pulid.ID,
	targetUserID pulid.ID,
	content string,
) (*ent.Message, error) {
	rm, err := s.chatRepository.FindOrCreatePersonalRoom(
		ctx,
		currentUserID,
		targetUserID,
	)
	if err != nil {
		return nil, err
	}

	msg, err := s.chatRepository.CreateMessage(
		ctx,
		rm.ID,
		currentUserID,
		content,
	)
	if err != nil {
		return nil, err
	}

	natsClient := s.natsService.Client()
	subject := fmt.Sprintf("room.%s", rm.ID)

	if err := natsClient.Publish(subject, msg); err != nil {
		return msg, err
	}

	return msg, nil
}

func (s *service) Subscribe(roomID pulid.ID) (<-chan *ent.Message, error) {
	natsClient := s.natsService.Client()
	ch := make(chan *ent.Message)

	subject := fmt.Sprintf("room.%s", roomID)

	_, err := natsClient.Subscribe(subject, func(msg *ent.Message) {
		ch <- msg
	})
	if err != nil {
		return ch, err
	}

	return ch, nil
}
