package docker

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
	logger = zerolog.New(output).With().Str("package", "docker").Caller().Timestamp().Logger()
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
