package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

// Config holds the application configuration
type Config struct {
	APIKey         string   `json:"apiKey,omitempty"`
	Location       string   `json:"location"`
	Logger         bool     `json:"logger"`
	CurrentFmt     string   `json:"currentFmt,omitempty"`
	CurrentFields  []any    `json:"current,omitempty"`
	ForecastFmt    string   `json:"forecastFmt,omitempty"`
	ForecastFields []any    `json:"forecast,omitempty"`
	IsOutputJSON   bool			`json:"outputJSON,omitempty"` // to be deprecated
  Output         string   `json:"output"` 
}

// ConfigPath holds paths related to application configuration files.
type ConfigPath struct {
	DefConf string
	Custom  string
}

// UserConfigPath returns a ConfigPath struct with the default configuration file path pre-filled.
func UserConfigPath() (ConfigPath, error) {
	configPath := ConfigPath{}

	uConfigDir, err := os.UserConfigDir()
	if err != nil {
		return configPath, err
	}
	configPath.DefConf = filepath.Join(uConfigDir, "wayther", "config.json")
	configPath.Custom = ""

	return configPath, nil
}

// LoadConfig loads the application configuration from the default path or a custom path.
// It merges configurations if a custom path is provided.
// It returns a Config struct or an error if loading/creating fails.
var LoadConfig = func(configPath ConfigPath) (*Config, error) {
	var err error

	//Load or Create the defaultConfig
	config := &Config{}
	if PathExists(configPath.DefConf) {
		config, err = LoadConfigFromFile(configPath.DefConf)
		if err != nil {
			return nil, err
		}
	} else {
		if configPath.Custom == "" {
			config, err = CreateConfig(configPath)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("exit with message like: first create the default config")
		}
	}

	// If there is a custom config specified
	if configPath.Custom != "" {
		customConfig := &Config{}
		if PathExists(configPath.Custom) {
			customConfig, err = LoadConfigFromFile(configPath.Custom)
			if err != nil {
				return nil, err
			}
		} else {
			customConfig, err = CreateConfig(configPath)
			if err != nil {
				return nil, err
			}
		}

		//merge the configs
		if customConfig.APIKey != "" {
			config.APIKey = customConfig.APIKey
		}
		config.Location = customConfig.Location
		config.Logger = customConfig.Logger

		if len(customConfig.CurrentFields) > 0 {
			config.CurrentFields = customConfig.CurrentFields
		}

		if len(customConfig.ForecastFields) > 0 {
			config.ForecastFields = customConfig.ForecastFields
		}


		if config.Output == "" {
			if customConfig.Output != "" {
				config.Output = customConfig.Output
			} else {
				config.Output = "table"
			}
		}
	}
	
	return config, nil
}

// LoadConfigFromFile loads configuration from a specified file path.
// It returns a Config struct or an error if the file cannot be read or decoded.
func LoadConfigFromFile(path string) (*Config, error) {
	config := &Config{}
	confFile, err := os.Open(path)
	if err == nil {
		defer confFile.Close()
		err = json.NewDecoder(confFile).Decode(config)
	}
	return config, err
}

// PathExists checks if a given file or directory path exists.
// It returns true if the path exists, false otherwise.
func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}


// GetUserInput prompts the user for input with a given prompt.
// It can handle secret input (like passwords) and default values.
// prompt: The message to display to the user.
// defaultValue: The value to return if the user enters an empty string.
// secret: If true, the input will be read from a password terminal.
// canBeEmpty: If true, the user can provide an empty value.
// It returns the user's input as a string, or an error if one occurred.
func GetUserInput(prompt string, defaultValue string, secret bool, canBeEmpty bool) (string, error) {

	var userInput string
	var input     string
	var err       error

	for {
		fmt.Printf("Enter %s:", prompt)

		if secret {
			byteString, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return string(byteString), err
			}
			userInput = string(byteString)
		} else {
			reader := bufio.NewReader(os.Stdin)
			userInput, err = reader.ReadString('\n')
			if err != nil {
				return userInput, err
			}
		}
	  
		input = strings.TrimSpace(userInput)
		fmt.Println()

		if input == "" && canBeEmpty == true {
			input = defaultValue
		}

		if input != "" || canBeEmpty == true {
			break
		}

	fmt.Printf("%s is required\n", prompt )
	}

	return input, nil
}

// WriteConfig writes the given Config struct to a file at the specified path.
// It creates the directory if it doesn't exist.
// config: The Config struct to write to the file.
// path: The path to the file to write to.
// It returns an error if one occurred.
func WriteConfig(config *Config, path string) error {
	
	//split out the folder and mkdir -p
	configDir, _ := filepath.Split(path)
		if configDir != "" {
		err := os.MkdirAll(configDir, 0750)
		if err != nil {
			return err
		}
	}

	//write the file to disk
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return err
	}

	return nil
}


// CreateConfig interactively prompts the user for configuration details
// and creates a new configuration file at the specified path.
// It returns the created Config struct or an error.
func CreateConfig(configPath ConfigPath) (*Config, error) {

	var err error

	config := &Config{}
	path := ""

	if configPath.Custom == "" {
		path = configPath.DefConf
	} else {
		path = configPath.Custom
	}


	//start the dialog
	if configPath.Custom == "" {
		config.APIKey, err = GetUserInput("weatheapi key", "", true, false)
		if err != nil {
			return config, err
		}
	}

	config.Location, err = GetUserInput("Location", "auto:ip", false, true)
	if err != nil {
		return config, err
	}

	config.Logger = false


	//Write to file
  err  = WriteConfig(config, path)
	if err != nil {
		return config, err
	}

	fmt.Printf("Created configuration file: %s\n\n", path)
	return config, nil
}

