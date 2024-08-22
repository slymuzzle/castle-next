package db

import (
	"fmt"

	"journeyhub/ent"
	"journeyhub/internal/platform/config"

	"entgo.io/ent/dialect"
)

type DatabaseConnection interface {
	Connect() (*ent.Client, error)
}

type PostgresConnection struct {
	config config.DatabaseConfig
}

func (p *PostgresConnection) Connect() (*ent.Client, error) {
	return ent.Open(
		dialect.Postgres,
		fmt.Sprintf(
			"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			p.config.Host,
			p.config.Port,
			p.config.Database,
			p.config.User,
			p.config.Password,
		),
	)
}
