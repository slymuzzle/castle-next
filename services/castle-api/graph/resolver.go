package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"journeyhub/graph/generated"
	"journeyhub/internal/modules/auth"
	"journeyhub/internal/modules/chat"
	"journeyhub/internal/modules/contacts"
	"journeyhub/internal/modules/roommembers"
	"journeyhub/internal/modules/rooms"
	"journeyhub/internal/platform/db"
	"journeyhub/internal/platform/validation"

	"github.com/99designs/gqlgen/graphql"
)

// Resolver is the resolver root.
type Resolver struct {
	dbService         db.Service
	validationService validation.Service
	authService       auth.Service
	roomsService      rooms.Service
	roomMemberService roommembers.Service
	chatService       chat.Service
	contactsService   contacts.Service
}

// NewSchema creates a graphql executable schema.
func NewSchema(
	dbService db.Service,
	validationService validation.Service,
	authService auth.Service,
	roomsService rooms.Service,
	roomMemberService roommembers.Service,
	chatService chat.Service,
	contactsService contacts.Service,
) graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{
			dbService,
			validationService,
			authService,
			roomsService,
			roomMemberService,
			chatService,
			contactsService,
		},
	})
}
