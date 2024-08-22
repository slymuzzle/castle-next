package chat

import (
	"context"
	"fmt"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/graph/model"
	"journeyhub/internal/modules/auth"
	"journeyhub/internal/modules/rooms"
	"journeyhub/internal/platform/media"
	"journeyhub/internal/platform/nats"
)

type (
	SendMessageInput   = model.SendMessageInput
	UpdateMessageInput = model.UpdateMessageInput
)

type Service interface {
	SendMessage(
		ctx context.Context,
		input SendMessageInput,
	) (*ent.MessageEdge, error)

	SubscribeToMessageAddedEvent(
		ctx context.Context,
		roomID pulid.ID,
	) (<-chan *ent.MessageEdge, error)

	UpdateMessage(
		ctx context.Context,
		messageID pulid.ID,
		input UpdateMessageInput,
	) (*ent.MessageEdge, error)

	SubscribeToMessageUpdatedEvent(
		ctx context.Context,
		roomID pulid.ID,
	) (<-chan *ent.MessageEdge, error)

	DeleteMessage(
		ctx context.Context,
		messageID pulid.ID,
	) (*ent.MessageEdge, error)

	SubscribeToMessageDeletedEvent(
		ctx context.Context,
		roomID pulid.ID,
	) (<-chan *ent.MessageEdge, error)
}

type service struct {
	chatRepository Repository
	authService    auth.Service
	roomsService   rooms.Service
	natsService    nats.Service
	mediaService   media.Service
}

func NewService(
	chatRepository Repository,
	authService auth.Service,
	roomsService rooms.Service,
	natsService nats.Service,
	mediaService media.Service,
) Service {
	return &service{
		chatRepository: chatRepository,
		authService:    authService,
		roomsService:   roomsService,
		natsService:    natsService,
		mediaService:   mediaService,
	}
}

func (s *service) SendMessage(
	ctx context.Context,
	input SendMessageInput,
) (*ent.MessageEdge, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	rm, err := s.roomsService.FindOrCreatePersonalRoom(
		ctx,
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

	_, err = s.roomsService.UpdateRoom(ctx, rm.ID, rooms.UpdateRoomInput{
		AddVersion:    1,
		LastMessageID: &msg.ID,
	})
	if err != nil {
		return nil, err
	}

	msgEdge := msg.ToEdge(ent.DefaultMessageOrder)

	natsClient := s.natsService.Client()

	subject := fmt.Sprintf("room.%s.message.created", rm.ID)
	if err := natsClient.Publish(subject, msgEdge); err != nil {
		return msgEdge, err
	}

	return msgEdge, nil
}

func (s *service) SubscribeToMessageAddedEvent(
	ctx context.Context,
	roomID pulid.ID,
) (<-chan *ent.MessageEdge, error) {
	subject := fmt.Sprintf("room.%s.message.created", roomID)

	return s.subscribe(ctx, subject)
}

func (s *service) UpdateMessage(
	ctx context.Context,
	messageID pulid.ID,
	input UpdateMessageInput,
) (*ent.MessageEdge, error) {
	_, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	rm, err := s.roomsService.FindRoomByMessage(
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

	_, err = s.roomsService.UpdateRoom(ctx, rm.ID, rooms.UpdateRoomInput{
		AddVersion: 1,
	})
	if err != nil {
		return nil, err
	}

	msgEdge := msg.ToEdge(ent.DefaultMessageOrder)

	natsClient := s.natsService.Client()

	subject := fmt.Sprintf("room.%s.message.updated", rm.ID)
	if err := natsClient.Publish(subject, msgEdge); err != nil {
		return msgEdge, err
	}

	return msgEdge, nil
}

func (s *service) SubscribeToMessageUpdatedEvent(
	ctx context.Context,
	roomID pulid.ID,
) (<-chan *ent.MessageEdge, error) {
	subject := fmt.Sprintf("room.%s.message.updated", roomID)

	return s.subscribe(ctx, subject)
}

func (s *service) DeleteMessage(
	ctx context.Context,
	messageID pulid.ID,
) (*ent.MessageEdge, error) {
	_, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	rm, err := s.roomsService.FindRoomByMessage(
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

	_, err = s.roomsService.UpdateRoom(ctx, rm.ID, rooms.UpdateRoomInput{
		AddVersion: 1,
	})
	if err != nil {
		return nil, err
	}

	msgEdge := msg.ToEdge(ent.DefaultMessageOrder)

	natsClient := s.natsService.Client()

	subject := fmt.Sprintf("room.%s.message.deleted", rm.ID)
	if err := natsClient.Publish(subject, msgEdge); err != nil {
		return msgEdge, err
	}

	return msgEdge, nil
}

func (s *service) SubscribeToMessageDeletedEvent(
	ctx context.Context,
	roomID pulid.ID,
) (<-chan *ent.MessageEdge, error) {
	subject := fmt.Sprintf("room.%s.message.deleted", roomID)

	return s.subscribe(ctx, subject)
}

func (s *service) subscribe(
	ctx context.Context,
	subject string,
) (<-chan *ent.MessageEdge, error) {
	natsClient := s.natsService.Client()

	ch := make(chan *ent.MessageEdge, 1)

	sub, err := natsClient.Subscribe(subject, func(msg *ent.MessageEdge) {
		ch <- msg
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
