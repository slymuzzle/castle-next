package config

import (
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

func (s *loggingService) LoadConfig() (cfg Config, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "LoadConfig",
			"path", configFilePath,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.LoadConfig()
}
