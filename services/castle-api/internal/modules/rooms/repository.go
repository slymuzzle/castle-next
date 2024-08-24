package rooms

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"journeyhub/ent"
	"journeyhub/ent/message"
	"journeyhub/ent/room"
	"journeyhub/ent/roommember"
	"journeyhub/ent/schema/pulid"
	"journeyhub/ent/usercontact"
)

type Repository interface {
	FindOrCreatePersonal(
		ctx context.Context,
		currentUserID pulid.ID,
		targetUserID pulid.ID,
	) (*ent.Room, error)

	FindByMessage(
		ctx context.Context,
		messageID pulid.ID,
	) (*ent.Room, error)

	IncrementVersion(
		ctx context.Context,
		ID pulid.ID,
		lastMessageID *pulid.ID,
	) (*ent.Room, error)

	IncrementUnreadMessagesCount(
		ctx context.Context,
		ID pulid.ID,
		currentUserID pulid.ID,
	) error

	DeleteRoomMember(
		ctx context.Context,
		roomMemberID pulid.ID,
	) (*ent.RoomMember, error)

	Delete(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.Room, error)
}

type repository struct {
	entClient *ent.Client
}

func NewRepository(entClient *ent.Client) Repository {
	return &repository{
		entClient: entClient,
	}
}

func (r *repository) FindOrCreatePersonal(
	ctx context.Context,
	currentUserID pulid.ID,
	targetUserID pulid.ID,
) (*ent.Room, error) {
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

	return rm.Unwrap(), nil
}

func (r *repository) FindByMessage(
	ctx context.Context,
	messageID pulid.ID,
) (*ent.Room, error) {
	return r.entClient.Room.
		Query().
		Where(
			room.HasMessagesWith(
				message.ID(messageID),
			),
		).
		Only(ctx)
}

func (r *repository) IncrementVersion(
	ctx context.Context,
	ID pulid.ID,
	lastMessageID *pulid.ID,
) (*ent.Room, error) {
	room, err := r.entClient.Room.
		UpdateOneID(ID).
		AddVersion(1).
		SetNillableLastMessageID(lastMessageID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *repository) IncrementUnreadMessagesCount(
	ctx context.Context,
	ID pulid.ID,
	currentUserID pulid.ID,
) error {
	return r.entClient.RoomMember.
		Update().
		Where(
			roommember.And(
				roommember.RoomID(ID),
				roommember.UserIDNEQ(currentUserID),
			),
		).
		AddUnreadMessagesCount(1).
		Exec(ctx)
}

func (r *repository) DeleteRoomMember(
	ctx context.Context,
	roomMemberID pulid.ID,
) (*ent.RoomMember, error) {
	roomMember, err := r.entClient.RoomMember.Get(ctx, roomMemberID)
	if err != nil {
		return nil, err
	}

	err = r.entClient.RoomMember.
		DeleteOneID(roomMemberID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return roomMember, nil
}

func (r *repository) Delete(
	ctx context.Context,
	ID pulid.ID,
) (*ent.Room, error) {
	room, err := r.entClient.Room.Get(ctx, ID)
	if err != nil {
		return nil, err
	}

	err = r.entClient.Room.
		DeleteOneID(ID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return room, nil
}
