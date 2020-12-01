package az

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	logger zerolog.Logger
)

func init() {
	// logLevel := zerolog.InfoLevel
	// zerolog.SetGlobalLevel(logLevel)
	// logger = zerolog.New(os.Stdout).With().Str("package", "az").Logger()

	logger = log.With().
		Str("package", "az").
		Logger()
}

func debug(format string, v ...interface{}) {
	logger.
		Debug().
		Msgf(format, v...)
}

func info(msg string) {
	logger.Info().Msg(msg)
}

func warnAssignRoleToServicePrincipal(err error) {
	logger.
		Warn().
		Err(err).
		Msg("Failed to assign role to Service Principal.")
}

func errFailedToGetEnvironment(err error) {
	logger.
		Panic().
		Err(err).
		Msg("Failed to get environment.")
}

func errFailedToGetAuthrorizerFromCli(err error) {
	logger.
		Panic().
		Err(err).
		Msg("Failed to get Authrorizer from CLI.")
}

func errFailedToGetGraphAuthrorizer(err error) {
	logger.
		Panic().
		Err(err).
		Msg("Failed to get Graph Authrorizer.")
}

func errFailedToGeneratePassword(err error) {
	logger.
		Panic().
		Err(err).
		Msg("Failed to generate password.")
}

func errFailedToCreateApplication(err error) {
	logger.
		Panic().
		Err(err).
		Msg("Failed to create application.")
}

func errFailedToMarshalJSON(err error) {
	logger.
		Panic().
		Err(err).
		Msg("Failed to marshal JSON.")
}

func errFailedToGetRoleDefinitionIterator(err error) {
	logger.
		Panic().
		Err(err).
		Msg("Failed to get role definition iterator.")
}

func errFailedToIterateOverRoleDefinitions(err error) {
	logger.
		Panic().
		Err(err).
		Msg("Failed to iterate over role definition iterator.")
}
