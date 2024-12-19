package main

import (
	"imgdd/buildflag"
	"imgdd/config"
	"imgdd/db"
	"imgdd/graph"
	"imgdd/httpserver"
	"imgdd/identity"
	"imgdd/logging"
	"imgdd/test_support"
	"log"
	"os"
	"os/exec"
	"strings"

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

func getGoBinPath() string {
	cmd := exec.Command("go", "env", "GOBIN")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	goBinPath := strings.TrimSpace(string(output))

	if goBinPath == "" {
		cmd = exec.Command("go", "env", "GOPATH")
		output, err = cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
		goBinPath = strings.TrimSpace(string(output)) + "/bin"
	}
	return goBinPath
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
					db.MigrateToVersion(&dbConf, migrateVersion)
					return nil
				} else {
					db.RunMigrationUp(&dbConf)
				}
				return nil
			},
		},
		{
			Name: "populate-built-in-roles",
			Action: func(ctx *cli.Context) error {
				dbConf := db.ReadConfigFromEnv()
				db.PopulateBuiltInRoles(&dbConf)
				return nil
			},
		},
		{
			Name: "add-user-to-group",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "group-key",
					Usage:    "Key of the group",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "user-email",
					Usage:    "Email of the user",
					Required: true,
				},
			},
			Action: func(ctx *cli.Context) error {
				dbConf := db.ReadConfigFromEnv()
				db.PopulateBuiltInRoles(&dbConf)
				err := identity.AddUserToGroup(ctx.String("group-key"), ctx.String("user-email"), &dbConf)
				return err
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
			&cli.Command{
				Name: "reset-db",
				Action: func(ctx *cli.Context) error {
					dbConf := db.ReadConfigFromEnv()
					test_support.ResetDatabase(&dbConf)
					return nil
				},
			},
		)
	}
	if buildflag.IsDev {
		commands = append(commands,
			&cli.Command{
				Name: "dev-server",
				Action: func(ctx *cli.Context) error {
					goBinPath := getGoBinPath()
					cmd := exec.Command(goBinPath+"/air", "-c", ".air.toml", "serve")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.Stdin = os.Stdin
					cmd.Run()
					return nil
				},
			},
			&cli.Command{
				Name: "foo",
				Action: func(ctx *cli.Context) error {
					config.ReadFromTomlFile(ctx.Path("config"))
					return nil
				},
			},
		)
	}
	app := &cli.App{
		Name:        "imgdd",
		Description: "imgdd command line tool",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to the config file",
			},
		},
		Commands: commands,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
