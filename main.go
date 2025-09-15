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
		weatherProvider := &weatherapiProvider{}
		configProvider  := &FileConfigProvider{}
		return runApp(cmd, args, weatherProvider, configProvider, isatty.IsTerminal, time.Now)
	},
}

func init() {
	rootCmd.Flags().StringP("config", "c", "", "Provide a custom config")
	rootCmd.Flags().StringP("output", "o", "table", "Output format (json, table)")
	rootCmd.Flags().IntP("forecast-hours", "n", 23, "Number of forecast hours to display (1-23). 0 means no hourly forecast.")
}

func runApp(cmd *cobra.Command, args []string, weatherProvider WeatherProvider, configProvider ConfigProvider, isTerminal func(uintptr) bool, nowFunc func() time.Time) error {

	configPath, err := NewConfigPath()
	configPath.Custom, _ = cmd.Flags().GetString("config")

	config, err := configProvider.LoadConfig(configPath)
	if err != nil {
		exitOnJSON(config, err)
		return err
	}
	config.ParseCommand(cmd, args, isTerminal)

	weather, err := NewWeather(weatherProvider, config)
	if err != nil {
		exitOnJSON(config, err)
		return err
	}


	// Format output based on flags or TTY
	if config.OutputType == "json" {
		fmt.Println(formatJSON(weather, config, nowFunc))
	} else {
		fmt.Println(formatTable(weather, config, nowFunc))
	}

	return nil
}

func exitOnJSON(config *Config, err error ) {
	
	if config == nil {
		return
	}

	if config.OutputType == "json" {
		fmt.Printf("{\"text\":\"N/A ☢\",\"tooltip\":\" error fetching weather: %s \"}", err)
		os.Exit(0)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
