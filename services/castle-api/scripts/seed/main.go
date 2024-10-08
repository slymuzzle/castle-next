package main

import (
	"os"

	"journeyhub/internal/platform/config"
	"journeyhub/internal/platform/db"
	"journeyhub/scripts/seed/seeddata"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	_ "journeyhub/ent/runtime"

	_ "github.com/lib/pq"
)

const (
	dir = "migrations"
)

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowDebug())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	var configService config.Service
	configService = config.NewService()
	configService = config.NewServiceLogging(
		log.With(logger, "component", "config"),
		configService,
	)

	config, err := configService.LoadConfig()
	if err != nil {
		logger.Log(err)
		os.Exit(1)
	}

	var dbService db.Service
	dbService = db.NewService(config.Database)
	dbService = db.NewServiceLogging(
		log.With(logger, "component", "database"),
		dbService,
	)

	if err = dbService.Connect(); err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}
	defer dbService.Close()

	err = seeddata.SeedUsers0(dbService)
	if err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}

	err = seeddata.SeedUsers1(dbService)
	if err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}

	err = seeddata.SeedRooms(dbService)
	if err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}
}
