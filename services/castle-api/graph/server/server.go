package server

import (
	"context"
	"net/http"
	"time"

	"journeyhub/internal/modules/auth/jwtauth"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

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

func NewDefaultServer(es graphql.ExecutableSchema, logger log.Logger, ja *jwtauth.JWTAuth) *handler.Server {
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
			ctx, payload, err := webSocketInit(ctx, ja, initPayload)
			level.Debug(logger).Log("method", "InitFunc", "auth", payload.Authorization())
			return ctx, payload, err
		},
		CloseFunc: func(ctx context.Context, closeCode int) {
			level.Error(logger).Log("method", "CloseFunc", "code", closeCode)
		},
		ErrorFunc: func(ctx context.Context, err error) {
			level.Error(logger).Log("method", "ErrorFunc", "err", err)
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
