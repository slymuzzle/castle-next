package nats

import (
	"context"
	"fmt"
	"log"

	"journeyhub/internal/platform/config"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Service interface {
	Connect(ctx context.Context) error
	Client() *nats.EncodedConn
	JetStream() *jetstream.JetStream
	Config() config.NatsConfig
	Close() error
}

type service struct {
	config    config.NatsConfig
	conn      *nats.EncodedConn
	jetStream *jetstream.JetStream
}

func NewService(config config.NatsConfig) Service {
	return &service{
		config: config,
	}
}

func (s *service) Connect(ctx context.Context) error {
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

	jetStream, err := jetstream.New(natsConn)
	if err != nil {
		log.Fatal(err)
	}
	s.jetStream = &jetStream

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

func (s *service) Client() *nats.EncodedConn {
	return s.conn
}

func (s *service) JetStream() *jetstream.JetStream {
	return s.jetStream
}

func (s *service) Config() config.NatsConfig {
	return s.config
}

func (s *service) Close() error {
	return s.conn.Drain()
}
