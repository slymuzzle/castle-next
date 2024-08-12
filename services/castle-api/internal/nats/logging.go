package nats

import (
	"journeyhub/internal/config"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/nats-io/nats.go"
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
			"url", s.Client().Conn.ConnectedUrlRedacted(),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Connect()
}

func (s *loggingService) Config() (config config.NatsConfig) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Config",
			"url", s.Client().Conn.ConnectedUrlRedacted(),
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Config()
}

func (s *loggingService) Client() (client *nats.EncodedConn) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Client",
			"url", client.Conn.ConnectedUrlRedacted(),
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Client()
}

func (s *loggingService) Close() (err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Close",
			"url", s.Client().Conn.ConnectedUrlRedacted(),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Close()
}
