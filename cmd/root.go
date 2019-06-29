package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var cfgFile string
var logLevel string
var debugShorthand bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "megantory",
	Short: "megantory searches your AWS accounts really fast",
	Long:  `TODO`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run:  searchCmd.Run,  // by default we run the "search" command
	Args: searchCmd.Args, // inherit searchCmds args rules, too
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLogrus)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file (default is $HOME/.megantory.yaml)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "warn", "Log-level (trace,debug,info,warn,error,fatal,panic)")
	rootCmd.PersistentFlags().BoolVarP(&debugShorthand, "debug", "d", false, "Shorthand for --loglevel debug. Overrides '--loglevel'")

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".megantory" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".megantory")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// initLogrus initializes the logging library
func initLogrus() {

	logLevels := map[string]log.Level{
		"trace": log.TraceLevel,
		"debug": log.DebugLevel,
		"info":  log.InfoLevel,
		"warn":  log.WarnLevel,
		"error": log.ErrorLevel,
		"fatal": log.FatalLevel,
		"panic": log.PanicLevel,
	}

	if debugShorthand {
		logLevel = "debug"
	}
	if logLevel == "trace" {
		log.SetReportCaller(true)
	}
	log.SetLevel(logLevels[logLevel])
	log.Debugf("Log level set to: %v", logLevel)
	if cfgFile != "" {
		log.Debugf("Config file set to: %v", cfgFile)
	}
}
