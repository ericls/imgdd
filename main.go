package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ericls/imgdd/buildflag"
	"github.com/ericls/imgdd/config"
	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/email"
	"github.com/ericls/imgdd/graph"
	"github.com/ericls/imgdd/httpserver"
	"github.com/ericls/imgdd/identity"
	"github.com/ericls/imgdd/logging"
	"github.com/ericls/imgdd/test_support"

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
	var migrateVersion uint = 0
	getConfig := func(ctx *cli.Context) *config.ConfigDef {
		conf, err := config.GetConfig(ctx.Path("config"))
		if err != nil {
			panic(err)
		}
		return conf
	}
	commands := []*cli.Command{
		{
			Name: "serve",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "bind",
					Value: "127.0.0.1:8000",
					Usage: "Which address to bind to when starting the server",
				},
			},
			Action: func(ctx *cli.Context) error {
				conf := getConfig(ctx)
				httpServerConf := conf.HttpServer
				bind := ctx.String("bind")
				if bind != "" {
					httpServerConf.Bind = bind
				}
				if buildflag.IsDev {
					httpServerConf.EnableGqlPlayground = true
				}
				httpServerConf.StaticFS = MoutingFS.Static
				httpServerConf.TemplatesFS = MoutingFS.Templates
				srv := httpserver.MakeServer(&httpServerConf, &conf.Db, &conf.Storage)
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
				dbConf := getConfig(ctx).Db
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
				dbConf := getConfig(ctx).Db
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
				dbConf := getConfig(ctx).Db
				db.PopulateBuiltInRoles(&dbConf)
				err := identity.AddUserToGroup(ctx.String("group-key"), ctx.String("user-email"), &dbConf)
				return err
			},
		},
		{
			Name: "gen-config",
			Action: func(ctx *cli.Context) error {
				if ctx.Path("config") == "" {
					return config.PrintEmptyConfig()
				}
				return config.GenerateEmptyConfigFile(ctx.Path("config"))
			},
		},
		{
			Name: "print-config",
			Action: func(ctx *cli.Context) error {
				conf := getConfig(ctx)
				conf.PrintConfig()
				return nil
			},
		},
		{
			Name: "send-test-email",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "to",
					Usage:    "Email address to send the test email to",
					Required: true,
				},
			},
			Action: func(ctx *cli.Context) error {
				conf := getConfig(ctx)
				emailBackend, err := email.GetEmailBackendFromConfig(&conf.Email)
				if err != nil {
					return err
				}
				err = email.SendEmail(emailBackend, "", []string{ctx.String("to")}, "IMGDD Test email", "This is a test email", "")
				return err
			},
		},
		{
			Name: "build-info",
			Action: func(ctx *cli.Context) error {
				buildflag.PrintBuildInfo()
				return nil
			},
		},
	}
	if buildflag.IsDev {
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
					dbConf := getConfig(ctx).Db
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
					dbConf := getConfig(ctx).Db
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
