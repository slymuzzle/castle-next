package auth

import (
	"context"
	"errors"
	"journeyhub/ent"
	"journeyhub/ent/user"
	"journeyhub/graph/model"
	"journeyhub/internal/config"
	"journeyhub/internal/db"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
}

var (
	ErrUserExist          = errors.New("user already exists")
	ErrPasswordConfirm    = errors.New("passwords do not match")
	ErrPasswordHash       = errors.New("failed to hash password")
	ErrCreateUser         = errors.New("failed to create user")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrGenerateToken      = errors.New("failed to generate token")
)

type service struct {
	config    config.AuthConfig
	dbService db.Service
}

func NewService(config config.AuthConfig, dbService db.Service) Service {
	return &service{
		config:    config,
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
		return existingUser, ErrUserExist
	}

	if password != passwordConfirmation {
		return existingUser, ErrPasswordConfirm
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return existingUser, ErrPasswordHash
	}

	user, err := entClient.User.
		Create().
		SetFirstName(firstName).
		SetLastName(lastName).
		SetEmail(email).
		SetNickname(nickname).
		SetHashedPassword(hashedPassword).
		Save(ctx)
	if err != nil {
		return user, ErrCreateUser
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
		return &model.LoginUser{}, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(
		existingUser.HashedPassword,
		[]byte(password),
	)
	if err != nil {
		return &model.LoginUser{}, ErrInvalidCredentials
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": existingUser.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
	})

	token, err := claims.SignedString([]byte(s.config.Secret))
	if err != nil {
		return &model.LoginUser{}, ErrGenerateToken
	}

	return &model.LoginUser{User: existingUser, Token: token}, nil
}
