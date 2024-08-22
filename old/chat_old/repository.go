package chat

import (
	"context"
	"errors"
	"fmt"
	"journeyhub/ent"
	"journeyhub/ent/message"
	"journeyhub/ent/room"
	"journeyhub/ent/roommember"
	"journeyhub/ent/usercontact"
	"journeyhub/internal/media"
	"strings"
	"unicode"
)

type (
	UploadAttachmentsFn func(*ent.Message) ([]*media.UploadInfo, error)
	UploadVoiceFn       func(*ent.Message) (*media.UploadInfo, error)
)

type Repository interface {
	FindRoomByID(
		ctx context.Context,
		roomID ID,
	) (*Room, error)
	FindRoomByMessage(
		ctx context.Context,
		messageID ID,
	) (*Room, error)
	FindOrCreatePersonalRoom(
		ctx context.Context,
		currentUserID ID,
		targetUserID ID,
	) (*Room, error)
	CreateMessage(
		ctx context.Context,
		roomID ID,
		currentUserID ID,
		replyTo *ID,
		content string,
		uploadAttachmentsFn UploadAttachmentsFn,
		uploadVoiceFn UploadVoiceFn,
	) (*Message, error)
	UpdateMessage(
		ctx context.Context,
		messageID ID,
		content string,
	) (*Message, error)
	DeleteMessage(
		ctx context.Context,
		messageID ID,
	) (*Message, error)
}

type repository struct {
	entClient *ent.Client
}

func NewRepository(entClient *ent.Client) Repository {
	return &repository{
		entClient: entClient,
	}
}

func (r *repository) FindRoomByID(
	ctx context.Context,
	roomID ID,
) (*Room, error) {
	return r.entClient.Room.
		Query().
		WithLastMessage().
		Where(
			room.ID(roomID),
		).
		Only(ctx)
}

func (r *repository) FindRoomByMessage(
	ctx context.Context,
	messageID ID,
) (*Room, error) {
	return r.entClient.Room.
		Query().
		Where(
			room.HasMessagesWith(
				message.ID(messageID),
			),
		).
		Only(ctx)
}

func (r *repository) FindOrCreatePersonalRoom(
	ctx context.Context,
	currentUserID ID,
	targetUserID ID,
) (*Room, error) {
	uc, err := r.entClient.UserContact.
		Query().
		Where(
			usercontact.UserID(currentUserID),
			usercontact.ContactID(targetUserID),
		).
		WithRoom().
		Only(ctx)
	if !ent.IsNotFound(err) {
		rm := uc.Edges.Room
		if rm != nil {
			return rm, nil
		}
	}

	currUsr, err := r.entClient.User.Get(ctx, currentUserID)
	if err != nil {
		return nil, err
	}

	rmName := strings.TrimRightFunc(
		fmt.Sprintf("%s %s", currUsr.FirstName, currUsr.LastName),
		unicode.IsSpace,
	)

	tx, err := r.entClient.Tx(ctx)
	if err != nil {
		return nil, err
	}

	rm, err := tx.Room.
		Create().
		SetName(rmName).
		SetType(room.TypePersonal).
		Save(ctx)
	if err != nil {
		return nil, errors.Join(tx.Rollback(), err)
	}

	// FIX: Wait for https://github.com/ent/ent/pull/4170
	err = tx.RoomMember.CreateBulk(
		tx.RoomMember.
			Create().
			SetUserID(currentUserID).
			SetRoomID(rm.ID),
		tx.RoomMember.
			Create().
			SetUserID(targetUserID).
			SetRoomID(rm.ID),
	).Exec(ctx)
	if err != nil {
		return nil, errors.Join(tx.Rollback(), err)
	}

	err = tx.UserContact.
		Create().
		SetUserID(currentUserID).
		SetContactID(targetUserID).
		SetRoomID(rm.ID).
		OnConflict().
		DoNothing().
		Exec(ctx)
	if err != nil {
		return nil, errors.Join(tx.Rollback(), err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return rm, nil
}

func (r *repository) CreateMessage(
	ctx context.Context,
	roomID ID,
	currentUserID ID,
	replyTo *ID,
	content string,
	uploadAttachmentsFn UploadAttachmentsFn,
	uploadVoiceFn UploadVoiceFn,
) (*Message, error) {
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

	rmMrErr := tx.RoomMember.
		Update().
		Where(
			roommember.And(
				roommember.RoomID(roomID),
				roommember.UserIDNEQ(currentUserID),
			),
		).
		AddUnreadMessagesCount(1).
		Exec(ctx)
	if rmMrErr != nil {
		return nil, errors.Join(tx.Rollback(), rmMrErr)
	}

	rmErr := tx.Room.
		UpdateOneID(roomID).
		SetLastMessage(msg).
		AddVersion(1).
		Exec(ctx)
	if rmErr != nil {
		return nil, errors.Join(tx.Rollback(), rmErr)
	}

	txcErr := tx.Commit()
	if txcErr != nil {
		return nil, txcErr
	}

	return msg, nil
}

func (r *repository) UpdateMessage(
	ctx context.Context,
	messageID ID,
	content string,
) (*Message, error) {
	tx, err := r.entClient.Tx(ctx)
	if err != nil {
		return nil, err
	}

	msg, err := tx.Message.
		UpdateOneID(messageID).
		SetContent(content).
		Save(ctx)
	if err != nil {
		return nil, errors.Join(tx.Rollback(), err)
	}

	err = tx.Room.
		UpdateOneID(msg.Edges.Room.ID).
		AddVersion(1).
		Exec(ctx)
	if err != nil {
		return nil, errors.Join(tx.Rollback(), err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (r *repository) DeleteMessage(
	ctx context.Context,
	messageID ID,
) (*Message, error) {
	msg, err := r.entClient.Message.Get(ctx, messageID)
	if err != nil {
		return nil, err
	}

	tx, err := r.entClient.Tx(ctx)
	if err != nil {
		return nil, err
	}

	err = tx.Message.
		DeleteOneID(messageID).
		Exec(ctx)
	if err != nil {
		return nil, errors.Join(tx.Rollback(), err)
	}

	err = tx.Room.
		UpdateOneID(msg.Edges.Room.ID).
		AddVersion(1).
		Exec(ctx)
	if err != nil {
		return nil, errors.Join(tx.Rollback(), err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return msg, nil
}
