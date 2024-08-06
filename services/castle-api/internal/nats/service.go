package nats

import (
	"fmt"
	"journeyhub/internal/config"

	"github.com/nats-io/nats.go"
)

type Service interface {
	Connect() (*nats.EncodedConn, error)
}

type service struct {
	config config.NatsConfig
}

func NewService(config config.NatsConfig) Service {
	return &service{
		config: config,
	}
}

func (s *service) Connect() (*nats.EncodedConn, error) {
	natsConn, err := nats.Connect(
		fmt.Sprintf(
			"nats://%s:%d",
			s.config.Host,
			s.config.Port,
		),
	)
	if err != nil {
		return &nats.EncodedConn{}, err
	}

	natsEncodedConn, err := nats.NewEncodedConn(
		natsConn,
		nats.JSON_ENCODER,
	)
	if err != nil {
		return natsEncodedConn, err
	}

	return natsEncodedConn, err
}
