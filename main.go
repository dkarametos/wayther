package main

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"strings"
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
	// rootCmd.Flags().Bool("json", false, "output to JSON [compatible with waybar]")  
	rootCmd.Flags().StringP("output", "o", "table", "Output format (json, table)")
}


func runApp(cmd *cobra.Command, args []string) error {

	configPath, err := UserConfigPath()	
	configPath.Custom, _ = cmd.Flags().GetString("config")	

	config, err := LoadConfig(configPath)
	if err != nil {
		return err 
	}

	//config.IsOutputJSON, _ =cmd.Flags().GetBool("json")

	//put a switch here.. 
	config.Output, _ =cmd.Flags().GetString("output")	
	if !isTerminal(os.Stdout.Fd()) {
		config.Output = "json"
	}

	if len(args) > 0 {
		config.Location = strings.Join(args, " ")
	}

	// Configure syslog if enabled
	if config.Logger {
		syslogWriter, err := syslog.New(syslog.LOG_NOTICE|syslog.LOG_DAEMON, "wayther")
		if err != nil {
			log.Printf("Failed to connect to syslog: %v", err)
		} else {
			log.SetOutput(syslogWriter)
			log.SetFlags(0) // Syslog adds its own timestamp and hostname
		}
	}

	weather, err := GetWeatherAPI(config.Location, config.APIKey)
	if err != nil {
		if config.Output == "json" {
			fmt.Printf("{\"text\":\" N/A 🌡 \",\"tooltip\":\" error fetching weather: %s \"}", err)
			return nil
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
	rootCmd.Execute()
}
