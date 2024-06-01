package main

import (
	"imgdd/buildflag"
	"imgdd/db"
	"imgdd/graph"
	"imgdd/httpserver"
	"imgdd/logging"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var logger = logging.GetLogger("main")

func init() {
	err := godotenv.Overload(".env")
	if err != nil {
		logger.Warn().Err(err).Msg("Could not load .env file")
	}
}

func main() {
	var bind string
	var migrateVersion uint = 0
	commands := []*cli.Command{
		{
			Name: "serve",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "bind",
					Value:       "127.0.0.1:8000",
					Usage:       "Which address to bind to when starting the server",
					Destination: &bind,
				},
			},
			Action: func(ctx *cli.Context) error {
				httpServerConf := httpserver.Config
				httpServerConf.Bind = bind
				httpServerConf.StaticFS = MoutingFS.Static
				httpServerConf.TemplatesFS = MoutingFS.Templates
				srv := httpserver.MakeServer(&httpServerConf)
				logger.Info().Str("bind", srv.Addr).Msg("Starting server")
				return srv.ListenAndServe()
			},
		},
		{
			Name: "migrate",
			Flags: []cli.Flag{
				&cli.UintFlag{
					Name:        "version",
					Value:       0,
					Usage:       "Which version to migrate to",
					Destination: &migrateVersion,
				},
			},
			Action: func(ctx *cli.Context) error {
				dbConf := db.ReadConfigFromEnv()
				if migrateVersion > 0 {
					db.MigrateToVersion(dbConf, migrateVersion)
					return nil
				} else {
					db.RunMigrationUp(dbConf)
				}
				return nil
			},
		},

		{
			Name: "populate-built-in-roles",
			Action: func(ctx *cli.Context) error {
				dbConf := db.ReadConfigFromEnv()
				db.PopulateBuiltInRoles(dbConf)
				return nil
			},
		},
	}
	if buildflag.IsDebug {
		commands = append(commands,
			&cli.Command{
				Name: "make-migration",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Usage:    "Name of the migration",
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					db.CreateMigration(ctx.String("name"))
					return nil
				},
			},
			&cli.Command{
				Name: "jet",
				Action: func(ctx *cli.Context) error {
					dbConf := db.ReadConfigFromEnv()
					db.JetGenerate(dbConf)
					return nil
				},
			},
			&cli.Command{
				Name: "gql",
				Action: func(ctx *cli.Context) error {
					graph.GenerateGqlCode()
					return nil
				},
			},
		)
	}
	app := &cli.App{
		Name:        "imgdd",
		Description: "imgdd command line tool",
		Commands:    commands,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
