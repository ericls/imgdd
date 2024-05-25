package db

import (
	"os/exec"
)

func JetGenerate(conf DBConfigDef) {
	cmd := exec.Command(
		"go", "run", "github.com/go-jet/jet/v2/cmd/jet", "-dsn", conf.URI(), "-schema", "public", "-path", "./db/.gen",
		"-ignore-tables", "schema_migrations",
	)
	cmd.Dir = "." // relative to project root (where go.mod is)
	out, err := cmd.Output()
	if err != nil {
		logger.Error().Err(err).Msg(string(out))
	} else {
		logger.Info().Str("Event", "Created jet files").Msg(string(out))
	}
}
