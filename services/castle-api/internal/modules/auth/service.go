package auth

import (
	"context"
	"errors"
	"time"

	"journeyhub/ent"
	"journeyhub/ent/device"
	"journeyhub/ent/schema/pulid"
	"journeyhub/ent/user"
	"journeyhub/graph/model"
	"journeyhub/internal/modules/auth/jwtauth"
	"journeyhub/internal/platform/config"

	pass "journeyhub/internal/modules/auth/password"
)

var (
	ErrUserExist                = errors.New("user already exists")
	ErrPasswordConfirm          = errors.New("passwords do not match")
	ErrPasswordHash             = errors.New("failed to hash password")
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrCreateUser               = errors.New("failed to create user")
	ErrUserNotFound             = errors.New("user not found")
	ErrGenerateToken            = errors.New("failed to generate token")
	ErrCreateOrUpdateUserDevice = errors.New("failed to create or update user device")
)

type Service interface {
	Register(
		ctx context.Context,
		input model.UserRegisterInput,
	) (*ent.User, error)

	Login(
		ctx context.Context,
		input model.UserLoginInput,
	) (*model.LoginUser, error)

	Auth(ctx context.Context) (pulid.ID, error)

	AuthUser(ctx context.Context) (*ent.User, error)

	JWTAuthClient() *jwtauth.JWTAuth
}

type service struct {
	config    config.AuthConfig
	entClient *ent.Client
	jwtAuth   *jwtauth.JWTAuth
}

func NewService(config config.AuthConfig, entClient *ent.Client) Service {
	jwtAuth := jwtauth.New("HS256", []byte(config.Secret), nil)
	return &service{
		config:    config,
		entClient: entClient,
		jwtAuth:   jwtAuth,
	}
}

func (s *service) Register(ctx context.Context, input model.UserRegisterInput) (*ent.User, error) {
	existingUser, _ := s.entClient.User.
		Query().
		Where(user.Nickname(input.Nickname)).
		Only(ctx)
	if existingUser != nil {
		return nil, ErrUserExist
	}

	if input.Password != input.PasswordConfirmation {
		return nil, ErrPasswordConfirm
	}

	hashedPassword, err := pass.Hash(input.Password)
	if err != nil {
		return nil, ErrPasswordHash
	}

	user, err := s.entClient.User.
		Create().
		SetFirstName(input.FirstName).
		SetLastName(input.LastName).
		SetNickname(input.Nickname).
		SetPassword(hashedPassword).
		Save(ctx)
	if err != nil {
		return nil, ErrCreateUser
	}

	return user, nil
}

func (s *service) Login(ctx context.Context, input model.UserLoginInput) (*model.LoginUser, error) {
	existingUser, err := s.entClient.User.
		Query().
		Where(user.Nickname(input.Nickname)).
		Only(ctx)
	if err != nil {
		return nil, errors.Join(ErrInvalidCredentials, err)
	}

	err = pass.Compare(existingUser.Password, input.Password)
	if err != nil {
		return nil, errors.Join(ErrInvalidCredentials, err)
	}

	claims := map[string]interface{}{
		"sub": string(existingUser.ID),
	}
	jwtauth.SetIssuedNow(claims)
	jwtauth.SetExpiryIn(claims, time.Hour*720)

	_, token, err := s.jwtAuth.Encode(claims)
	if err != nil {
		return nil, err
	}

	_, err = s.entClient.Device.
		Create().
		SetUserID(existingUser.ID).
		SetDeviceID(input.DeviceID).
		SetFcmToken(input.FcmToken).
		OnConflictColumns(device.FieldDeviceID).
		UpdateNewValues().
		ID(ctx)
	if err != nil {
		return nil, errors.Join(ErrCreateOrUpdateUserDevice, err)
	}

	return &model.LoginUser{User: existingUser, Token: token}, nil
}

func (s *service) Auth(ctx context.Context) (pulid.ID, error) {
	jwtToken, _, err := jwtauth.FromContext(ctx)
	if err != nil {
		return pulid.ID(""), err
	}

	subject := jwtToken.Subject()

	userID := pulid.ID(subject)

	return userID, nil
}

func (s *service) AuthUser(ctx context.Context) (*ent.User, error) {
	jwtToken, _, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	subject := jwtToken.Subject()

	user, err := s.entClient.User.
		Query().
		Where(user.ID(pulid.ID(subject))).
		WithDevice().
		Only(ctx)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *service) JWTAuthClient() *jwtauth.JWTAuth {
	return s.jwtAuth
}
