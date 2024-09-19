package nats

import (
	"context"
	"time"

	"journeyhub/internal/platform/config"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type serviceLogging struct {
	logger log.Logger
	Service
}

func NewServiceLogging(logger log.Logger, s Service) Service {
	return &serviceLogging{logger, s}
}

func (s *serviceLogging) Connect(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Connect",
			"url", s.Client().Conn.ConnectedUrlRedacted(),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Connect(ctx)
}

func (s *serviceLogging) Client() (client *nats.EncodedConn) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Client",
			"url", client.Conn.ConnectedUrlRedacted(),
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Client()
}

func (s *serviceLogging) JetStream() (jetStream *jetstream.JetStream) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "JetStream",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.JetStream()
}

func (s *serviceLogging) Config() (config config.NatsConfig) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Config",
			"url", s.Client().Conn.ConnectedUrlRedacted(),
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Config()
}

func (s *serviceLogging) Close() (err error) {
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
