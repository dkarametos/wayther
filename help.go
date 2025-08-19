package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func displayHelp() {
	fmt.Println("Usage: wayther [location] [flags]")
	fmt.Println("")
	fmt.Println("A command-line weather application.")
	fmt.Println("")
	fmt.Println("Arguments:")
	fmt.Println("  [location]    Optional. The city or location to get weather for. If not provided,")
	fmt.Println("                the default location from the configuration file will be used.")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  -c, --config <path>  Specify a custom path for the configuration file.")
	fmt.Println("  --json               Output weather data in JSON format.")
	fmt.Println("  -h, --help           Display this help message and exit.")
	fmt.Println("")

	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Error getting user config directory: %v\n", err)
		return
	}
	defaultConfigPath := filepath.Join(configDir, "wayther", "config.json")
	fmt.Printf("Default configuration path: %s\n", defaultConfigPath)
	fmt.Println("")
	fmt.Println("Configuration:")
	fmt.Println("  The application uses a configuration file to store your WeatherAPI key and default location.")
	fmt.Println("  If no configuration file is found, you will be prompted to create one interactively.")
	fmt.Println("  The 'logger' key in the config (boolean, defaults to false) enables syslog output if true.")
}