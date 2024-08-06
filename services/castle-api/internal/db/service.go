package db

import (
	"errors"
	"journeyhub/ent"
	"journeyhub/internal/config"
)

type Service interface {
	Connect() (*ent.Client, error)
}

var ErrUnsupportedDatabaseDriver = errors.New("unsupported database driver")

type service struct {
	config config.DatabaseConfig
}

func NewService(config config.DatabaseConfig) Service {
	return &service{
		config: config,
	}
}

func (s *service) Connect() (*ent.Client, error) {
	var dbConnection DatabaseConnection

	switch s.config.Driver {
	case "postgres":
		dbConnection = &PostgresConnection{config: s.config}
	default:
		return nil, ErrUnsupportedDatabaseDriver
	}

	client, err := dbConnection.Connect()
	if err != nil {
		return nil, err
	}

	return client, nil
}
