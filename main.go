package main

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
)

type GetWeatherFunc func(location, apiKey string) (*WeatherAPIResponse, error)

// nowFunc is a variable that holds the function to get the current time.
// It can be overridden in tests for deterministic behavior.
var nowFunc = time.Now


func main() {
	configPath, _ := UserConfigPath()

	args := []string{}

	// Check for help flag first
	for _, arg := range os.Args[1:] {
		if arg == "-h" || arg == "--help" {
			displayHelp()
			os.Exit(0)
		}
	}

	// Custom parsing to find the config path early
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-c", "--config":
			if i+1 < len(os.Args) {
				configPath.Custom = os.Args[i+1]
				i++ // Skip the path argument
			} else {
				log.Fatalf("Error: %s flag requires a path", os.Args[i])
			}
		default:
			args = append(args, os.Args[i])
		}
	}

	//Prepend the program name to the filtered args
	args = append([]string{os.Args[0]}, args...)

	config, err := LoadConfig(configPath)
	if err != nil {
		fmt.Errorf("Could not load file: %w", err)
		os.Exit(1)
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

	// Run the application logic with actual dependencies
	output, err := runApp(args, config, GetWeather)
	if err != nil {
		log.Fatalf("Application error: %v", err)
	}

	fmt.Println(output)

	// format := "test: %s %s %d"
	// value  := []any{"worked", "NA",3}
	//fmt.Println(config.CurrentFields[1])
	//config.CurrentFields[1] = 12
	//fmt.Printf(config.CurrentFormat, config.CurrentFields...)

}

func runApp(args []string, config *Config, getWeather GetWeatherFunc) (string, error) {
	var location string
	jsonOutput := false

	// Parse command-line arguments
	for _, arg := range args[1:] {
		if arg == "--json" {
			jsonOutput = true
		} else if !strings.HasPrefix(arg, "-") {
			location = arg
		}
	}

	if !isatty.IsTerminal(os.Stdout.Fd()) {
		jsonOutput = true
	}

	// If location is not provided via command line, use the one from config
	if location == "" {
		location = config.Location
	}

	if location == "" {
		return "", fmt.Errorf("no location provided. Please provide a location as an argument or set a default in config.json")
	}

	weather, err := getWeather(location, config.APIKey)
	if err != nil {
		if jsonOutput {
			fmt.Printf("{\"text\":\"🤔 ❓\",\"tooltip\":\" error fetching weather: %s \"}", err)
			os.Exit(0)
		} else {
			return "", fmt.Errorf("error fetching weather2: %w", err)
		}
	}

	// Format output based on flags or TTY
	if jsonOutput {
		return formatJSON(weather)
	}

	return formatTable(weather), nil
}
