package chat

import (
	"context"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/platform/media"
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
		roomID pulid.ID,
		currentUserID pulid.ID,
		replyTo *pulid.ID,
		content *string,
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
	roomID pulid.ID,
	currentUserID pulid.ID,
	replyTo *pulid.ID,
	content *string,
	uploadAttachmentsFn UploadAttachmentsFn,
	uploadVoiceFn UploadVoiceFn,
) (*ent.Message, error) {
	client := r.getClient(ctx)

	msg, msgErr := client.Message.
		Create().
		SetRoomID(roomID).
		SetUserID(currentUserID).
		SetNillableReplyToID(replyTo).
		SetNillableContent(content).
		Save(ctx)
	if msgErr != nil {
		return nil, msgErr
	}

	if uploadAttachmentsFn != nil {
		uAtchInfo, uAtchErr := uploadAttachmentsFn(msg)
		if uAtchErr != nil {
			return nil, uAtchErr
		}

		msgFiles, msgFilesErr := client.File.MapCreateBulk(
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
		if msgFilesErr != nil {
			return nil, msgFilesErr
		}

		msgAtchErr := client.MessageAttachment.MapCreateBulk(
			msgFiles,
			func(a *ent.MessageAttachmentCreate, i int) {
				a.SetMessage(msg).
					SetRoomID(roomID).
					SetFile(msgFiles[i]).
					SetOrder(uint(i))
			},
		).Exec(ctx)
		if msgAtchErr != nil {
			return nil, msgAtchErr
		}
	}

	if uploadVoiceFn != nil {
		uVoiceInfo, uVoiceErr := uploadVoiceFn(msg)
		if uVoiceErr != nil {
			return nil, uVoiceErr
		}

		voiceFile, voiceFileErr := client.File.
			Create().
			SetID(uVoiceInfo.ID).
			SetName(uVoiceInfo.Filename).
			SetContentType(uVoiceInfo.ContentType).
			SetSize(uint64(uVoiceInfo.Size)).
			SetLocation(uVoiceInfo.Location).
			SetBucket(uVoiceInfo.Bucket).
			SetPath(uVoiceInfo.Path).
			Save(ctx)
		if voiceFileErr != nil {
			return nil, voiceFileErr
		}

		msgVoiceErr := client.MessageVoice.
			Create().
			SetMessage(msg).
			SetRoomID(roomID).
			SetFile(voiceFile).
			Exec(ctx)
		if msgVoiceErr != nil {
			return nil, msgVoiceErr
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
