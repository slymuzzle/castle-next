package auth

import (
	"context"
	"journeyhub/ent"
	"journeyhub/graph/model"
	"journeyhub/internal/auth/jwtauth"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Register(
	ctx context.Context,
	firstName string,
	lastName string,
	email string,
	nickname string,
	password string,
	passwordConfirmation string,
) (user *ent.User, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Register",
			"email", email,
			"nickname", nickname,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.Service.Register(
		ctx,
		firstName,
		lastName,
		email,
		nickname,
		password,
		passwordConfirmation,
	)
}

func (s *loggingService) Login(
	ctx context.Context,
	nicknameOrEmail string,
	password string,
) (loginUser *model.LoginUser, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Login",
			"loginUser", loginUser,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.Service.Login(
		ctx,
		nicknameOrEmail,
		password,
	)
}

func (s *loggingService) Auth(ctx context.Context) (user *ent.User, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "User",
			"user", user,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.Service.Auth(ctx)
}

func (s *loggingService) JWTAuthClient() *jwtauth.JWTAuth {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "JWTAuth",
			"took", time.Since(begin),
			"err", nil,
		)
	}(time.Now())

	return s.Service.JWTAuthClient()
}
