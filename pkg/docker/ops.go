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

func debugJson(json []byte, format string, v ...interface{}) {
	logger.
		Debug().
		RawJSON("json", json).
		Msgf(format, v...)
}

func debug(format string, v ...interface{}) {
	logger.
		Debug().
		Msgf(format, v...)
}

func warnRemovingContainer(err error) {
	logger.
		Warn().
		Err(err).
		Msg("cannot remove container after it finished it's job")
}
