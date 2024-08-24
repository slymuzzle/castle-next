package chat

import (
	"context"
	"errors"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/platform/media"
)

type (
	UploadAttachmentsFn func(*ent.Message) ([]*media.UploadInfo, error)
	UploadVoiceFn       func(*ent.Message) (*media.UploadInfo, error)
)

type Repository interface {
	CreateMessage(
		ctx context.Context,
		roomID pulid.ID,
		currentUserID pulid.ID,
		replyTo *pulid.ID,
		content string,
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

func (r *repository) CreateMessage(
	ctx context.Context,
	roomID pulid.ID,
	currentUserID pulid.ID,
	replyTo *pulid.ID,
	content string,
	uploadAttachmentsFn UploadAttachmentsFn,
	uploadVoiceFn UploadVoiceFn,
) (*ent.Message, error) {
	tx, txErr := r.entClient.Tx(ctx)
	if txErr != nil {
		return nil, txErr
	}

	msg, msgErr := tx.Message.
		Create().
		SetRoomID(roomID).
		SetUserID(currentUserID).
		SetNillableReplyToID(replyTo).
		SetContent(content).
		Save(ctx)
	if msgErr != nil {
		return nil, errors.Join(tx.Rollback(), msgErr)
	}

	if uploadAttachmentsFn != nil {
		uAtchInfo, uAtchErr := uploadAttachmentsFn(msg)
		if uAtchErr != nil {
			return nil, errors.Join(tx.Rollback(), uAtchErr)
		}

		msgFiles, msgFilesErr := tx.File.MapCreateBulk(
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
			return nil, errors.Join(tx.Rollback(), msgFilesErr)
		}

		msgAtchErr := tx.MessageAttachment.MapCreateBulk(
			msgFiles,
			func(a *ent.MessageAttachmentCreate, i int) {
				a.SetMessage(msg).
					SetRoomID(roomID).
					SetFile(msgFiles[i]).
					SetOrder(uint(i))
			},
		).Exec(ctx)
		if msgAtchErr != nil {
			return nil, errors.Join(tx.Rollback(), msgAtchErr)
		}
	}

	if uploadVoiceFn != nil {
		uVoiceInfo, uVoiceErr := uploadVoiceFn(msg)
		if uVoiceErr != nil {
			return nil, errors.Join(tx.Rollback(), uVoiceErr)
		}

		voiceFile, voiceFileErr := tx.File.
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
			return nil, errors.Join(tx.Rollback(), voiceFileErr)
		}

		msgVoiceErr := tx.MessageVoice.
			Create().
			SetMessage(msg).
			SetRoomID(roomID).
			SetFile(voiceFile).
			Exec(ctx)
		if msgVoiceErr != nil {
			return nil, errors.Join(tx.Rollback(), msgVoiceErr)
		}
	}

	txcErr := tx.Commit()
	if txcErr != nil {
		return nil, txcErr
	}

	return msg.Unwrap(), nil
}

func (r *repository) UpdateMessage(
	ctx context.Context,
	messageID pulid.ID,
	content string,
) (*ent.Message, error) {
	return r.entClient.Message.
		UpdateOneID(messageID).
		SetContent(content).
		Save(ctx)
}

func (r *repository) DeleteMessage(
	ctx context.Context,
	messageID pulid.ID,
) (*ent.Message, error) {
	msg, err := r.entClient.Message.Get(ctx, messageID)
	if err != nil {
		return nil, err
	}

	err = r.entClient.Message.
		DeleteOneID(messageID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
