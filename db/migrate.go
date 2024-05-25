package db

import (
	"embed"
	"imgdd/logging"
	"os/exec"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

var logger = logging.GetLogger("main")

func getMigrateInstance(conf DBConfigDef) *migrate.Migrate {
	d, err := iofs.New(fs, "migrations")
	if err != nil {
		logger.Panic().Err(err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, conf.URI())
	if err != nil {
		logger.Panic().Err(err)
	}
	return m
}

func RunMigrationUp(conf DBConfigDef) {
	m := getMigrateInstance(conf)
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info().Msg("No migrations to run")
			return
		}
		logger.Error().Err(err).Msg("Migrate up failed")
	} else {
		logger.Info().Msg("Migrated up")
	}
}

func MigrateToVersion(conf DBConfigDef, version uint) {
	m := getMigrateInstance(conf)
	if err := m.Migrate(version); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info().Msg("No migrations to run")
			return
		}
		logger.Error().Err(err).Msg("Migration failed")
	} else {
		logger.Info().Uint("version", version).Msg("Migrated to version")
	}
}

func CreateMigration(name string) {
	// This is used in development to create a new migration file
	cmd := exec.Command("go", "run", "github.com/golang-migrate/migrate/v4/cmd/migrate", "create", "-ext", "sql", "-dir", "db/migrations", "-seq", name)
	cmd.Dir = "." // relative to project root (where go.mod is)
	out, err := cmd.Output()
	if err != nil {
		logger.Panic().Err(err).Msg("Could not create migration")
	} else {
		logger.Info().Str("Event", "Migration created").Msg(string(out))
	}
}
