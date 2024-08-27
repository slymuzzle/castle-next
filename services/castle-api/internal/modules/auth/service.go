package auth

import (
	"context"
	"errors"
	"time"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/graph/model"
	"journeyhub/internal/modules/auth/jwtauth"
	"journeyhub/internal/platform/config"

	pass "journeyhub/internal/modules/auth/password"
)

var (
	ErrUserExist          = errors.New("user already exists")
	ErrPasswordConfirm    = errors.New("passwords do not match")
	ErrPasswordHash       = errors.New("failed to hash password")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrCreateUser         = errors.New("failed to create user")
	ErrUserNotFound       = errors.New("user not found")
	ErrGenerateToken      = errors.New("failed to generate token")
)

type Service interface {
	Register(
		ctx context.Context,
		firstName string,
		lastName string,
		nickname string,
		password string,
		passwordConfirmation string,
	) (*ent.User, error)
	Login(
		ctx context.Context,
		nickname string,
		password string,
	) (*model.LoginUser, error)
	Auth(ctx context.Context) (pulid.ID, error)
	AuthUser(ctx context.Context) (*ent.User, error)
	JWTAuthClient() *jwtauth.JWTAuth
}

type service struct {
	config         config.AuthConfig
	jwtAuth        *jwtauth.JWTAuth
	authRepository Repository
}

func NewService(config config.AuthConfig, authRepository Repository) Service {
	jwtAuth := jwtauth.New("HS256", []byte(config.Secret), nil)
	return &service{
		config:         config,
		jwtAuth:        jwtAuth,
		authRepository: authRepository,
	}
}

func (s *service) Register(
	ctx context.Context,
	firstName string,
	lastName string,
	nickname string,
	password string,
	passwordConfirmation string,
) (*ent.User, error) {
	existingUser, _ := s.authRepository.
		FindUserByNickname(ctx, nickname)
	if existingUser != nil {
		return nil, ErrUserExist
	}

	if password != passwordConfirmation {
		return nil, ErrPasswordConfirm
	}

	hashedPassword, err := pass.Hash(password)
	if err != nil {
		return nil, ErrPasswordHash
	}

	user, err := s.authRepository.CreateUser(
		ctx,
		firstName,
		lastName,
		nickname,
		hashedPassword,
	)
	if err != nil {
		return nil, ErrCreateUser
	}

	return user, nil
}

func (s *service) Login(
	ctx context.Context,
	nickname string,
	password string,
) (*model.LoginUser, error) {
	existingUser, err := s.authRepository.
		FindUserByNickname(ctx, nickname)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = pass.Compare(existingUser.Password, password)
	if err != nil {
		return nil, ErrInvalidCredentials
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

	user, err := s.authRepository.FindUserByID(
		ctx,
		pulid.ID(subject),
	)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *service) JWTAuthClient() *jwtauth.JWTAuth {
	return s.jwtAuth
}
