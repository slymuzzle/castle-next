package db

import (
	"errors"
	"journeyhub/ent"
	"journeyhub/internal/config"
)

type Service interface {
	Connect() error
	Close() error
	Config() config.DatabaseConfig
	Client() *ent.Client
}

var ErrUnsupportedDatabaseDriver = errors.New("unsupported database driver")

type service struct {
	config    config.DatabaseConfig
	entClient *ent.Client
}

func NewService(config config.DatabaseConfig) Service {
	return &service{
		config:    config,
		entClient: &ent.Client{},
	}
}

func (s *service) Connect() error {
	var dbConnection DatabaseConnection

	switch s.config.Driver {
	case "postgres":
		dbConnection = &PostgresConnection{config: s.config}
	default:
		return ErrUnsupportedDatabaseDriver
	}

	client, err := dbConnection.Connect()
	if err != nil {
		return err
	}

	s.entClient = client

	return err
}

func (s *service) Config() config.DatabaseConfig {
	return s.config
}

func (s *service) Client() *ent.Client {
	return s.entClient
}

func (s *service) Close() error {
	return s.entClient.Close()
}
