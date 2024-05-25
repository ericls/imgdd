package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func GetLogger(name string) zerolog.Logger {
	sublogger := log.With().
		Str("logger", name).
		Logger()
	return sublogger
}
