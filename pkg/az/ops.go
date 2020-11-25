package az

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	logger zerolog.Logger
)

func init() {
	// logLevel := zerolog.InfoLevel
	// zerolog.SetGlobalLevel(logLevel)
	// logger = zerolog.New(os.Stdout).With().Str("package", "az").Logger()

	logger = log.With().
		Str("package", "az").
		Logger()
}

func debug(format string, v ...interface{}) {
	logger.
		Debug().
		Msgf(format, v...)
}

func info(msg string) {
	logger.Info().Msg(msg)
}
