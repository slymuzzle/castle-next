package middleware

import (
	"context"
	"errors"
	"journeyhub/ent"
	"journeyhub/internal/auth"
	"net/http"
	"strings"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var (
	userCtxKey      = &contextKey{"user"}
	ErrAccessDenied = errors.New("access denied")
)

type contextKey struct {
	name string
}

// Middleware decodes the share session cookie and packs the session into context
func JwtMiddleware(authService auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := TokenFromHTTPRequest(r)

			// Allow unauthenticated users in
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Get the user from the database
			user, err := authService.User(r.Context(), token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}

			// Put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, user)

			// And call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// Get jwt token from request
func TokenFromHTTPRequest(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	var tokenString string

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) > 1 {
		tokenString = splitToken[1]
	}

	return tokenString
}

// Finds the user from the context. REQUIRES Middleware to have run.
func JwtUserForContext(ctx context.Context) *ent.User {
	raw, _ := ctx.Value(userCtxKey).(*ent.User)
	return raw
}
