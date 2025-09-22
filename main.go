package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wayther [Location]",
	Short: "A simple weatherapi.com cli client",
	Long: `wayther is a CLI tool for retrieving current weather and forecasts from weatherapi.com.

You You can provide location as argument.
Multiple options can be applied simultaneously.

Configuration:
  The application uses a configuration file to store your WeatherAPI key and default location.
  If no configuration file is found, you will be prompted to create one interactively.
  The 'logger' key in the config (boolean, defaults to false) enables syslog output if true.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := NewConfigPath()
		if err != nil {
			return err
		}
		configPath.Custom, _ = cmd.Flags().GetString("config")

		cache, err := NewCache(configPath.GetPath())
		if err != nil {
			return err
		}

		weatherProvider := &weatherapiProvider{cache: cache}
		configProvider := &FileConfigProvider{}
		isTerminal := isatty.IsTerminal(os.Stdout.Fd())

		return runApp(cmd, args, configPath, weatherProvider, configProvider, isTerminal, time.Now)
	},
}

func init() {
	rootCmd.Flags().StringP("config",         "c", "",      "Provide a custom config")
	rootCmd.Flags().StringP("output",         "o", "table", "Output format (json, table)")
	rootCmd.Flags().IntP(   "forecast-hours", "n", 23,      "Number of forecast hours to display (1-23). 0 means no hourly forecast.")
	rootCmd.Flags().BoolP(  "no-cache",       "f", false,   "Force a refresh of the data from the API")
	rootCmd.Flags().BoolP(  "clean-cache",    "C", false,   "Clean cache entries older than 1h")
}

// runApp is the main application logic.
func runApp(cmd *cobra.Command, args []string, configPath ConfigPath, weatherProvider WeatherProvider, configProvider ConfigProvider, isTerminal bool, nowFunc func() time.Time) error {

	cleanCache, _ := cmd.Flags().GetBool("clean-cache")
	if cleanCache {
		weatherProvider.CleanCache(time.Hour) // Clean entries older than 1 hour
	}

	config, err := configProvider.LoadConfig(configPath)
	if err != nil {
		return handleExitError(config, err, isTerminal) 
	}

	config.ParseCommand(cmd, args, isTerminal)
	weather, err := NewWeather(weatherProvider, config)
	if err != nil {
		return handleExitError(config, err, isTerminal) 
	}

	// Format output
	output, err := FormatOutput(weather, config, nowFunc)
	if err != nil {
		return handleExitError(config, err, isTerminal)
	}
	fmt.Println(output)

	return nil
}

// handleExitError provides a centralized way to handle errors and exit the application.
// It considers whether the output is to a terminal or if JSON output is requested.
func handleExitError(config *Config, err error, isTerminal bool) error {

	if (config == nil && !isTerminal) || (config != nil && config.Output == "json") {
		fmt.Printf("{\"text\":\"N/A ☢\",\"tooltip\":\" error fetching weather: %s \"}", err)
		return nil
	}

	return err
}

// main is the entry point of the application.
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
