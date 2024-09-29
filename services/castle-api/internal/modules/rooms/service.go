package rooms

import (
	"context"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/modules/auth"
	"journeyhub/internal/modules/roommembers"
)

type Service interface {
	// FindOrCreatePersonalRoom(
	// 	ctx context.Context,
	// 	targetUserID pulid.ID,
	// ) (*ent.Room, error)

	IncrementRoomVersion(
		ctx context.Context,
		ID pulid.ID,
		lastMessage *ent.Message,
	) (*ent.Room, error)

	DeleteRoom(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.Room, error)
}

type service struct {
	entClient          *ent.Client
	authService        auth.Service
	roomMembersService roommembers.Service
}

func NewService(
	entClient *ent.Client,
	roomMembersService roommembers.Service,
	authService auth.Service,
) Service {
	return &service{
		entClient:          entClient,
		roomMembersService: roomMembersService,
		authService:        authService,
	}
}

// func (s *service) FindOrCreatePersonalRoom(
// 	ctx context.Context,
// 	targetUserID pulid.ID,
// ) (*ent.Room, error) {
// 	currentUserID, err := s.authService.Auth(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	room, err := s.roomsRepository.FindPersonal(ctx, currentUserID, targetUserID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if room != nil {
// 		_, err = s.roomMembersService.RestoreRoomMemberByRoom(ctx, targetUserID, room.ID)
// 		if !ent.IsNotFound(err) && err != nil {
// 			return nil, err
// 		}
// 		return room, nil
// 	}
//
// 	room, err = s.roomsRepository.CreatePersonal(ctx, currentUserID, targetUserID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	s.roomMembersService.CreateRoomMembers(ctx, []roommembers.CreateRoomMemberInput{
// 		{
// 			UserID: currentUserID,
// 			RoomID: room.ID,
// 		},
// 		{
// 			UserID: targetUserID,
// 			RoomID: room.ID,
// 		},
// 	})
//
// 	return room, nil
// }

func (s *service) IncrementRoomVersion(
	ctx context.Context,
	ID pulid.ID,
	lastMessage *ent.Message,
) (*ent.Room, error) {
	_, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	var lastMessageID *pulid.ID
	if lastMessage == nil {
		lastMessageID = nil
	} else {
		lastMessageID = &lastMessage.ID
	}

	repository := s.entClient

	room, err := repository.Room.
		UpdateOneID(ID).
		AddVersion(1).
		SetNillableLastMessageID(lastMessageID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *service) DeleteRoom(
	ctx context.Context,
	ID pulid.ID,
) (*ent.Room, error) {
	_, err := s.authService.Auth(ctx)
	if err != nil {
		return nil, err
	}

	repository := s.entClient

	room, err := repository.Room.Get(ctx, ID)
	if err != nil {
		return nil, err
	}

	err = repository.Room.
		DeleteOneID(ID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return room, nil
}
