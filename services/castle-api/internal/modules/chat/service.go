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

	UpdateMessage(
		ctx context.Context,
		messageID pulid.ID,
		input UpdateMessageInput,
	) (*ent.Message, error)

	DeleteMessage(
		ctx context.Context,
		messageID pulid.ID,
	) (*ent.Message, error)

	Subscriptions() Subscriptions
}

type service struct {
	subscriptions      Subscriptions
	chatRepository     Repository
	authService        auth.Service
	roomsService       rooms.Service
	roomMembersService roommembers.Service
	mediaService       media.Service
}

func NewService(
	subscriptions Subscriptions,
	chatRepository Repository,
	authService auth.Service,
	roomsService rooms.Service,
	roomMembersService roommembers.Service,
	mediaService media.Service,
) Service {
	return &service{
		subscriptions:      subscriptions,
		chatRepository:     chatRepository,
		authService:        authService,
		roomsService:       roomsService,
		roomMembersService: roomMembersService,
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
		uploadAttachmentsFn = func(message *ent.Message) ([]*media.UploadInfo, error) {
			uploadPrefix := fmt.Sprintf("rooms/%s/%s/attachments", input.RoomID, message.ID)
			return s.mediaService.UploadMessageFiles(ctx, uploadPrefix, input.Files)
		}
	}

	var uploadVoiceFn UploadVoiceFn
	if input.Voice != nil {
		uploadVoiceFn = func(message *ent.Message) (*media.UploadInfo, error) {
			uploadPrefix := fmt.Sprintf("rooms/%s/%s/voice", input.RoomID, message.ID)
			return s.mediaService.UploadFile(ctx, uploadPrefix, input.Voice)
		}
	}

	msg, err := s.chatRepository.CreateMessage(
		ctx,
		currentUserID,
		input,
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

	_, err = s.subscriptions.PublishMessageCreatedEvent(ctx, input.RoomID, msg.ID)
	if err != nil {
		return nil, err
	}

	return msg, nil
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

	_, err = s.roomsService.IncrementRoomVersion(ctx, room.ID, nil)
	if err != nil {
		return nil, err
	}

	_, err = s.subscriptions.PublishMessageUpdatedEvent(ctx, room.ID, message.ID)
	if err != nil {
		return nil, err
	}

	return message, nil
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

	_, err = s.subscriptions.PublishMessageDeletedEvent(ctx, room.ID, message.ID)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *service) Subscriptions() Subscriptions {
	return s.subscriptions
}
