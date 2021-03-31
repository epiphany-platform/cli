package logger

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

var (
	l           zerolog.Logger
	initialized bool
)

func Initialize() {
	if !initialized {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		l = zerolog.New(output).With().Caller().Timestamp().Logger()
		initialized = true
	}
}

func Panic() *zerolog.Event {
	return l.Panic()
}

func Fatal() *zerolog.Event {
	return l.Fatal()
}

func Error() *zerolog.Event {
	return l.Error()
}

func Warn() *zerolog.Event {
	return l.Warn()
}

func Info() *zerolog.Event {
	return l.Info()
}

func Debug() *zerolog.Event {
	return l.Debug()
}
