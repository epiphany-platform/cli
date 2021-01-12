package repository

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

var (
	logger zerolog.Logger
)

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger = zerolog.New(output).With().Str("package", "repository").Caller().Timestamp().Logger()
}

func debug(format string, v ...interface{}) {
	logger.
		Debug().
		Msgf(format, v...)
}

func errGetRepository(err error) {
	logger.
		Panic().
		Err(err).
		Msg("get repository failed")
}
