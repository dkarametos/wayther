package main

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"strings"
	"time"
	"path/filepath"

	"github.com/mattn/go-isatty"
)

// nowFunc is a variable that holds the function to get the current time.
// It can be overridden in tests for deterministic behavior.
var nowFunc = time.Now

func ProcessArgs(args[] string) (*Config, error) {
  configPath, err := UserConfigPath()

	// Custom parsing has to be refactored...
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-h", "--help":
			displayHelp()
			os.Exit(0)
		case "-c", "--config":
			if i+1 < len(args) {
				if filepath.IsAbs(args[i+1]) {
					configPath.Custom = args[i+1]
				} else {
					configPath.Custom = "./"+args[i+1]
				}
				i++ // Skip the path argument
			} else {
				log.Fatalf("Error: %s flag requires a path", args[i-1])
			}
		}
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	
	for i, arg := range args[1:] {
		if arg == "--json" {
			config.IsOutputJSON = true
		} else if !strings.HasPrefix(arg, "-") && strings.HasPrefix(args[i-1], "-") {
			config.Location = arg
		}
	}

	if config.Location == "" {
		return nil, fmt.Errorf("no location provided. Please provide a location as an argument or set a default in config.json")
	}
	
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		config.IsOutputJSON = true
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

	return config, nil
}


func runApp(config *Config) (string, error) {

	//call the weather API
	weather, err := GetWeather(config.Location, config.APIKey)
	if err != nil {
		return "", err
	}

	// Format output based on flags or TTY
	if config.IsOutputJSON {
		return formatJSON(weather)
	}
	return formatTable(weather), nil
}


func main() {

	// Process the Args and get a config
	config, err := ProcessArgs(os.Args)
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// Run the application logic with actual dependencies
	output, err := runApp(config)
	if err != nil {
		if config.IsOutputJSON {
			output = fmt.Sprintf("{\"text\":\"🤔 ❓\",\"tooltip\":\" error fetching weather: %s \"}", err)	
		} else {
		  log.Fatalf("Application error: %v", err)
		}
	}

	fmt.Println(output)
}
