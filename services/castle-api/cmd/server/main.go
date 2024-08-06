package main

import (
	"fmt"
	"journeyhub/graph"
	"journeyhub/internal/auth"
	"journeyhub/internal/chat"
	"journeyhub/internal/config"
	"journeyhub/internal/db"
	"journeyhub/internal/nats"
	"journeyhub/internal/validation"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}

	var natsService nats.Service
	natsService = nats.NewService(config.Nats)
	natsService = nats.NewLoggingService(
		log.With(logger, "component", "nats"),
		natsService,
	)

	natsConn, err := natsService.Connect()
	if err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}
	defer natsConn.Drain()

	var dbService db.Service
	dbService = db.NewService(config.Database)
	dbService = db.NewLoggingService(
		log.With(logger, "component", "database"),
		dbService,
	)

	entClient, err := dbService.Connect()
	if err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}
	defer entClient.Close()

	var validationService validation.Service
	validationService = validation.NewService()
	validationService = validation.NewLoggingService(
		log.With(logger, "component", "validation"),
		validationService,
	)

	var authService auth.Service
	authService = auth.NewService(entClient, config.Auth)
	authService = auth.NewLoggingService(
		log.With(logger, "component", "auth"),
		authService,
	)

	var chatService chat.Service
	chatService = chat.NewService(entClient, natsConn)
	chatService = chat.NewLoggingService(
		log.With(logger, "component", "chat"),
		chatService,
	)

	server := echo.New()
	httpLogger := log.With(logger, "component", "http")

	server.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogProtocol:      true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogStatus:        true,
		LogLatency:       true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			httpLogger.Log(
				"protocol", v.Protocol,
				"remote_ip", v.RemoteIP,
				"host", v.Host,
				"method", v.Method,
				"uri", v.URI,
				"status", v.Status,
				"latency", v.Latency,
				"bytes_in", v.ContentLength,
				"bytes_out", v.ResponseSize,
			)
			return nil
		},
	}))
	server.Use(middleware.Recover())

	graphqlQueryHandler := handler.NewDefaultServer(
		graph.NewSchema(
			entClient,
			validationService,
			authService,
			chatService,
		),
	)
	graphqlQueryHandler.AddTransport(&transport.Websocket{})

	server.Any("/query", func(c echo.Context) error {
		graphqlQueryHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	graphqlPlaygroundHandler := playground.Handler("GraphQL", "/query")

	server.GET("/", func(c echo.Context) error {
		graphqlPlaygroundHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	if err := server.Start(fmt.Sprintf(":%d", config.Server.Port)); err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}
}
