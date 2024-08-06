package auth

import (
	"context"
	"journeyhub/ent"
	"journeyhub/graph/model"
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
			"nicknameOrEmail", nicknameOrEmail,
			"token", loginUser.Token,
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
