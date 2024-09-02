package chat

import (
	"context"
	"fmt"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/graph/model"
	"journeyhub/internal/modules/auth"
	"journeyhub/internal/modules/media"
	"journeyhub/internal/modules/roommembers"
	"journeyhub/internal/modules/rooms"
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
	) (*ent.Message, error)

	SubscribeToMessageCreatedEvent(
		ctx context.Context,
		roomID pulid.ID,
	) (<-chan *ent.MessageEdge, error)

	UpdateMessage(
		ctx context.Context,
		messageID pulid.ID,
		input UpdateMessageInput,
	) (*ent.Message, error)

	SubscribeToMessageUpdatedEvent(
		ctx context.Context,
		roomID pulid.ID,
	) (<-chan *ent.MessageEdge, error)

	DeleteMessage(
		ctx context.Context,
		messageID pulid.ID,
	) (*ent.Message, error)

	SubscribeToMessageDeletedEvent(
		ctx context.Context,
		roomID pulid.ID,
	) (<-chan pulid.ID, error)
}

type service struct {
	chatRepository     Repository
	authService        auth.Service
	roomsService       rooms.Service
	roomMembersService roommembers.Service
	natsService        nats.Service
	mediaService       media.Service
}

func NewService(
	chatRepository Repository,
	authService auth.Service,
	roomsService rooms.Service,
	roomMembersService roommembers.Service,
	natsService nats.Service,
	mediaService media.Service,
) Service {
	return &service{
		chatRepository:     chatRepository,
		authService:        authService,
		roomsService:       roomsService,
		roomMembersService: roomMembersService,
		natsService:        natsService,
		mediaService:       mediaService,
	}
}

func (s *service) SendMessage(
	ctx context.Context,
	input SendMessageInput,
) (*ent.Message, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	_, err = s.roomMembersService.RestoreRoomMembersByRoom(ctx, input.RoomID)
	if err != nil {
		return nil, err
	}

	var uploadAttachmentsFn UploadAttachmentsFn
	if input.Files != nil && len(input.Files) > 0 {
		uploadAttachmentsFn = func(m *ent.Message) ([]*media.UploadInfo, error) {
			uploadPrefix := fmt.Sprintf("rooms/%s/%s/attachments", input.RoomID, m.ID)
			return s.mediaService.UploadMessageFiles(ctx, uploadPrefix, input.Files)
		}
	}

	var uploadVoiceFn UploadVoiceFn
	if input.Voice != nil {
		uploadVoiceFn = func(m *ent.Message) (*media.UploadInfo, error) {
			uploadPrefix := fmt.Sprintf("rooms/%s/%s/voice", input.RoomID, m.ID)
			return s.mediaService.UploadFile(ctx, uploadPrefix, input.Voice)
		}
	}

	msg, err := s.chatRepository.CreateMessage(
		ctx,
		input.RoomID,
		currentUserID,
		input.ReplyTo,
		input.Content,
		uploadAttachmentsFn,
		uploadVoiceFn,
	)
	if err != nil {
		return nil, err
	}

	_, err = s.roomMembersService.IncrementUnreadMessagesCount(ctx, input.RoomID)
	if err != nil {
		return nil, err
	}

	_, err = s.roomsService.IncrementRoomVersion(ctx, input.RoomID, msg)
	if err != nil {
		return nil, err
	}

	natsClient := s.natsService.Client()
	subject := fmt.Sprintf("room.%s.message.created", input.RoomID)
	if err := natsClient.Publish(subject, msg.ID); err != nil {
		return msg, err
	}

	return msg, nil
}

func (s *service) SubscribeToMessageCreatedEvent(
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
) (*ent.Message, error) {
	_, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	room, err := s.roomsService.FindRoomByMessage(ctx, messageID)
	if err != nil {
		return nil, err
	}

	message, err := s.chatRepository.UpdateMessage(ctx, messageID, input.Content)
	if err != nil {
		return nil, err
	}

	_, err = s.roomsService.IncrementRoomVersion(ctx, room.ID, message)
	if err != nil {
		return nil, err
	}

	natsClient := s.natsService.Client()
	subject := fmt.Sprintf("room.%s.message.updated", room.ID)
	if err := natsClient.Publish(subject, message.ID); err != nil {
		return message, err
	}

	return message, nil
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
) (*ent.Message, error) {
	_, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	room, err := s.roomsService.FindRoomByMessage(ctx, messageID)
	if err != nil {
		return nil, err
	}

	message, err := s.chatRepository.DeleteMessage(ctx, messageID)
	if err != nil {
		return nil, err
	}

	_, err = s.roomsService.IncrementRoomVersion(ctx, room.ID, nil)
	if err != nil {
		return nil, err
	}

	natsClient := s.natsService.Client()
	subject := fmt.Sprintf("room.%s.message.deleted", room.ID)
	if err := natsClient.Publish(subject, message.ID); err != nil {
		return message, err
	}

	return message, nil
}

func (s *service) SubscribeToMessageDeletedEvent(
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

func (s *service) subscribe(
	ctx context.Context,
	subject string,
) (<-chan *ent.MessageEdge, error) {
	natsClient := s.natsService.Client()

	ch := make(chan *ent.MessageEdge, 1)

	sub, err := natsClient.Subscribe(subject, func(messageID pulid.ID) {
		message, err := s.chatRepository.FindByID(ctx, messageID)
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
