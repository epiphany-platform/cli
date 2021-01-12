package configuration

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
	logger = zerolog.New(output).With().Str("package", "configuration").Caller().Timestamp().Logger()
}

func debug(format string, v ...interface{}) {
	logger.
		Debug().
		Msgf(format, v...)
}

func errCreateEnvironment(err error) {
	logger.
		Panic().
		Err(err).
		Msg("creation of new environment failed")
}

func errIncorrectInitialization(err error) {
	logger.
		Panic().
		Err(err).
		Msg("incorrect initialization")
}

func errSave(err error) {
	logger.
		Panic().
		Err(err).
		Msg("failed to save")
}
