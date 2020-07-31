/*
 * Copyright © 2020 Mateusz Kyc
 */

package docker

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	logger zerolog.Logger
)

func init() {
	logger = log.With().
		Str("package", "docker").
		Logger()
}

func warnRemovingContainer(err error) {
	logger.
		Warn().
		Err(err).
		Msg("cannot remove container after it finished it's job")
}
