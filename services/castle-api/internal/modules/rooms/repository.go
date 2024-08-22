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

	Create(
		ctx context.Context,
		input CreateRoomInput,
	) (*ent.Room, error)

	Update(
		ctx context.Context,
		ID pulid.ID,
		input UpdateRoomInput,
	) (*ent.Room, error)

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

func (r *repository) Create(
	ctx context.Context,
	input CreateRoomInput,
) (*ent.Room, error) {
	u := r.entClient.Room.Create()
	input.Mutate(u.Mutation())
	return u.Save(ctx)
}

func (r *repository) Update(
	ctx context.Context,
	ID pulid.ID,
	input UpdateRoomInput,
) (*ent.Room, error) {
	u := r.entClient.Room.UpdateOneID(ID)
	input.Mutate(u.Mutation())
	return u.Save(ctx)
}

func (r *repository) Delete(
	ctx context.Context,
	ID pulid.ID,
) (*ent.Room, error) {
	rm, err := r.entClient.Room.Get(ctx, ID)
	if err != nil {
		return nil, err
	}

	err = r.entClient.Room.
		DeleteOneID(ID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return rm, nil
}
