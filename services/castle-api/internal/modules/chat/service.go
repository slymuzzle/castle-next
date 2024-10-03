package chat

import (
	"context"
	"fmt"

	"journeyhub/ent"
	"journeyhub/ent/message"
	"journeyhub/ent/room"
	"journeyhub/ent/schema/pulid"
	"journeyhub/graph/model"
	"journeyhub/internal/modules/auth"
	"journeyhub/internal/modules/media"
	"journeyhub/internal/modules/notifications"
	"journeyhub/internal/modules/roommembers"
	"journeyhub/internal/modules/rooms"
)

type Service interface {
	SendMessage(
		ctx context.Context,
		input model.SendMessageInput,
	) (*ent.Message, error)

	UpdateMessage(
		ctx context.Context,
		messageID pulid.ID,
		input model.UpdateMessageInput,
	) (*ent.Message, error)

	DeleteMessage(
		ctx context.Context,
		messageID pulid.ID,
	) (*ent.Message, error)

	Subscriptions() Subscriptions
}

type service struct {
	entClient           *ent.Client
	subscriptions       Subscriptions
	authService         auth.Service
	roomsService        rooms.Service
	roomMembersService  roommembers.Service
	mediaService        media.Service
	notificationService notifications.Service
}

func NewService(
	entClient *ent.Client,
	subscriptions Subscriptions,
	authService auth.Service,
	roomsService rooms.Service,
	roomMembersService roommembers.Service,
	mediaService media.Service,
) Service {
	return &service{
		entClient:          entClient,
		subscriptions:      subscriptions,
		authService:        authService,
		roomsService:       roomsService,
		roomMembersService: roomMembersService,
		mediaService:       mediaService,
	}
}

func (s *service) SendMessage(
	ctx context.Context,
	input model.SendMessageInput,
) (*ent.Message, error) {
	currentUserID, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	_, err = s.roomMembersService.RestoreRoomMembersByRoom(ctx, input.RoomID)
	if err != nil {
		return nil, err
	}

	repository := s.entClient

	message, err := repository.Message.
		Create().
		SetRoomID(input.RoomID).
		SetUserID(currentUserID).
		SetNillableReplyToID(input.ReplyTo).
		SetNillableContent(input.Content).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	err = repository.MessageLink.MapCreateBulk(
		input.Links,
		func(a *ent.MessageLinkCreate, i int) {
			a.SetLink(input.Links[i].Link).
				SetNillableTitle(input.Links[i].Title).
				SetNillableDescription(input.Links[i].Description).
				SetNillableImageURL(input.Links[i].ImageURL).
				SetMessage(message).
				SetRoomID(input.RoomID)
		},
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	if input.Files != nil && len(input.Files) > 0 {
		uploadPrefix := fmt.Sprintf("rooms/%s/%s/attachments", input.RoomID, message.ID)
		uAtchInfo, uErr := s.mediaService.UploadMessageFiles(ctx, uploadPrefix, input.Files)
		if uErr != nil {
			return nil, uErr
		}

		msgFiles, uErr := repository.File.MapCreateBulk(
			uAtchInfo,
			func(a *ent.FileCreate, i int) {
				a.SetID(uAtchInfo[i].ID).
					SetName(uAtchInfo[i].Filename).
					SetContentType(uAtchInfo[i].ContentType).
					SetSize(uint64(uAtchInfo[i].Size)).
					SetLocation(uAtchInfo[i].Location).
					SetBucket(uAtchInfo[i].Bucket).
					SetPath(uAtchInfo[i].Path)
			},
		).Save(ctx)
		if uErr != nil {
			return nil, uErr
		}

		uErr = repository.MessageAttachment.MapCreateBulk(
			msgFiles,
			func(a *ent.MessageAttachmentCreate, i int) {
				a.SetMessage(message).
					SetType(uAtchInfo[i].Type).
					SetRoomID(input.RoomID).
					SetFile(msgFiles[i]).
					SetOrder(uint(i))
			},
		).Exec(ctx)
		if uErr != nil {
			return nil, uErr
		}
	}

	if input.Voice != nil {
		uploadPrefix := fmt.Sprintf("rooms/%s/%s/voice", input.RoomID, message.ID)
		uVoiceInfo, uErr := s.mediaService.UploadFile(ctx, uploadPrefix, input.Voice)
		if uErr != nil {
			return nil, uErr
		}

		voiceFile, uErr := repository.File.
			Create().
			SetID(uVoiceInfo.ID).
			SetName(uVoiceInfo.Filename).
			SetContentType(uVoiceInfo.ContentType).
			SetSize(uint64(uVoiceInfo.Size)).
			SetLocation(uVoiceInfo.Location).
			SetBucket(uVoiceInfo.Bucket).
			SetPath(uVoiceInfo.Path).
			Save(ctx)
		if uErr != nil {
			return nil, uErr
		}

		uErr = repository.MessageVoice.
			Create().
			SetMessage(message).
			SetRoomID(input.RoomID).
			SetFile(voiceFile).
			Exec(ctx)
		if uErr != nil {
			return nil, uErr
		}
	}

	_, err = s.roomMembersService.IncrementUnreadMessagesCount(ctx, input.RoomID)
	if err != nil {
		return nil, err
	}

	_, err = s.roomsService.IncrementRoomVersion(ctx, input.RoomID, message)
	if err != nil {
		return nil, err
	}

	_, err = s.subscriptions.PublishMessageCreatedEvent(ctx, input.RoomID, message.ID)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *service) UpdateMessage(
	ctx context.Context,
	messageID pulid.ID,
	input model.UpdateMessageInput,
) (*ent.Message, error) {
	_, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	repository := s.entClient

	room, err := repository.Room.
		Query().
		Where(room.HasMessagesWith(message.ID(messageID))).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	message, err := repository.Message.
		UpdateOneID(messageID).
		SetContent(input.Content).
		Save(ctx)
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

	repository := s.entClient

	room, err := repository.Room.
		Query().
		Where(room.HasMessagesWith(message.ID(messageID))).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	message, err := repository.Message.Get(ctx, messageID)
	if err != nil {
		return nil, err
	}

	err = repository.Message.DeleteOneID(messageID).Exec(ctx)
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
