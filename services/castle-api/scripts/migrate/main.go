package main

import (
	"context"
	"fmt"
	"os"

	"journeyhub/ent/migrate"
	"journeyhub/internal/platform/config"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	atlas "ariga.io/atlas/sql/migrate"

	"entgo.io/ent/dialect/sql/schema"

	_ "journeyhub/ent/runtime"

	_ "github.com/lib/pq"
)

const (
	dir = "migrations"
)

func main() {
	ctx := context.Background()

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

	cfg, err := configService.LoadConfig()
	if err != nil {
		logger.Log(err)
		os.Exit(1)
	}

	dsn := fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.Driver,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		"atlas",
	)

	// Create a local migration directory able to understand Atlas migration file format for replay
	if err = os.MkdirAll(dir, 0o755); err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}

	dir, err := atlas.NewLocalDir(dir)
	if err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}

	// Migrate diff options.
	opts := []schema.MigrateOption{
		schema.WithDir(dir),
		schema.WithMigrationMode(schema.ModeReplay),
		schema.WithDialect(cfg.Database.Driver),
		schema.WithFormatter(atlas.DefaultFormatter),
		schema.WithIndent("  "),
	}

	if len(os.Args) != 2 {
		level.Error(logger).Log(
			"msg", "migration name is required. Use: 'go run -mod=mod ent/migrate/main.go <name>'",
		)
		os.Exit(1)
	}

	// Generate migrations
	err = migrate.NamedDiff(ctx, dsn, os.Args[1], opts...)
	if err != nil {
		level.Error(logger).Log(
			"msg", fmt.Sprintf("failed generating migration file: %v", err),
		)
		os.Exit(1)
	}

	// Generate data migrations
	// err = migratedata.SeedInitialUsers(cfg.Database.Driver, dir)
	// if err != nil {
	// 	level.Error(logger).Log("msg", err)
	// 	os.Exit(1)
	// }
	//
	// err = migratedata.SeedInitialRooms(cfg.Database.Driver, dir)
	// if err != nil {
	// 	level.Error(logger).Log("msg", err)
	// 	os.Exit(1)
	// }
}
