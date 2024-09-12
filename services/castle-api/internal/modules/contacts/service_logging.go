package contacts

import (
	"context"
	"time"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type ServiceLogging struct {
	logger log.Logger
	Service
}

func NewServiceLogging(logger log.Logger, s Service) Service {
	return &ServiceLogging{logger, s}
}

func (s *ServiceLogging) GenerateUserPinCode(
	ctx context.Context,
) (pincode *string, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "GenerateUserPinCode",
			"pincode", pincode,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GenerateUserPinCode(ctx)
}

func (s *ServiceLogging) AddUserContact(
	ctx context.Context,
	pincode string,
) (usr *ent.UserContact, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "AddUserContact",
			"pincode", pincode,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AddUserContact(ctx, pincode)
}

func (s *ServiceLogging) DeleteUserContact(
	ctx context.Context,
	ID pulid.ID,
) (usr *ent.UserContact, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "DeleteUserContact",
			"ID", ID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.DeleteUserContact(ctx, ID)
}
