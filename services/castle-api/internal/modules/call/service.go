package call

import (
	"context"
	"fmt"
	"time"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/ent/usercontact"
	"journeyhub/internal/modules/auth"
	"journeyhub/internal/platform/config"
	"journeyhub/internal/platform/notification"

	"github.com/appleboy/gorush/rpc/proto"

	livekitauth "github.com/livekit/protocol/auth"
	"google.golang.org/protobuf/types/known/structpb"
)

type CallNotificationType int

const (
	Start CallNotificationType = iota
	End
	Decline
)

type Service interface {
	StartCall(
		ctx context.Context,
		roomID pulid.ID,
	) (bool, error)

	EndCall(
		ctx context.Context,
		roomID pulid.ID,
	) (bool, error)

	DeclineCall(
		ctx context.Context,
		roomID pulid.ID,
	) (bool, error)

	GetCallJoinToken(
		ctx context.Context,
		roomID pulid.ID,
	) (string, error)
}

type service struct {
	config              config.LivekitConfig
	entClient           *ent.Client
	authService         auth.Service
	notificationService notification.Service
}

func NewService(
	config config.LivekitConfig,
	entClient *ent.Client,
	authService auth.Service,
	notificationService notification.Service,
) Service {
	return &service{
		config:              config,
		entClient:           entClient,
		authService:         authService,
		notificationService: notificationService,
	}
}

func (s *service) StartCall(
	ctx context.Context,
	roomID pulid.ID,
) (bool, error) {
	currentUser, err := s.authService.AuthUser(ctx)
	if err != nil {
		return false, err
	}

	userContact, err := s.getInversedUserContact(ctx, roomID, currentUser.ID)
	if err != nil {
		return false, err
	}

	_, err = s.notificationService.Client().Send(ctx, &proto.NotificationRequest{
		Platform: 2,
		Tokens:   []string{userContact.Edges.User.Edges.Device.FcmToken},
		Title:    "Test Title",
		Message:  "Test Message",
		Priority: proto.NotificationRequest_HIGH,
		Data:     s.getCallData(Start, userContact),
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *service) EndCall(
	ctx context.Context,
	roomID pulid.ID,
) (bool, error) {
	currentUser, err := s.authService.AuthUser(ctx)
	if err != nil {
		return false, err
	}

	userContact, err := s.getInversedUserContact(ctx, roomID, currentUser.ID)
	if err != nil {
		return false, err
	}

	_, err = s.notificationService.Client().Send(ctx, &proto.NotificationRequest{
		Platform: 2,
		Tokens:   []string{userContact.Edges.User.Edges.Device.FcmToken},
		Title:    "Test Title",
		Message:  "Test Message",
		Priority: proto.NotificationRequest_HIGH,
		Data:     s.getCallData(End, userContact),
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *service) DeclineCall(
	ctx context.Context,
	roomID pulid.ID,
) (bool, error) {
	currentUser, err := s.authService.AuthUser(ctx)
	if err != nil {
		return false, err
	}

	userContact, err := s.getInversedUserContact(ctx, roomID, currentUser.ID)
	if err != nil {
		return false, err
	}

	_, err = s.notificationService.Client().Send(ctx, &proto.NotificationRequest{
		Platform: 2,
		Tokens:   []string{userContact.Edges.User.Edges.Device.FcmToken},
		Title:    "Test Title",
		Message:  "Test Message",
		Priority: proto.NotificationRequest_HIGH,
		Data:     s.getCallData(Decline, userContact),
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *service) GetCallJoinToken(
	ctx context.Context,
	roomID pulid.ID,
) (string, error) {
	user, err := s.authService.AuthUser(ctx)
	if err != nil {
		return "", err
	}

	at := livekitauth.NewAccessToken(s.config.Access, s.config.Secret)
	grant := &livekitauth.VideoGrant{
		RoomJoin: true,
		Room:     string(roomID),
	}
	at.AddGrant(grant).
		SetIdentity(string(user.ID)).
		SetName(fmt.Sprintf("%s %s", user.FirstName, user.LastName)).
		SetValidFor(time.Hour)

	return at.ToJWT()
}

func (s *service) getInversedUserContact(
	ctx context.Context,
	roomID pulid.ID,
	currentUserID pulid.ID,
) (*ent.UserContact, error) {
	return s.entClient.UserContact.
		Query().
		WithUser(func(q *ent.UserQuery) {
			q.WithDevice()
		}).
		Where(
			usercontact.RoomID(roomID),
			usercontact.ContactID(currentUserID),
		).
		Only(ctx)
}

func (s *service) getCallData(
	callType CallNotificationType,
	userContact *ent.UserContact,
) *structpb.Struct {
	return &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"type": {
				Kind: &structpb.Value_NumberValue{NumberValue: float64(callType)},
			},
			"contactID": {
				Kind: &structpb.Value_StringValue{StringValue: string(userContact.ID)},
			},
			"roomID": {
				Kind: &structpb.Value_StringValue{StringValue: string(userContact.RoomID)},
			},
			"userID": {
				Kind: &structpb.Value_StringValue{StringValue: string(userContact.Edges.User.ID)},
			},
			"userFirstName": {
				Kind: &structpb.Value_StringValue{StringValue: string(userContact.Edges.User.FirstName)},
			},
			"userLastName": {
				Kind: &structpb.Value_StringValue{StringValue: string(userContact.Edges.User.LastName)},
			},
		},
	}
}
