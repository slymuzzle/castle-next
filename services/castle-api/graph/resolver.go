package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"journeyhub/ent"
	"journeyhub/graph/generated"
	"journeyhub/internal/auth"
	"journeyhub/internal/chat"
	"journeyhub/internal/validation"

	"github.com/99designs/gqlgen/graphql"
)

// Resolver is the resolver root.
type Resolver struct {
	client            *ent.Client
	validationService validation.Service
	authService       auth.Service
	chatService       chat.Service
}

// NewSchema creates a graphql executable schema.
func NewSchema(
	client *ent.Client,
	validationService validation.Service,
	authService auth.Service,
	chatService chat.Service,
) graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{
			client,
			validationService,
			authService,
			chatService,
		},
	})
}
