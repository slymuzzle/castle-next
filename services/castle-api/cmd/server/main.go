package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"journeyhub/graph"
	"journeyhub/graph/server"
	"journeyhub/internal/modules/auth"
	"journeyhub/internal/modules/auth/jwtauth"
	"journeyhub/internal/modules/calls"
	"journeyhub/internal/modules/chat"
	"journeyhub/internal/modules/contacts"
	"journeyhub/internal/modules/media"
	"journeyhub/internal/modules/notifications"
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

	// Initialize config service
	var configService config.Service
	configService = config.NewService()
	configService = config.NewServiceLogging(
		log.With(logger, "component", "config"),
		configService,
	)

	config, cErr := configService.LoadConfig()
	if cErr != nil {
		level.Error(logger).Log("exit", cErr)
		os.Exit(1)
	}

	// Initialize nats service
	var natsService nats.Service
	natsService = nats.NewService(config.Nats)
	natsService = nats.NewServiceLogging(
		log.With(logger, "component", "nats"),
		natsService,
	)

	if nErr := natsService.Connect(context.TODO()); nErr != nil {
		level.Error(logger).Log("exit", nErr)
		os.Exit(1)
	}
	defer natsService.Close()

	// Initialize db service
	var dbService db.Service
	dbService = db.NewService(config.Database)
	dbService = db.NewServiceLogging(
		log.With(logger, "component", "db"),
		dbService,
	)

	if dbErr := dbService.Connect(); dbErr != nil {
		level.Error(logger).Log("exit", dbErr)
		os.Exit(1)
	}
	defer dbService.Close()

	// Initialize ent client
	entClient := dbService.Client()

	// Ent hooks stack
	entLogger := log.With(logger, "component", "ent")
	entClient.Use(db.LoggingHook(entLogger))

	// Initialize media service
	var mediaService media.Service
	mediaService, mErr := media.NewService(config.S3)
	if mErr != nil {
		level.Error(logger).Log("exit", mErr)
		os.Exit(1)
	}
	mediaService = media.NewServiceLogging(
		log.With(logger, "component", "media"),
		mediaService,
	)

	// Initialize validation service
	var validationService validation.Service
	validationService = validation.NewService()
	validationService = validation.NewServiceLogging(
		log.With(logger, "component", "validation"),
		validationService,
	)

	// Initialize notifications service
	var notificationsService notifications.Service
	notificationsService = notifications.NewService(config.Notifications)
	notificationsService = notifications.NewServiceLogging(
		log.With(logger, "component", "notifications"),
		notificationsService,
	)

	if nErr := notificationsService.Connect(); nErr != nil {
		level.Error(logger).Log("exit", nErr)
		os.Exit(1)
	}
	defer notificationsService.Close()

	// Initialize auth service
	var authService auth.Service
	authService = auth.NewService(config.Auth, entClient)
	authService = auth.NewServiceLogging(
		log.With(logger, "component", "auth"),
		authService,
	)

	// Initialize room members service
	var roomMembersSubscriptions roommembers.Subscriptions
	roomMembersSubscriptions = roommembers.NewSubscriptions(entClient, authService, natsService)
	roomMembersSubscriptions = roommembers.NewSubscriptionsLogging(
		log.With(logger, "component", "roommembers-subscriptions"),
		roomMembersSubscriptions,
	)
	var roomMembersService roommembers.Service
	roomMembersService = roommembers.NewService(entClient, roomMembersSubscriptions, authService, notificationsService)
	roomMembersService = roommembers.NewServiceLogging(
		log.With(logger, "component", "roommembers"),
		roomMembersService,
	)

	// Initialize rooms service
	var roomsService rooms.Service
	roomsService = rooms.NewService(entClient, roomMembersService, authService)
	roomsService = rooms.NewServiceLogging(
		log.With(logger, "component", "rooms"),
		roomsService,
	)

	// Initialize contacts service
	var contactsService contacts.Service
	contactsService = contacts.NewService(entClient, authService)
	contactsService = contacts.NewServiceLogging(
		log.With(logger, "component", "contacts"),
		contactsService,
	)

	// Initialize call service
	var callsService calls.Service
	callsService = calls.NewService(config.Livekit, entClient, authService, notificationsService)
	callsService = calls.NewServiceLogging(
		log.With(logger, "component", "calls"),
		callsService,
	)

	// Initialize chat service
	var chatSubscriptions chat.Subscriptions
	chatSubscriptions = chat.NewSubscriptions(entClient, natsService)
	chatSubscriptions = chat.NewSubscriptionsLogging(
		log.With(logger, "component", "chat-subscriptions"),
		chatSubscriptions,
	)
	var chatService chat.Service
	chatService = chat.NewService(entClient, chatSubscriptions, authService, roomsService, roomMembersService, mediaService)
	chatService = chat.NewServiceLogging(
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
	// router.Use(middleware.Timeout(60 * time.Second))

	graphqlLogger := log.With(logger, "component", "graphql")
	graphqlQueryHandler := server.NewDefaultServer(
		graph.NewSchema(
			dbService,
			validationService,
			authService,
			roomsService,
			roomMembersService,
			contactsService,
			callsService,
			chatService,
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
