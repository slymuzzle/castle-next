package calls

import (
	"context"
	"time"

	"journeyhub/ent/schema/pulid"

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

func (s *serviceLogging) GetCallJoinToken(
	ctx context.Context,
	roomID pulid.ID,
) (token string, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "GetCallJoinToken",
			"roomID", roomID,
			"token", token,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetCallJoinToken(ctx, roomID)
}
