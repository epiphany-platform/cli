package configuration

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	logger zerolog.Logger
)

func init() {
	logger = log.
		With().
		Str("package", "configuration").
		Logger()
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
