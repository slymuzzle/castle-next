package chat

import (
	"context"
	"fmt"
	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/nats"
)

type Service interface {
	SendMessage(
		ctx context.Context,
		currentUserID pulid.ID,
		targetUserID pulid.ID,
		replyTo pulid.ID,
		content string,
	) (*ent.Message, error)
	SubscribeToMessageAddedEvent(roomID pulid.ID) (<-chan *ent.Message, error)
	UpdateMessage(
		ctx context.Context,
		messageID pulid.ID,
		content string,
	) (*ent.Message, error)
	SubscribeToMessageUpdatedEvent(roomID pulid.ID) (<-chan *ent.Message, error)
	DeleteMessage(
		ctx context.Context,
		messageID pulid.ID,
	) (*ent.Message, error)
	SubscribeToMessageDeletedEvent(roomID pulid.ID) (<-chan *ent.Message, error)
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
	replyTo pulid.ID,
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
		replyTo,
		content,
	)
	if err != nil {
		return nil, err
	}

	natsClient := s.natsService.Client()
	subject := fmt.Sprintf("room.%s.message.add", rm.ID)

	if err := natsClient.Publish(subject, msg); err != nil {
		return msg, err
	}

	return msg, nil
}

func (s *service) SubscribeToMessageAddedEvent(roomID pulid.ID) (<-chan *ent.Message, error) {
	subject := fmt.Sprintf("room.%s.message.add", roomID)

	return s.subscribe(subject)
}

func (s *service) UpdateMessage(
	ctx context.Context,
	messageID pulid.ID,
	content string,
) (*ent.Message, error) {
	rm, err := s.chatRepository.FindRoomByMessage(
		ctx,
		messageID,
	)
	if err != nil {
		return nil, err
	}

	msg, err := s.chatRepository.UpdateMessage(
		ctx,
		messageID,
		content,
	)
	if err != nil {
		return nil, err
	}

	natsClient := s.natsService.Client()
	subject := fmt.Sprintf("room.%s.message.update", rm.ID)

	if err := natsClient.Publish(subject, msg); err != nil {
		return msg, err
	}

	return msg, nil
}

func (s *service) SubscribeToMessageUpdatedEvent(roomID pulid.ID) (<-chan *ent.Message, error) {
	subject := fmt.Sprintf("room.%s.message.update", roomID)

	return s.subscribe(subject)
}

func (s *service) DeleteMessage(
	ctx context.Context,
	messageID pulid.ID,
) (*ent.Message, error) {
	rm, err := s.chatRepository.FindRoomByMessage(
		ctx,
		messageID,
	)
	if err != nil {
		return nil, err
	}

	msg, err := s.chatRepository.DeleteMessage(
		ctx,
		messageID,
	)
	if err != nil {
		return nil, err
	}

	natsClient := s.natsService.Client()
	subject := fmt.Sprintf("room.%s.message.delete", rm.ID)

	if err := natsClient.Publish(subject, msg); err != nil {
		return msg, err
	}

	return msg, nil
}

func (s *service) SubscribeToMessageDeletedEvent(roomID pulid.ID) (<-chan *ent.Message, error) {
	subject := fmt.Sprintf("room.%s.message.delete", roomID)

	return s.subscribe(subject)
}

func (s *service) subscribe(subject string) (<-chan *ent.Message, error) {
	natsClient := s.natsService.Client()
	ch := make(chan *ent.Message)

	_, err := natsClient.Subscribe(subject, func(msg *ent.Message) {
		ch <- msg
	})
	if err != nil {
		return ch, err
	}

	return ch, nil
}
