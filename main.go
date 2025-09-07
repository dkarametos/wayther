package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

// isTerminal can be overridden in tests for deterministic behavior.
var isTerminal = isatty.IsTerminal

// nowFunc is a variable that holds the function to get the current time.
// It can be overridden in tests for deterministic behavior.
var nowFunc = time.Now

var rootCmd = &cobra.Command{
	Use:   "wayther [Location]",
	Short: "A simple weather API client",
	Long: `wayther is a CLI tool for retrieving current weather and forecasts.

You can provide location as argument.
Multiple options can be applied simultaneously.

Configuration:
  The application uses a configuration file to store your WeatherAPI key and default location.
  If no configuration file is found, you will be prompted to create one interactively.
  The 'logger' key in the config (boolean, defaults to false) enables syslog output if true.`,
	RunE: runApp,
}

func init() {
	rootCmd.Flags().StringP("config", "c", "", "Provide a custom config")
	rootCmd.Flags().StringP("output", "o", "table", "Output format (json, table)")
}

func runApp(cmd *cobra.Command, args []string) error {

	configPath, err := NewConfigPath()	
	configPath.Custom, _ = cmd.Flags().GetString("config")	

	config, err := LoadConfig(configPath)
	if err != nil {
		return err 
	}
	config.ParseCommand(cmd, args)

	weather, err := GetWeatherAPI(config.Location, config.APIKey)
	if err != nil {
		if config.Output == "json" {
			fmt.Printf("{\"text\":\" N/A 🌡 \",\"tooltip\":\" error fetching weather: %s \"}", err)
			os.Exit(0)
		}
		return err 
	}

	// Format output based on flags or TTY
	if config.Output == "json" {
		fmt.Println(formatJSON(weather))
	} else {
		fmt.Println(formatTable(weather))
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
