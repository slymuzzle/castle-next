package auth

import (
	"context"
	"errors"
	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/ent/user"
	"journeyhub/graph/model"
	"journeyhub/internal/auth/jwtauth"
	"journeyhub/internal/config"
	"journeyhub/internal/db"
	"time"

	"github.com/k0kubun/pp"
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
		email string,
		nickname string,
		password string,
		passwordConfirmation string,
	) (*ent.User, error)
	Login(
		ctx context.Context,
		nicknameOrEmail string,
		password string,
	) (*model.LoginUser, error)
	Auth(ctx context.Context) (*ent.User, error)
	JWTAuthClient() *jwtauth.JWTAuth
}

type service struct {
	config    config.AuthConfig
	jwtAuth   *jwtauth.JWTAuth
	dbService db.Service
}

func NewService(config config.AuthConfig, dbService db.Service) Service {
	jwtAuth := jwtauth.New("HS256", []byte(config.Secret), nil)
	return &service{
		config:    config,
		jwtAuth:   jwtAuth,
		dbService: dbService,
	}
}

func (s *service) Register(
	ctx context.Context,
	firstName string,
	lastName string,
	email string,
	nickname string,
	password string,
	passwordConfirmation string,
) (*ent.User, error) {
	entClient := s.dbService.Client()

	existingUser, _ := entClient.User.
		Query().
		Where(
			user.Or(
				user.Email(email),
				user.Nickname(nickname),
			),
		).Only(ctx)
	if existingUser != nil {
		return nil, ErrUserExist
	}

	if password != passwordConfirmation {
		return nil, ErrPasswordConfirm
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, ErrPasswordHash
	}

	user, err := entClient.User.
		Create().
		SetFirstName(firstName).
		SetLastName(lastName).
		SetEmail(email).
		SetNickname(nickname).
		SetPassword(hashedPassword).
		Save(ctx)
	if err != nil {
		return nil, ErrCreateUser
	}

	return user, nil
}

func (s *service) Login(
	ctx context.Context,
	nicknameOrEmail string,
	password string,
) (*model.LoginUser, error) {
	entClient := s.dbService.Client()

	existingUser, err := entClient.User.
		Query().
		Where(
			user.Or(
				user.Email(nicknameOrEmail),
				user.Nickname(nicknameOrEmail),
			),
		).Only(ctx)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = CompareHashAndPassword(existingUser.Password, password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	claims := map[string]interface{}{
		"sub": string(existingUser.ID),
	}
	jwtauth.SetIssuedNow(claims)
	jwtauth.SetExpiryIn(claims, time.Hour*24)

	_, token, err := s.jwtAuth.Encode(claims)
	if err != nil {
		return nil, err
	}

	return &model.LoginUser{User: existingUser, Token: token}, nil
}

func (s *service) Auth(ctx context.Context) (*ent.User, error) {
	pp.Print("Hello")

	jwtToken, _, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	subject := jwtToken.Subject()

	entClient := s.dbService.Client()

	user, err := entClient.User.
		Query().
		Where(user.ID(pulid.ID(subject))).
		Only(ctx)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *service) JWTAuthClient() *jwtauth.JWTAuth {
	return s.jwtAuth
}
