package calls

import (
	"context"
	"fmt"
	"time"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/ent/usercontact"
	"journeyhub/graph/model"
	"journeyhub/internal/modules/auth"
	"journeyhub/internal/modules/notifications"
	"journeyhub/internal/platform/config"

	"github.com/appleboy/gorush/rpc/proto"

	livekitauth "github.com/livekit/protocol/auth"
	"google.golang.org/protobuf/types/known/structpb"
)

type Service interface {
	StartCall(
		ctx context.Context,
		roomID pulid.ID,
		callType model.CallType,
	) (bool, error)

	EndCall(
		ctx context.Context,
		roomID pulid.ID,
		callType model.CallType,
	) (bool, error)

	DeclineCall(
		ctx context.Context,
		roomID pulid.ID,
		callType model.CallType,
	) (bool, error)

	AnswerCall(
		ctx context.Context,
		roomID pulid.ID,
		callType model.CallType,
	) (bool, error)

	GetCallJoinToken(
		ctx context.Context,
		roomID pulid.ID,
	) (string, error)
}

type service struct {
	config               config.LivekitConfig
	entClient            *ent.Client
	authService          auth.Service
	notificationsService notifications.Service
}

func NewService(
	config config.LivekitConfig,
	entClient *ent.Client,
	authService auth.Service,
	notificationsService notifications.Service,
) Service {
	return &service{
		config:               config,
		entClient:            entClient,
		authService:          authService,
		notificationsService: notificationsService,
	}
}

func (s *service) StartCall(
	ctx context.Context,
	roomID pulid.ID,
	callType model.CallType,
) (bool, error) {
	currentUser, err := s.authService.AuthUser(ctx)
	if err != nil {
		return false, err
	}

	userContact, err := s.getInversedUserContact(ctx, roomID, currentUser.ID)
	if err != nil {
		return false, err
	}

	_, err = s.notificationsService.Client().Send(ctx, &proto.NotificationRequest{
		Platform:         2,
		Tokens:           []string{userContact.Edges.User.Edges.Device.FcmToken},
		Priority:         proto.NotificationRequest_HIGH,
		ContentAvailable: true,
		Data:             s.getCallData(model.CallNotificationTypeStartCall, callType, userContact),
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *service) EndCall(
	ctx context.Context,
	roomID pulid.ID,
	callType model.CallType,
) (bool, error) {
	currentUser, err := s.authService.AuthUser(ctx)
	if err != nil {
		return false, err
	}

	userContact, err := s.getInversedUserContact(ctx, roomID, currentUser.ID)
	if err != nil {
		return false, err
	}

	_, err = s.notificationsService.Client().Send(ctx, &proto.NotificationRequest{
		Platform:         2,
		Tokens:           []string{userContact.Edges.User.Edges.Device.FcmToken},
		Priority:         proto.NotificationRequest_HIGH,
		ContentAvailable: true,
		Data:             s.getCallData(model.CallNotificationTypeEndCall, callType, userContact),
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *service) DeclineCall(
	ctx context.Context,
	roomID pulid.ID,
	callType model.CallType,
) (bool, error) {
	currentUser, err := s.authService.AuthUser(ctx)
	if err != nil {
		return false, err
	}

	userContact, err := s.getInversedUserContact(ctx, roomID, currentUser.ID)
	if err != nil {
		return false, err
	}

	_, err = s.notificationsService.Client().Send(ctx, &proto.NotificationRequest{
		Platform:         2,
		Tokens:           []string{userContact.Edges.User.Edges.Device.FcmToken},
		Priority:         proto.NotificationRequest_HIGH,
		ContentAvailable: true,
		Data:             s.getCallData(model.CallNotificationTypeDeclineCall, callType, userContact),
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *service) AnswerCall(
	ctx context.Context,
	roomID pulid.ID,
	callType model.CallType,
) (bool, error) {
	currentUser, err := s.authService.AuthUser(ctx)
	if err != nil {
		return false, err
	}

	userContact, err := s.getInversedUserContact(ctx, roomID, currentUser.ID)
	if err != nil {
		return false, err
	}

	_, err = s.notificationsService.Client().Send(ctx, &proto.NotificationRequest{
		Platform:         2,
		Tokens:           []string{userContact.Edges.User.Edges.Device.FcmToken},
		Priority:         proto.NotificationRequest_HIGH,
		ContentAvailable: true,
		Data:             s.getCallData(model.CallNotificationTypeAnswerCall, callType, userContact),
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
	notificationType model.CallNotificationType,
	callType model.CallType,
	userContact *ent.UserContact,
) *structpb.Struct {
	return &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"callNotificationType": {
				Kind: &structpb.Value_StringValue{StringValue: notificationType.String()},
			},
			"callType": {
				Kind: &structpb.Value_StringValue{StringValue: callType.String()},
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
