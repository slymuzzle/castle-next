package main

import (
	"fmt"
	"journeyhub/graph"
	"journeyhub/graph/server"
	"journeyhub/internal/auth"
	"journeyhub/internal/auth/jwtauth"
	"journeyhub/internal/chat"
	"journeyhub/internal/config"
	"journeyhub/internal/db"
	"journeyhub/internal/media"
	"journeyhub/internal/nats"
	"journeyhub/internal/validation"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

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

	config, cErr := configService.LoadConfig()
	if cErr != nil {
		level.Error(logger).Log("exit", cErr)
		os.Exit(1)
	}

	var dbService db.Service
	dbService = db.NewService(config.Database)
	dbService = db.NewLoggingService(
		log.With(logger, "component", "database"),
		dbService,
	)

	if dbErr := dbService.Connect(); dbErr != nil {
		level.Error(logger).Log("exit", dbErr)
		os.Exit(1)
	}
	defer dbService.Close()

	var mediaService media.Service
	mediaService, mErr := media.NewService(config.S3)
	if mErr != nil {
		level.Error(logger).Log("exit", mErr)
		os.Exit(1)
	}
	mediaService = media.NewLoggingService(
		log.With(logger, "component", "media"),
		mediaService,
	)

	var natsService nats.Service
	natsService = nats.NewService(config.Nats)
	natsService = nats.NewLoggingService(
		log.With(logger, "component", "nats"),
		natsService,
	)

	if nErr := natsService.Connect(); nErr != nil {
		level.Error(logger).Log("exit", nErr)
		os.Exit(1)
	}
	defer natsService.Close()

	var validationService validation.Service
	validationService = validation.NewService()
	validationService = validation.NewLoggingService(
		log.With(logger, "component", "validation"),
		validationService,
	)

	var authService auth.Service
	authRepository := auth.NewRepository(dbService.Client())
	authService = auth.NewService(config.Auth, authRepository)
	authService = auth.NewLoggingService(
		log.With(logger, "component", "auth"),
		authService,
	)

	var chatService chat.Service
	chatRepository := chat.NewRepository(dbService.Client())
	chatService = chat.NewService(chatRepository, natsService, mediaService)
	chatService = chat.NewLoggingService(
		log.With(logger, "component", "chat"),
		chatService,
	)

	httpLogger := log.With(logger, "component", "http")

	// Initialize chi router
	router := chi.NewRouter()

	// Initialize auth client
	jwtAuth := authService.JWTAuthClient()

	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Auth middleware stack
	router.Use(jwtauth.Verifier(jwtAuth))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	graphqlQueryHandler := server.NewDefaultServer(
		graph.NewSchema(
			dbService,
			validationService,
			authService,
			chatService,
		),
		jwtAuth,
	)
	router.Handle("/query", graphqlQueryHandler)

	graphqlPlaygroundHandler := playground.AltairHandler("GraphQL", "/query")
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
