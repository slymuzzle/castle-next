package chat

import (
	"context"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/modules/media"
)

type (
	UploadAttachmentsFn func(*ent.Message) ([]*media.UploadInfo, error)
	UploadVoiceFn       func(*ent.Message) (*media.UploadInfo, error)
)

type Repository interface {
	FindByID(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.Message, error)

	CreateMessage(
		ctx context.Context,
		currentUserID pulid.ID,
		input SendMessageInput,
		uploadAttachmentsFn UploadAttachmentsFn,
		uploadVoiceFn UploadVoiceFn,
	) (*ent.Message, error)

	UpdateMessage(
		ctx context.Context,
		messageID pulid.ID,
		content string,
	) (*ent.Message, error)

	DeleteMessage(
		ctx context.Context,
		messageID pulid.ID,
	) (*ent.Message, error)
}

type repository struct {
	entClient *ent.Client
}

func NewRepository(entClient *ent.Client) Repository {
	return &repository{
		entClient: entClient,
	}
}

func (r *repository) FindByID(
	ctx context.Context,
	ID pulid.ID,
) (*ent.Message, error) {
	return r.getClient(ctx).Message.Get(ctx, ID)
}

func (r *repository) CreateMessage(
	ctx context.Context,
	currentUserID pulid.ID,
	input SendMessageInput,
	uploadAttachmentsFn UploadAttachmentsFn,
	uploadVoiceFn UploadVoiceFn,
) (*ent.Message, error) {
	client := r.getClient(ctx)

	msg, err := client.Message.
		Create().
		SetRoomID(input.RoomID).
		SetUserID(currentUserID).
		SetNillableReplyToID(input.ReplyTo).
		SetNillableContent(input.Content).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	err = client.MessageLink.MapCreateBulk(
		input.Links,
		func(a *ent.MessageLinkCreate, i int) {
			a.SetLink(input.Links[i].Link).
				SetNillableTitle(input.Links[i].Title).
				SetNillableDescription(input.Links[i].Description).
				SetNillableImageURL(input.Links[i].ImageURL).
				SetMessage(msg).
				SetRoomID(input.RoomID)
		},
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	if uploadAttachmentsFn != nil {
		uAtchInfo, err := uploadAttachmentsFn(msg)
		if err != nil {
			return nil, err
		}

		msgFiles, err := client.File.MapCreateBulk(
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
		if err != nil {
			return nil, err
		}

		err = client.MessageAttachment.MapCreateBulk(
			msgFiles,
			func(a *ent.MessageAttachmentCreate, i int) {
				a.SetMessage(msg).
					SetType(uAtchInfo[i].Type).
					SetRoomID(input.RoomID).
					SetFile(msgFiles[i]).
					SetOrder(uint(i))
			},
		).Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	if uploadVoiceFn != nil {
		uVoiceInfo, err := uploadVoiceFn(msg)
		if err != nil {
			return nil, err
		}

		voiceFile, err := client.File.
			Create().
			SetID(uVoiceInfo.ID).
			SetName(uVoiceInfo.Filename).
			SetContentType(uVoiceInfo.ContentType).
			SetSize(uint64(uVoiceInfo.Size)).
			SetLocation(uVoiceInfo.Location).
			SetBucket(uVoiceInfo.Bucket).
			SetPath(uVoiceInfo.Path).
			Save(ctx)
		if err != nil {
			return nil, err
		}

		err = client.MessageVoice.
			Create().
			SetMessage(msg).
			SetRoomID(input.RoomID).
			SetFile(voiceFile).
			Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	return msg, nil
}

func (r *repository) UpdateMessage(
	ctx context.Context,
	messageID pulid.ID,
	content string,
) (*ent.Message, error) {
	return r.getClient(ctx).Message.
		UpdateOneID(messageID).
		SetContent(content).
		Save(ctx)
}

func (r *repository) DeleteMessage(
	ctx context.Context,
	messageID pulid.ID,
) (*ent.Message, error) {
	client := r.getClient(ctx)

	msg, err := client.Message.Get(ctx, messageID)
	if err != nil {
		return nil, err
	}

	err = client.Message.
		DeleteOneID(messageID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (r *repository) getClient(ctx context.Context) *ent.Client {
	var client *ent.Client
	if clientFromCtx := ent.FromContext(ctx); clientFromCtx != nil {
		client = clientFromCtx
	} else {
		client = r.entClient
	}
	return client
}
