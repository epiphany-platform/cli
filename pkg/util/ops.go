/*
 * Copyright © 2020 Mateusz Kyc
 */

package util

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
		Str("package", "util").
		Logger()
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
