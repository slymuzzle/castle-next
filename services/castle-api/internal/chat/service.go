package chat

import (
	"context"
	"fmt"
	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/graph/model"
	"journeyhub/internal/media"
	"journeyhub/internal/nats"
)

type (
	ID                 = pulid.ID
	SendMessageInput   = model.SendMessageInput
	UpdateMessageInput = model.UpdateMessageInput
	Message            = ent.Message
	Room               = ent.Room
)

type Service interface {
	SendMessage(
		ctx context.Context,
		currentUserID ID,
		input SendMessageInput,
	) (*Message, error)
	SubscribeToMessageAddedEvent(
		roomID ID,
	) (<-chan *Message, error)
	UpdateMessage(
		ctx context.Context,
		messageID ID,
		input UpdateMessageInput,
	) (*Message, error)
	SubscribeToMessageUpdatedEvent(
		roomID ID,
	) (<-chan *Message, error)
	DeleteMessage(
		ctx context.Context,
		messageID ID,
	) (*Message, error)
	SubscribeToMessageDeletedEvent(
		roomID ID,
	) (<-chan *Message, error)
}

type service struct {
	chatRepository Repository
	natsService    nats.Service
	mediaService   media.Service
}

func NewService(
	chatRepository Repository,
	natsService nats.Service,
	mediaService media.Service,
) Service {
	return &service{
		chatRepository: chatRepository,
		natsService:    natsService,
		mediaService:   mediaService,
	}
}

func (s *service) SendMessage(
	ctx context.Context,
	currentUserID ID,
	input SendMessageInput,
) (*Message, error) {
	rm, err := s.chatRepository.FindOrCreatePersonalRoom(
		ctx,
		currentUserID,
		input.TargetUserID,
	)
	if err != nil {
		return nil, err
	}

	var uploadAttachmentsFn UploadAttachmentsFn
	if input.Files != nil && len(input.Files) > 0 {
		uploadAttachmentsFn = func(m *ent.Message) ([]*media.UploadInfo, error) {
			uploadPrefix := fmt.Sprintf("rooms/%s/%s/attachments", rm.ID, m.ID)
			return s.mediaService.UploadFiles(ctx, uploadPrefix, input.Files)
		}
	}

	var uploadVoiceFn UploadVoiceFn
	if input.Voice != nil {
		uploadVoiceFn = func(m *ent.Message) (*media.UploadInfo, error) {
			uploadPrefix := fmt.Sprintf("rooms/%s/%s/voice", rm.ID, m.ID)
			return s.mediaService.UploadFile(ctx, uploadPrefix, input.Voice)
		}
	}

	msg, err := s.chatRepository.CreateMessage(
		ctx,
		rm.ID,
		currentUserID,
		input.ReplyTo,
		input.Content,
		uploadAttachmentsFn,
		uploadVoiceFn,
	)
	if err != nil {
		return nil, err
	}

	natsClient := s.natsService.Client()
	subject := fmt.Sprintf("room.%s.message.send", rm.ID)

	if err := natsClient.Publish(subject, msg); err != nil {
		return msg, err
	}

	return msg, nil
}

func (s *service) SubscribeToMessageAddedEvent(roomID ID) (<-chan *Message, error) {
	subject := fmt.Sprintf("room.%s.message.send", roomID)

	return s.subscribe(subject)
}

func (s *service) UpdateMessage(
	ctx context.Context,
	messageID ID,
	input UpdateMessageInput,
) (*Message, error) {
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
		input.Content,
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

func (s *service) SubscribeToMessageUpdatedEvent(
	roomID ID,
) (<-chan *Message, error) {
	subject := fmt.Sprintf("room.%s.message.update", roomID)

	return s.subscribe(subject)
}

func (s *service) DeleteMessage(
	ctx context.Context,
	messageID ID,
) (*Message, error) {
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

func (s *service) SubscribeToMessageDeletedEvent(
	roomID ID,
) (<-chan *Message, error) {
	subject := fmt.Sprintf("room.%s.message.delete", roomID)

	return s.subscribe(subject)
}

func (s *service) subscribe(
	subject string,
) (<-chan *Message, error) {
	natsClient := s.natsService.Client()
	ch := make(chan *Message)

	_, err := natsClient.Subscribe(subject, func(msg *Message) {
		ch <- msg
	})
	if err != nil {
		return ch, err
	}

	return ch, nil
}
