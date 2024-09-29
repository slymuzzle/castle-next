package chat

import (
	"context"
	"fmt"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/platform/nats"
)

type Subscriptions interface {
	PublishMessageCreatedEvent(
		ctx context.Context,
		roomID pulid.ID,
		messageID pulid.ID,
	) (string, error)

	SubscribeToMessageCreatedEvent(
		ctx context.Context,
		roomID pulid.ID,
	) (<-chan *ent.MessageEdge, error)

	PublishMessageUpdatedEvent(
		ctx context.Context,
		roomID pulid.ID,
		messageID pulid.ID,
	) (string, error)

	SubscribeToMessageUpdatedEvent(
		ctx context.Context,
		roomID pulid.ID,
	) (<-chan *ent.MessageEdge, error)

	PublishMessageDeletedEvent(
		ctx context.Context,
		roomID pulid.ID,
		messageID pulid.ID,
	) (string, error)

	SubscribeToMessageDeletedEvent(
		ctx context.Context,
		roomID pulid.ID,
	) (<-chan pulid.ID, error)
}

type subscriptions struct {
	entClient   *ent.Client
	natsService nats.Service
}

func NewSubscriptions(
	entClient *ent.Client,
	natsService nats.Service,
) Subscriptions {
	return &subscriptions{
		entClient:   entClient,
		natsService: natsService,
	}
}

func (s *subscriptions) PublishMessageCreatedEvent(
	ctx context.Context,
	roomID pulid.ID,
	messageID pulid.ID,
) (string, error) {
	natsClient := s.natsService.Client()

	subject := fmt.Sprintf("room.%s.message.created", roomID)
	if err := natsClient.Publish(subject, messageID); err != nil {
		return "", err
	}

	return subject, nil
}

func (s *subscriptions) SubscribeToMessageCreatedEvent(
	ctx context.Context,
	roomID pulid.ID,
) (<-chan *ent.MessageEdge, error) {
	subject := fmt.Sprintf("room.%s.message.created", roomID)

	return s.subscribe(ctx, subject)
}

func (s *subscriptions) PublishMessageUpdatedEvent(
	ctx context.Context,
	roomID pulid.ID,
	messageID pulid.ID,
) (string, error) {
	natsClient := s.natsService.Client()

	subject := fmt.Sprintf("room.%s.message.updated", roomID)
	if err := natsClient.Publish(subject, messageID); err != nil {
		return "", err
	}

	return subject, nil
}

func (s *subscriptions) SubscribeToMessageUpdatedEvent(
	ctx context.Context,
	roomID pulid.ID,
) (<-chan *ent.MessageEdge, error) {
	subject := fmt.Sprintf("room.%s.message.updated", roomID)

	return s.subscribe(ctx, subject)
}

func (s *subscriptions) PublishMessageDeletedEvent(
	ctx context.Context,
	roomID pulid.ID,
	messageID pulid.ID,
) (string, error) {
	natsClient := s.natsService.Client()

	subject := fmt.Sprintf("room.%s.message.deleted", roomID)
	if err := natsClient.Publish(subject, messageID); err != nil {
		return "", err
	}

	return subject, nil
}

func (s *subscriptions) SubscribeToMessageDeletedEvent(
	ctx context.Context,
	roomID pulid.ID,
) (<-chan pulid.ID, error) {
	subject := fmt.Sprintf("room.%s.message.deleted", roomID)

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
) (<-chan *ent.MessageEdge, error) {
	natsClient := s.natsService.Client()

	ch := make(chan *ent.MessageEdge, 1)

	sub, err := natsClient.Subscribe(subject, func(messageID pulid.ID) {
		message, err := s.entClient.Message.Get(ctx, messageID)
		if err != nil {
			return
		}
		ch <- message.ToEdge(ent.DefaultMessageOrder)
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
