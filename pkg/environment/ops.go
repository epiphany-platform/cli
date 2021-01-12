package environment

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
	logger = zerolog.New(output).With().Str("package", "environment").Caller().Timestamp().Logger()
}

func debug(format string, v ...interface{}) {
	logger.
		Debug().
		Msgf(format, v...)
}

func warnEnvironmentConfigFileNotFound(err error, path string) {
	logger.
		Warn().
		Err(err).
		Msgf("expected file %s not found", path)
}

func warnNotEnvironmentDirectory(err error) {
	logger.
		Warn().
		Err(err).
		Msg("does not seam like environment directory")
}

func errFailedToWriteFile(err error) {
	logger.
		Panic().
		Err(err).
		Msg("failed to write file")
}

func errSaveEnvironment(err error, uuid string) {
	logger.
		Panic().
		Err(err).
		Msgf("wasn't able to save environment %s", uuid)
}
