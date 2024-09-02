package main

import (
	"fmt"
	"net/http"
	"os"

	"journeyhub/graph"
	"journeyhub/graph/server"
	"journeyhub/internal/modules/auth"
	"journeyhub/internal/modules/auth/jwtauth"
	"journeyhub/internal/modules/chat"
	"journeyhub/internal/modules/contacts"
	"journeyhub/internal/modules/media"
	"journeyhub/internal/modules/roommembers"
	"journeyhub/internal/modules/rooms"
	"journeyhub/internal/platform/config"
	"journeyhub/internal/platform/db"
	"journeyhub/internal/platform/nats"
	"journeyhub/internal/platform/validation"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "journeyhub/ent/runtime"

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

	var dbService db.Service
	dbService = db.NewService(config.Database)
	dbService = db.NewLoggingService(
		log.With(logger, "component", "db"),
		dbService,
	)

	if dbErr := dbService.Connect(); dbErr != nil {
		level.Error(logger).Log("exit", dbErr)
		os.Exit(1)
	}
	defer dbService.Close()

	// Initialize ent client
	entClient := dbService.Client().Debug()

	// Ent hooks stack
	entLogger := log.With(logger, "component", "ent")
	entClient.Use(db.LoggingHook(entLogger))

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

	var validationService validation.Service
	validationService = validation.NewService()
	validationService = validation.NewLoggingService(
		log.With(logger, "component", "validation"),
		validationService,
	)

	var authService auth.Service
	authRepository := auth.NewRepository(entClient)
	authService = auth.NewService(config.Auth, authRepository)
	authService = auth.NewLoggingService(
		log.With(logger, "component", "auth"),
		authService,
	)

	var roomMembersService roommembers.Service
	roomMembersRepository := roommembers.NewRepository(entClient)
	roomMembersService = roommembers.NewService(roomMembersRepository, authService, natsService)
	roomMembersService = roommembers.NewLoggingService(
		log.With(logger, "component", "roommembers"),
		roomMembersService,
	)

	var roomsService rooms.Service
	roomsRepository := rooms.NewRepository(entClient)
	roomsService = rooms.NewService(roomsRepository, roomMembersService, authService, natsService)
	roomsService = rooms.NewLoggingService(
		log.With(logger, "component", "rooms"),
		roomsService,
	)

	var chatService chat.Service
	chatRepository := chat.NewRepository(entClient)
	chatService = chat.NewService(chatRepository, authService, roomsService, roomMembersService, natsService, mediaService)
	chatService = chat.NewLoggingService(
		log.With(logger, "component", "chat"),
		chatService,
	)

	var contactsService contacts.Service
	contactsRepository := contacts.NewRepository(entClient)
	contactsService = contacts.NewService(contactsRepository, authService)

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
	// router.Use(middleware.Timeout(60 * time.Second))

	graphqlLogger := log.With(logger, "component", "graphql")
	graphqlQueryHandler := server.NewDefaultServer(
		graph.NewSchema(
			dbService,
			validationService,
			authService,
			roomsService,
			roomMembersService,
			chatService,
			contactsService,
		),
		graphqlLogger,
		jwtAuth,
		entClient,
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
