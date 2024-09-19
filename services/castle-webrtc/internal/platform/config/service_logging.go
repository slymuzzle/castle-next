package config

import (
	"time"

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

func (s *serviceLogging) LoadConfig() (cfg Config, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "LoadConfig",
			"path", configFilePath,
			"config", cfg,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.LoadConfig()
}
