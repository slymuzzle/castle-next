package db

import (
	"time"

	"journeyhub/ent"
	"journeyhub/internal/platform/config"

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

func (s *serviceLogging) Connect() (err error) {
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

func (s *serviceLogging) Config() (config config.DatabaseConfig) {
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

func (s *serviceLogging) Client() (client *ent.Client) {
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

func (s *serviceLogging) Close() (err error) {
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
