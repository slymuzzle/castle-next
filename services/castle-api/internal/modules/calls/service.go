package calls

import (
	"context"
	"errors"
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

var ErrFcmTokenNotExists = errors.New("fcm token not exist")

type Service interface {
	StartCall(
		ctx context.Context,
		input model.CallParamsInput,
	) (bool, error)

	EndCall(
		ctx context.Context,
		input model.CallParamsInput,
	) (bool, error)

	DeclineCall(
		ctx context.Context,
		input model.CallParamsInput,
	) (bool, error)

	AnswerCall(
		ctx context.Context,
		input model.CallParamsInput,
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
	input model.CallParamsInput,
) (bool, error) {
	currentUser, err := s.authService.AuthUser(ctx)
	if err != nil {
		return false, err
	}

	userContact, err := s.getInversedUserContact(ctx, input.RoomID, currentUser.ID)
	if err != nil {
		return false, err
	}

	if userContact != nil &&
		userContact.Edges.User != nil &&
		userContact.Edges.User.Edges.Device != nil &&
		userContact.Edges.User.Edges.Device.FcmToken != "" {

		_, err = s.notificationsService.Client().Send(ctx, &proto.NotificationRequest{
			Platform:         2,
			Tokens:           []string{userContact.Edges.User.Edges.Device.FcmToken},
			Priority:         proto.NotificationRequest_HIGH,
			ContentAvailable: true,
			Data: s.getCallData(
				model.CallNotificationTypeStartCall,
				input.CallType,
				input.CallID,
				userContact,
			),
		})
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, ErrFcmTokenNotExists
}

func (s *service) EndCall(
	ctx context.Context,
	input model.CallParamsInput,
) (bool, error) {
	currentUser, err := s.authService.AuthUser(ctx)
	if err != nil {
		return false, err
	}

	userContact, err := s.getInversedUserContact(ctx, input.RoomID, currentUser.ID)
	if err != nil {
		return false, err
	}

	if userContact != nil &&
		userContact.Edges.User != nil &&
		userContact.Edges.User.Edges.Device != nil &&
		userContact.Edges.User.Edges.Device.FcmToken != "" {

		_, err = s.notificationsService.Client().Send(ctx, &proto.NotificationRequest{
			Platform:         2,
			Tokens:           []string{userContact.Edges.User.Edges.Device.FcmToken},
			Priority:         proto.NotificationRequest_HIGH,
			ContentAvailable: true,
			Data: s.getCallData(
				model.CallNotificationTypeEndCall,
				input.CallType,
				input.CallID,
				userContact,
			),
		})
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, ErrFcmTokenNotExists
}

func (s *service) DeclineCall(
	ctx context.Context,
	input model.CallParamsInput,
) (bool, error) {
	currentUser, err := s.authService.AuthUser(ctx)
	if err != nil {
		return false, err
	}

	userContact, err := s.getInversedUserContact(ctx, input.RoomID, currentUser.ID)
	if err != nil {
		return false, err
	}

	if userContact != nil &&
		userContact.Edges.User != nil &&
		userContact.Edges.User.Edges.Device != nil &&
		userContact.Edges.User.Edges.Device.FcmToken != "" {

		_, err = s.notificationsService.Client().Send(ctx, &proto.NotificationRequest{
			Platform:         2,
			Tokens:           []string{userContact.Edges.User.Edges.Device.FcmToken},
			Priority:         proto.NotificationRequest_HIGH,
			ContentAvailable: true,
			Data: s.getCallData(
				model.CallNotificationTypeDeclineCall,
				input.CallType,
				input.CallID,
				userContact,
			),
		})
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, ErrFcmTokenNotExists
}

func (s *service) AnswerCall(
	ctx context.Context,
	input model.CallParamsInput,
) (bool, error) {
	currentUser, err := s.authService.AuthUser(ctx)
	if err != nil {
		return false, err
	}

	userContact, err := s.getInversedUserContact(ctx, input.RoomID, currentUser.ID)
	if err != nil {
		return false, err
	}

	if userContact != nil &&
		userContact.Edges.User != nil &&
		userContact.Edges.User.Edges.Device != nil &&
		userContact.Edges.User.Edges.Device.FcmToken != "" {

		_, err = s.notificationsService.Client().Send(ctx, &proto.NotificationRequest{
			Platform:         2,
			Tokens:           []string{userContact.Edges.User.Edges.Device.FcmToken},
			Priority:         proto.NotificationRequest_HIGH,
			ContentAvailable: true,
			Data: s.getCallData(
				model.CallNotificationTypeAnswerCall,
				input.CallType,
				input.CallID,
				userContact,
			),
		})
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, ErrFcmTokenNotExists
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
		WithContact().
		Where(
			usercontact.RoomID(roomID),
			usercontact.ContactID(currentUserID),
		).
		Only(ctx)
}

func (s *service) getCallData(
	notificationType model.CallNotificationType,
	callType model.CallType,
	callID string,
	userContact *ent.UserContact,
) *structpb.Struct {
	return &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"callNotificationType": {
				Kind: &structpb.Value_StringValue{StringValue: notificationType.String()},
			},
			"callID": {
				Kind: &structpb.Value_StringValue{StringValue: callID},
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
				Kind: &structpb.Value_StringValue{StringValue: string(userContact.Edges.Contact.ID)},
			},
			"userFirstName": {
				Kind: &structpb.Value_StringValue{StringValue: string(userContact.Edges.Contact.FirstName)},
			},
			"userLastName": {
				Kind: &structpb.Value_StringValue{StringValue: string(userContact.Edges.Contact.LastName)},
			},
		},
	}
}
