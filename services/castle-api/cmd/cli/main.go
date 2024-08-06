package main

import (
	"fmt"
	"journeyhub/ent/migrate"
	"journeyhub/internal/config"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/urfave/cli/v2"

	"ariga.io/atlas-go-sdk/atlasexec"
	atlas "ariga.io/atlas/sql/migrate"

	"entgo.io/ent/dialect/sql/schema"

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

	var cs config.Service
	cs = config.NewService()
	cs = config.NewLoggingService(log.With(logger, "component", "config"), cs)

	cfg, err := cs.LoadConfig()
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
		cfg.Database.Database,
	)

	migrateCommands := []*cli.Command{
		{
			Name:  "diff",
			Usage: "Generate diff migrations",
			Action: func(cCtx *cli.Context) error {
				if cCtx.Args().Len() != 1 {
					level.Error(logger).Log("msg", "migration name is required")
					cli.Exit(nil, 1)
				}

				if err := os.MkdirAll(dir, 0755); err != nil {
					level.Error(logger).Log("msg", err)
					cli.Exit(nil, 1)
				}

				dir, err := atlas.NewLocalDir(dir)
				if err != nil {
					level.Error(logger).Log("msg", err)
					cli.Exit(nil, 1)
				}

				opts := []schema.MigrateOption{
					schema.WithDir(dir),
					schema.WithMigrationMode(schema.ModeReplay),
					schema.WithDialect(cfg.Database.Driver),
					schema.WithFormatter(atlas.DefaultFormatter),
				}

				err = migrate.NamedDiff(cCtx.Context, dsn, cCtx.Args().Get(0), opts...)
				if err != nil {
					level.Error(logger).Log("msg", err)
					cli.Exit(nil, 1)
				}
				level.Info(logger).Log("created_migration", cCtx.Args().Get(0))

				return nil
			},
		},
		{
			Name:  "apply",
			Usage: "Apply migrations",
			Action: func(cCtx *cli.Context) error {
				logger = log.With(logger, "command", cCtx.Command.Name)

				workdir, err := atlasexec.NewWorkingDir(
					atlasexec.WithMigrations(
						os.DirFS(dir),
					),
				)
				if err != nil {
					level.Error(logger).Log("msg", err)
					cli.Exit(nil, 1)
				}
				defer workdir.Close()

				client, err := atlasexec.NewClient(workdir.Path(), "atlas")
				if err != nil {
					level.Error(logger).Log("msg", err)
					cli.Exit(nil, 1)
				}

				res, err := client.MigrateApply(cCtx.Context, &atlasexec.MigrateApplyParams{
					URL: dsn,
				})
				if err != nil {
					level.Error(logger).Log("msg", err)
					cli.Exit(nil, 1)
				}
				level.Info(logger).Log("applied_migrations", len(res.Applied))

				return nil
			},
		},
	}

	app := &cli.App{
		Name:                 "Cli",
		Usage:                "Cli commands",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:        "migrate",
				Usage:       "Migrate database commands",
				Subcommands: migrateCommands,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
