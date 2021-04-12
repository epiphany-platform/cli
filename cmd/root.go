package cmd

import (
	"fmt"
	"path"

	"github.com/epiphany-platform/cli/pkg/environment"

	"github.com/epiphany-platform/cli/internal/janitor"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/internal/util"
	"github.com/epiphany-platform/cli/pkg/configuration"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgDir             string
	logLevel           string
	config             *configuration.Config
	currentEnvironment *environment.Environment
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "e",
	Long: `E wrapper allows to interact with epiphany`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("root PersistentPreRun")
		var usedConfigDir string
		if cfgDir != "" {
			logger.Trace().Msg("configDir parameter not empty")
			usedConfigDir = cfgDir
		} else {
			logger.Trace().Msg("configDir parameter empty")
			usedConfigDir = path.Join(util.GetHomeDirectory(), util.DefaultConfigurationDirectory)
		}

		err := janitor.InitializeStructure(usedConfigDir)
		if err != nil {
			logger.Fatal().Err(err).Msg("initialization failed")
		}
		logger.Trace().Msg("will configuration.GetConfig()")
		config, err = configuration.GetConfig()
		if err != nil {
			logger.Fatal().Err(err).Msg("get config failed")
		}
		logger.Trace().Msg("will environment.Get(config.CurrentEnvironment)")
		currentEnvironment, err = environment.Get(config.CurrentEnvironment)
		if err != nil {
			logger.Fatal().Err(err).Msg("get current environment failed")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Debug().Err(err).Msg("root execute failed")
	}
}

func init() {
	logger.Initialize()

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgDir, "configDir", "", fmt.Sprintf("config directory (default is %s)", util.DefaultConfigurationDirectory))
	rootCmd.PersistentFlags().StringVar(&logLevel, "logLevel", "", fmt.Sprintf("log level (default is warn, values: [trace, debug, info, error, fatal])"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	switch logLevel {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}

	logger.Debug().Msg("read config variables")
	viper.AutomaticEnv() // read in environment variables that match
}
