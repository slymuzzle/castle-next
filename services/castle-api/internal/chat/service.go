package chat

import (
	"context"
	"fmt"
	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"

	"github.com/nats-io/nats.go"
)

type Service interface {
	SendMessage(
		ctx context.Context,
		roomID pulid.ID,
		userID pulid.ID,
		content string,
	) (*ent.Message, error)
	Subscribe(roomID pulid.ID) (<-chan *ent.Message, error)
}

type service struct {
	entClient *ent.Client
	natsConn  *nats.EncodedConn
}

func NewService(entClient *ent.Client, natsConn *nats.EncodedConn) Service {
	return &service{
		entClient: entClient,
		natsConn:  natsConn,
	}
}

func (s *service) SendMessage(
	ctx context.Context,
	roomID pulid.ID,
	userID pulid.ID,
	content string,
) (*ent.Message, error) {
	msg, err := s.entClient.Message.
		Create().
		SetRoomID(roomID).
		SetUserID(userID).
		SetContent(content).
		Save(ctx)
	if err != nil {
		return msg, err
	}

	subject := fmt.Sprintf("room.%s", roomID)

	if err := s.natsConn.Publish(subject, msg); err != nil {
		return msg, err
	}

	return msg, nil
}

func (s *service) Subscribe(roomID pulid.ID) (<-chan *ent.Message, error) {
	ch := make(chan *ent.Message)

	subject := fmt.Sprintf("room.%s", roomID)

	_, err := s.natsConn.Subscribe(subject, func(msg *ent.Message) {
		ch <- msg
	})
	if err != nil {
		return ch, err
	}

	return ch, nil
}
