package notifications

import (
	"time"

	"journeyhub/internal/platform/config"

	"github.com/appleboy/gorush/rpc/proto"
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
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Connect()
}

func (s *serviceLogging) Client() (client proto.GorushClient) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Client",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Client()
}

func (s *serviceLogging) Config() (config config.NotificationsConfig) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Config",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Config()
}

func (s *serviceLogging) Close() (err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Close",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Close()
}
