package util

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
	logger = zerolog.New(output).With().Str("package", "util").Caller().Timestamp().Logger()
}

func debug(format string, v ...interface{}) {
	logger.
		Debug().
		Msgf(format, v...)
}

func errDirectoryCreation(err error, directory string) {
	logger.
		Panic().
		Err(err).
		Msgf("directory %s creation failed", directory)
}

func errFindingHome(err error) {
	logger.
		Panic().
		Err(err).
		Msg("cannot determine home directory")
}
