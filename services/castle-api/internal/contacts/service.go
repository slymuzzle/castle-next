package contacts

import (
	"context"
	"journeyhub/internal/db"
)

type Service interface {
	GenereatePinCode(context.Context) error
}

type service struct {
	dbService db.Service
}

func NewService(dbService db.Service) Service {
	return &service{
		dbService: dbService,
	}
}

func (s *service) GenereatePinCode(ctx context.Context) error {
	return nil
}
