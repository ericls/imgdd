package graph

import (
	"imgdd/logging"
	"os/exec"
)

var logger = logging.GetLogger("graph_commands")

func GenerateGqlCode() {
	cmd := exec.Command(
		"go", "run", "github.com/99designs/gqlgen",
	)
	cmd.Dir = "." // relative to project root (where go.mod is)
	out, err := cmd.Output()
	if err != nil {
		logger.Error().Err(err).Msg(string(out))
	} else {
		logger.Info().Str("Event", "Generated graphql code").Msg(string(out))
	}
}
