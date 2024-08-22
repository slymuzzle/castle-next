package contacts

import (
	"context"
	"time"

	"journeyhub/ent"

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

func (s *loggingService) GenerateUserPinCode(
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

func (s *loggingService) AddUserContact(
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
