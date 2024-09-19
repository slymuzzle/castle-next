package rooms

import (
	"context"
	"os"
	"time"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/modules/auth"
	"journeyhub/internal/modules/roommembers"
	"journeyhub/internal/platform/nats"

	livekitauth "github.com/livekit/protocol/auth"
)

type Service interface {
	// FindOrCreatePersonalRoom(
	// 	ctx context.Context,
	// 	targetUserID pulid.ID,
	// ) (*ent.Room, error)

	FindRoomByMessage(
		ctx context.Context,
		messageID pulid.ID,
	) (*ent.Room, error)

	IncrementRoomVersion(
		ctx context.Context,
		ID pulid.ID,
		lastMessage *ent.Message,
	) (*ent.Room, error)

	DeleteRoom(
		ctx context.Context,
		ID pulid.ID,
	) (*ent.Room, error)

	GetJoinToken(
		ctx context.Context,
		ID pulid.ID,
	) (string, error)
}

type service struct {
	roomsRepository    Repository
	authService        auth.Service
	roomMembersService roommembers.Service
	natsService        nats.Service
}

func NewService(
	roomsRepository Repository,
	roomMembersService roommembers.Service,
	authService auth.Service,
	natsService nats.Service,
) Service {
	return &service{
		roomsRepository:    roomsRepository,
		roomMembersService: roomMembersService,
		authService:        authService,
		natsService:        natsService,
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

func (s *service) FindRoomByMessage(
	ctx context.Context,
	messageID pulid.ID,
) (*ent.Room, error) {
	return s.roomsRepository.FindByMessage(ctx, messageID)
}

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

	room, err := s.roomsRepository.IncrementVersion(ctx, ID, lastMessageID)
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

	return s.roomsRepository.Delete(ctx, ID)
}

func (s *service) GetJoinToken(
	ctx context.Context,
	ID pulid.ID,
) (string, error) {
	userID, err := s.authService.Auth(ctx)
	if err != nil {
		return "", err
	}

	at := livekitauth.NewAccessToken(os.Getenv("LIVEKIT_API_KEY"), os.Getenv("LIVEKIT_API_SECRET"))
	grant := &livekitauth.VideoGrant{
		RoomJoin: true,
		Room:     string(ID),
	}
	at.AddGrant(grant).
		SetIdentity(string(userID)).
		SetValidFor(time.Hour)

	return at.ToJWT()
}
