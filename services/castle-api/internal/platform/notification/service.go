package notification

import (
	"journeyhub/internal/platform/config"

	"github.com/appleboy/gorush/rpc/proto"
	"github.com/nats-io/nats.go/jetstream"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service interface {
	Connect() error

	Client() proto.GorushClient

	Config() config.NotificationsConfig

	Close() error
}

type service struct {
	config    config.NotificationsConfig
	conn      *grpc.ClientConn
	client    proto.GorushClient
	jetStream *jetstream.JetStream
}

func NewService(config config.NotificationsConfig) Service {
	return &service{
		config: config,
	}
}

func (s *service) Connect() error {
	conn, err := grpc.NewClient(
		s.config.Host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	s.conn = conn

	client := proto.NewGorushClient(conn)
	s.client = client

	return nil
}

func (s *service) Client() proto.GorushClient {
	return s.client
}

func (s *service) Config() config.NotificationsConfig {
	return s.config
}

func (s *service) Close() error {
	return s.conn.Close()
}
