package auth

import (
	"context"
	"time"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/graph/model"
	"journeyhub/internal/modules/auth/jwtauth"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type serviceLogging struct {
	logger log.Logger
	Service
}

func NewServiceLogging(logger log.Logger, s Service) Service {
	return &serviceLogging{logger, s}
}

func (s *serviceLogging) Register(
	ctx context.Context,
	input model.UserRegisterInput,
) (user *ent.User, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Register",
			"nickname", input.Nickname,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.Service.Register(ctx, input)
}

func (s *serviceLogging) Login(
	ctx context.Context,
	input model.UserLoginInput,
) (loginUser *model.LoginUser, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Login",
			"loginUser", input.Nickname,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.Service.Login(ctx, input)
}

func (s *serviceLogging) Auth(ctx context.Context) (userID pulid.ID, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "User",
			"userID", userID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.Service.Auth(ctx)
}

func (s *serviceLogging) AuthUser(ctx context.Context) (user *ent.User, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "User",
			// "userID", user.ID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.Service.AuthUser(ctx)
}

func (s *serviceLogging) JWTAuthClient() *jwtauth.JWTAuth {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "JWTAuth",
			"took", time.Since(begin),
			"err", nil,
		)
	}(time.Now())

	return s.Service.JWTAuthClient()
}
