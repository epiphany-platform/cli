package cmd

import (
	"fmt"

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
	Use:   "e",
	Short: "[root description never actually showed]",
	Long:  `E wrapper allows to interact with epiphany`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		errRootExecute(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgDir, "configDir", "", fmt.Sprintf("config directory (default is %s)", util.DefaultConfigurationDirectory))
	rootCmd.PersistentFlags().StringVar(&logLevel, "logLevel", "", fmt.Sprintf("log level (default is warn, values: [debug, info, error, fatal])"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	switch logLevel {
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

	debug("initializing root config")
	if cfgDir != "" {
		config, err := configuration.SetConfigDirectory(cfgDir)
		if err != nil {
			errSetConfigFile(err)
		}
		// Use config file from the flag.
		viper.SetConfigFile(config.GetConfigFilePath())
	} else {
		config, err := configuration.GetConfig()
		if err != nil {
			errGetConfig(err)
		}
		// setup default
		viper.SetConfigFile(config.GetConfigFilePath())
	}
	debug("read config variables")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		infoConfigFile(viper.ConfigFileUsed())
	}
}
