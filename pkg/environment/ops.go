package environment

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	logger zerolog.Logger
)

func init() {
	logger = log.With().
		Str("package", "environment").
		Logger()
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
