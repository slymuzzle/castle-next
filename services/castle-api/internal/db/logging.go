package db

import (
	"journeyhub/ent"
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

func (s *loggingService) Connect() (client *ent.Client, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Connect",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Connect()
}
