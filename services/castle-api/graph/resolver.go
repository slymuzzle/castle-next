package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"journeyhub/graph/generated"
	"journeyhub/internal/auth"
	"journeyhub/internal/chat"
	"journeyhub/internal/db"
	"journeyhub/internal/validation"

	"github.com/99designs/gqlgen/graphql"
)

// Resolver is the resolver root.
type Resolver struct {
	dbService         db.Service
	validationService validation.Service
	authService       auth.Service
	chatService       chat.Service
}

// NewSchema creates a graphql executable schema.
func NewSchema(
	dbService db.Service,
	validationService validation.Service,
	authService auth.Service,
	chatService chat.Service,
) graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{
			dbService,
			validationService,
			authService,
			chatService,
		},
	})
}
