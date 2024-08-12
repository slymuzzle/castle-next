package nats

import (
	"fmt"
	"journeyhub/internal/config"

	"github.com/nats-io/nats.go"
)

type Service interface {
	Connect() error
	Config() config.NatsConfig
	Client() *nats.EncodedConn
	Close() error
}

type service struct {
	config config.NatsConfig
	conn   *nats.EncodedConn
}

func NewService(config config.NatsConfig) Service {
	return &service{
		config: config,
	}
}

func (s *service) Connect() error {
	natsConn, err := nats.Connect(
		fmt.Sprintf(
			"nats://%s:%d",
			s.config.Host,
			s.config.Port,
		),
	)
	if err != nil {
		return err
	}

	natsEncodedConn, err := nats.NewEncodedConn(
		natsConn,
		nats.JSON_ENCODER,
	)
	if err != nil {
		return err
	}

	s.conn = natsEncodedConn

	return err
}

func (s *service) Config() config.NatsConfig {
	return s.config
}

func (s *service) Client() *nats.EncodedConn {
	return s.conn
}

func (s *service) Close() error {
	return s.conn.Drain()
}
