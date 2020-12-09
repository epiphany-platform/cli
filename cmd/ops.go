package cmd

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	logger zerolog.Logger
)

func init() {
	logger = log.With().
		Str("package", "cmd").
		Logger()
}

func debug(format string, v ...interface{}) {
	logger.
		Debug().
		Msgf(format, v...)
}

func errSetConfigFile(err error) {
	logger.
		Fatal().
		Err(err).
		Msg("set config failed")
}

func errGetConfig(err error) {
	logger.
		Fatal().
		Err(err).
		Msg("get config failed")
}

func errRootExecute(err error) {
	logger.
		Fatal().
		Err(err).
		Msg("root execute failed")
}

func errGetComponentByName(err error) {
	logger.
		Fatal().
		Err(err).
		Msg("getting component by name failed")
}

func errGetComponentWithLatestVersion(err error) {
	logger.
		Fatal().
		Err(err).
		Msg("getting component with latest version failed")
}

func errTooFewArguments(err error) {
	logger.
		Fatal().
		Err(err).
		Msg("too few arguments")
}

func errGetEnvironments(err error) {
	logger.
		Fatal().
		Err(err).
		Msg("environments get failed")
}

func errIncorrectNumberOfArguments(err error) {
	logger.
		Fatal().
		Err(err).
		Msg("incorrect number of arguments")
}

func errInstallComponent(err error) {
	logger.
		Fatal().
		Err(err).
		Msg("install component in environment failed")
}

func errNilEnvironment() {
	logger.
		Fatal().
		Msg("no environment used")
}

func errGetEnvironmentDetails(err error) {
	logger.
		Fatal().
		Err(err).
		Msg("get environments details failed")
}

func errPrompt(err error) {
	logger.
		Fatal().
		Err(err).
		Msg("prompt failed")
}

func errCreateEnvironment(err error) {
	logger.
		Fatal().
		Err(err).
		Msg("create new environment failed")
}

func errRunCommand(err error) {
	logger.
		Fatal().
		Err(err).
		Msg("run command failed")
}

func errSetEnvironment(err error) {
	logger.
		Fatal().
		Err(err).
		Msg("setting used environment failed")
}

func errGeneratePassword(err error) {
	logger.
		Panic().
		Err(err).
		Msg("failed to generate password")
}

func infoConfigFile(filePath string) {
	logger.
		Info().
		Msgf("used config file: %s", filePath)
}

func infoRunFinished(component string, command string) {
	logger.
		Info().
		Msgf("running %s %s finished", component, command)
}

func infoChosenEnvironment(uuid string) {
	logger.
		Info().
		Msgf("Chosen environment UUID is %s", uuid)
}
