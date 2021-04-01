package cmd

import (
	"fmt"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/util"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgDir   string
	logLevel string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "e",
	Long: `E wrapper allows to interact with epiphany`,
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

	logger.Debug().Msg("initializing root config")
	if cfgDir != "" {
		config, err := configuration.SetConfigDirectory(cfgDir)
		if err != nil {
			logger.Fatal().Err(err).Msg("set config failed")
		}
		// Use config file from the flag.
		viper.SetConfigFile(config.GetConfigFilePath())
	} else {
		config, err := configuration.GetConfig()
		if err != nil {
			logger.Fatal().Err(err).Msg("get config failed")
		}
		// setup default
		viper.SetConfigFile(config.GetConfigFilePath())
	}
	logger.Debug().Msg("read config variables")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Info().Msgf("used config file: %s", viper.ConfigFileUsed())
	}
}
