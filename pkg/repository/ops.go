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

func errGetRepository(err error) {
	logger.
		Fatal().
		Err(err).
		Msg("get repository failed")
}

func errInitRepository(err error) {
	logger.
		Fatal().
		Err(err).
		Msg("initialization of repository failed")
}
