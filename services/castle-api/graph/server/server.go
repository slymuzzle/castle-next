package server

import (
	"context"
	"journeyhub/internal/auth/jwtauth"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
)

func webSocketInit(ctx context.Context, ja *jwtauth.JWTAuth, initPayload transport.InitPayload) (context.Context, *transport.InitPayload, error) {
	tokenString := initPayload.Authorization()

	jwtToken, err := jwtauth.VerifyToken(ja, tokenString)
	ctxNew := jwtauth.NewContext(ctx, jwtToken, err)

	return ctxNew, &initPayload, err
}

func NewDefaultServer(es graphql.ExecutableSchema, ja *jwtauth.JWTAuth) *handler.Server {
	srv := handler.New(es)

	srv.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		InitFunc: func(
			ctx context.Context,
			initPayload transport.InitPayload,
		) (context.Context, *transport.InitPayload, error) {
			return webSocketInit(ctx, ja, initPayload)
		},
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	return srv
}
