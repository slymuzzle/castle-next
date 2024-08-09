package db

import (
	"journeyhub/ent"
	"journeyhub/internal/config"
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

func (s *loggingService) Connect() (err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Connect",
			"driver", s.Service.Config().Driver,
			"database", s.Service.Config().Database,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Connect()
}

func (s *loggingService) Config() (config config.DatabaseConfig) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Config",
			"driver", config.Driver,
			"database", config.Database,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Config()
}

func (s *loggingService) Client() (client *ent.Client) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Client",
			"driver", s.Service.Config().Driver,
			"database", s.Service.Config().Database,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Client()
}

func (s *loggingService) Close() (err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Close",
			"driver", s.Service.Config().Driver,
			"database", s.Service.Config().Database,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Close()
}
