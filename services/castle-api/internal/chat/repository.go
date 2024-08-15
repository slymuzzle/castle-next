package chat

import (
	"context"
	"errors"
	"fmt"
	"journeyhub/ent"
	"journeyhub/ent/friendship"
	"journeyhub/ent/message"
	"journeyhub/ent/room"
	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/db"
	"strings"
	"unicode"
)

type Repository interface {
	FindRoomByMessage(
		ctx context.Context,
		messageID pulid.ID,
	) (*ent.Room, error)
	FindOrCreatePersonalRoom(
		ctx context.Context,
		currentUserID pulid.ID,
		targetUserID pulid.ID,
	) (*ent.Room, error)
	CreateMessage(
		ctx context.Context,
		roomID pulid.ID,
		currentUserID pulid.ID,
		replyTo pulid.ID,
		content string,
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
	dbService db.Service
}

func NewRepository(dbService db.Service) Repository {
	return &repository{
		dbService: dbService,
	}
}

func (r *repository) FindRoomByMessage(ctx context.Context, messageID pulid.ID) (*ent.Room, error) {
	entClient := r.dbService.Client()

	return entClient.Room.
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
	currentUserID pulid.ID,
	targetUserID pulid.ID,
) (*ent.Room, error) {
	entClient := r.dbService.Client()

	frs, err := entClient.Friendship.
		Query().
		Where(
			friendship.UserID(currentUserID),
			friendship.FriendID(targetUserID),
		).
		WithRoom().
		Only(ctx)
	if !ent.IsNotFound(err) {
		rm := frs.Edges.Room
		if rm != nil {
			return rm, nil
		}
	}

	tx, err := entClient.Tx(ctx)
	if err != nil {
		return nil, err
	}

	currUsr, err := tx.User.Get(ctx, currentUserID)
	if err != nil {
		return nil, err
	}

	rmName := strings.TrimRightFunc(
		fmt.Sprintf("%s %s", currUsr.FirstName, currUsr.LastName),
		unicode.IsSpace,
	)

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

	err = tx.Friendship.
		Create().
		SetUserID(currentUserID).
		SetFriendID(targetUserID).
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
	roomID pulid.ID,
	currentUserID pulid.ID,
	replyTo pulid.ID,
	content string,
) (*ent.Message, error) {
	entClient := r.dbService.Client()

	msg, err := entClient.Message.
		Create().
		SetRoomID(roomID).
		SetUserID(currentUserID).
		SetReplyToID(replyTo).
		SetContent(content).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (r *repository) UpdateMessage(
	ctx context.Context,
	messageID pulid.ID,
	content string,
) (*ent.Message, error) {
	entClient := r.dbService.Client()

	msg, err := entClient.Message.
		UpdateOneID(messageID).
		SetContent(content).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (r *repository) DeleteMessage(
	ctx context.Context,
	messageID pulid.ID,
) (*ent.Message, error) {
	entClient := r.dbService.Client()

	msg, err := entClient.Message.Get(ctx, messageID)
	if err != nil {
		return nil, err
	}

	err = entClient.Message.
		DeleteOneID(messageID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
