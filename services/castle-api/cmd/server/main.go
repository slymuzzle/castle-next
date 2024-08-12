package main

import (
	"fmt"
	"journeyhub/graph"
	graphmiddleware "journeyhub/graph/middleware"
	"journeyhub/internal/auth"
	"journeyhub/internal/chat"
	"journeyhub/internal/config"
	"journeyhub/internal/db"
	"journeyhub/internal/nats"
	"journeyhub/internal/validation"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/websocket"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/lib/pq"
)

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowDebug())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	var configService config.Service
	configService = config.NewService()
	configService = config.NewLoggingService(
		log.With(logger, "component", "config"),
		configService,
	)

	config, err := configService.LoadConfig()
	if err != nil {
		level.Error(logger).Log("exit", err)
		os.Exit(1)
	}

	var natsService nats.Service
	natsService = nats.NewService(config.Nats)
	natsService = nats.NewLoggingService(
		log.With(logger, "component", "nats"),
		natsService,
	)

	if err = natsService.Connect(); err != nil {
		level.Error(logger).Log("exit", err)
		os.Exit(1)
	}
	defer natsService.Close()

	var dbService db.Service
	dbService = db.NewService(config.Database)
	dbService = db.NewLoggingService(
		log.With(logger, "component", "database"),
		dbService,
	)

	if err = dbService.Connect(); err != nil {
		level.Error(logger).Log("exit", err)
		os.Exit(1)
	}
	defer dbService.Close()

	var validationService validation.Service
	validationService = validation.NewService()
	validationService = validation.NewLoggingService(
		log.With(logger, "component", "validation"),
		validationService,
	)

	var authService auth.Service
	authService = auth.NewService(config.Auth, dbService)
	authService = auth.NewLoggingService(
		log.With(logger, "component", "auth"),
		authService,
	)

	var chatService chat.Service
	chatRepository := chat.NewRepository(dbService)
	chatService = chat.NewService(chatRepository, natsService)
	chatService = chat.NewLoggingService(
		log.With(logger, "component", "chat"),
		chatService,
	)

	httpLogger := log.With(logger, "component", "http")

	// Initialize chi router
	router := chi.NewRouter()

	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// GraphQL middleware stack
	router.Use(graphmiddleware.JwtMiddleware(authService))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	graphqlQueryHandler := handler.NewDefaultServer(
		graph.NewSchema(
			dbService,
			validationService,
			authService,
			chatService,
		),
	)
	graphqlQueryHandler.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		// InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
		// 	return webSocketInit(ctx, initPayload)
		// },
	})
	router.Handle("/query", graphqlQueryHandler)

	graphqlPlaygroundHandler := playground.Handler("GraphQL", "/query")
	router.Get("/", graphqlPlaygroundHandler)

	addr := fmt.Sprintf(":%d", config.Server.Port)

	level.Info(httpLogger).Log(
		"msg", "start server",
		"addr", addr,
	)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), router); err != nil {
		level.Error(logger).Log("exit", err)
		os.Exit(1)
	}
}
