package rooms

import (
	"context"
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
	// FindOrCreatePersonal(
	// 	ctx context.Context,
	// 	currentUserID pulid.ID,
	// 	targetUserID pulid.ID,
	// ) (*ent.Room, error)

	FindByMessage(
		ctx context.Context,
		messageID pulid.ID,
	) (*ent.Room, error)

	FindPersonal(
		ctx context.Context,
		currentUserID pulid.ID,
		targetUserID pulid.ID,
	) (*ent.Room, error)

	CreatePersonal(
		ctx context.Context,
		currentUserID pulid.ID,
		targetUserID pulid.ID,
	) (*ent.Room, error)

	IncrementVersion(
		ctx context.Context,
		ID pulid.ID,
		lastMessageID *pulid.ID,
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

// func (r *repository) FindOrCreatePersonal(
// 	ctx context.Context,
// 	currentUserID pulid.ID,
// 	targetUserID pulid.ID,
// ) (*ent.Room, error) {
// 	client := r.getClient(ctx)
//
// 	uc, err := client.UserContact.
// 		Query().
// 		Where(
// 			usercontact.UserID(currentUserID),
// 			usercontact.ContactID(targetUserID),
// 		).
// 		WithRoom().
// 		Only(ctx)
// 	if !ent.IsNotFound(err) {
// 		rm := uc.Edges.Room
//
// 		err = client.RoomMember.
// 			Update().
// 			Where(
// 				roommember.And(
// 					roommember.RoomID(rm.ID),
// 					roommember.DeletedAtNotNil(),
// 				),
// 			).
// 			ClearDeletedAt().
// 			Exec(ctx)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		if rm != nil {
// 			return rm, nil
// 		}
// 	}
//
// 	currUsr, err := client.User.Get(ctx, currentUserID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	rmName := strings.TrimRightFunc(
// 		fmt.Sprintf("%s %s", currUsr.FirstName, currUsr.LastName),
// 		unicode.IsSpace,
// 	)
//
// 	rm, err := client.Room.
// 		Create().
// 		SetName(rmName).
// 		SetType(room.TypePersonal).
// 		Save(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// FIX: Wait for https://github.com/ent/ent/pull/4170
// 	err = client.RoomMember.CreateBulk(
// 		client.RoomMember.
// 			Create().
// 			SetUserID(currentUserID).
// 			SetRoomID(rm.ID),
// 		client.RoomMember.
// 			Create().
// 			SetUserID(targetUserID).
// 			SetRoomID(rm.ID),
// 	).Exec(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	err = client.UserContact.
// 		Create().
// 		SetUserID(currentUserID).
// 		SetContactID(targetUserID).
// 		SetRoomID(rm.ID).
// 		OnConflict().
// 		DoNothing().
// 		Exec(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return rm, nil
// }

func (r *repository) FindByMessage(
	ctx context.Context,
	messageID pulid.ID,
) (*ent.Room, error) {
	return r.getClient(ctx).Room.
		Query().
		Where(
			room.HasMessagesWith(
				message.ID(messageID),
			),
		).
		Only(ctx)
}

func (r *repository) FindPersonal(
	ctx context.Context,
	currentUserID pulid.ID,
	targetUserID pulid.ID,
) (*ent.Room, error) {
	client := r.getClient(ctx)

	userContact, err := client.UserContact.
		Query().
		Where(
			usercontact.UserID(currentUserID),
			usercontact.ContactID(targetUserID),
		).
		WithRoom().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return userContact.Room(ctx)
}

func (r *repository) CreatePersonal(
	ctx context.Context,
	currentUserID pulid.ID,
	targetUserID pulid.ID,
) (*ent.Room, error) {
	client := r.getClient(ctx)

	currUsr, err := client.User.Get(ctx, currentUserID)
	if err != nil {
		return nil, err
	}

	rmName := strings.TrimRightFunc(
		fmt.Sprintf("%s %s", currUsr.FirstName, currUsr.LastName),
		unicode.IsSpace,
	)

	rm, err := client.Room.
		Create().
		SetName(rmName).
		SetType(room.TypePersonal).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	err = client.UserContact.
		Update().
		Where(
			usercontact.Or(
				usercontact.And(
					usercontact.UserID(currentUserID),
					usercontact.ContactID(targetUserID),
				),
				usercontact.And(
					usercontact.UserID(targetUserID),
					usercontact.ContactID(currentUserID),
				),
			),
		).
		SetRoomID(rm.ID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return rm, nil
}

func (r *repository) IncrementVersion(
	ctx context.Context,
	ID pulid.ID,
	lastMessageID *pulid.ID,
) (*ent.Room, error) {
	client := r.getClient(ctx)

	room, err := client.Room.
		UpdateOneID(ID).
		AddVersion(1).
		SetNillableLastMessageID(lastMessageID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *repository) Delete(
	ctx context.Context,
	ID pulid.ID,
) (*ent.Room, error) {
	client := r.getClient(ctx)

	room, err := client.Room.Get(ctx, ID)
	if err != nil {
		return nil, err
	}

	err = client.Room.
		DeleteOneID(ID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return room, nil
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
