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

func errDirectoryCreation(err error, directory string) {
	logger.
		Fatal().
		Err(err).
		Msgf("directory %s creation failed", directory)
}

func errFindingHome(err error) {
	logger.
		Fatal().
		Err(err).
		Msg("cannot determine home directory")
}
