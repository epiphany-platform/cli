/*
 * Copyright © 2020 Mateusz Kyc
 */

package repository

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	logger zerolog.Logger
)

func init() {
	logger = log.With().
		Str("package", "repository").
		Logger()
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
