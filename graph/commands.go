package graph

import (
	"os/exec"
)

func GenerateGqlCode() {
	cmd := exec.Command(
		"go", "run", "github.com/99designs/gqlgen",
	)
	cmd.Dir = "." // relative to project root (where go.mod is)
	out, err := cmd.Output()
	if err != nil {
		commandLogger.Error().Err(err).Msg(string(out))
	} else {
		commandLogger.Info().Str("Event", "Generated graphql code").Msg(string(out))
	}
}
