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
	IsOutputJSON   bool			`json:"outputJSON,omitempty"`
}

type ConfigPath struct {
	DefConf string
	Custom  string
}

// userConfigPath returns the ConfigPath struct with default config prefilled
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
	}
	
	return config, nil
}

func LoadConfigFromFile(path string) (*Config, error) {
	config := &Config{}
	confFile, err := os.Open(path)
	if err == nil {
		defer confFile.Close()
		err = json.NewDecoder(confFile).Decode(config)
	}
	return config, err
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// createConfig(path)
func CreateConfig(configPath ConfigPath) (*Config, error) {

	reader := bufio.NewReader(os.Stdin)
	config := &Config{}
	path := ""

	if configPath.Custom == "" {
		path = configPath.DefConf
	} else {
		path = configPath.Custom
	}

	//split out the folder and mkdir -p
	configDir, _ := filepath.Split(path)
	err := os.MkdirAll(configDir, 0750)
	if err != nil {
		return config, err
	}

	//start the dialog

	if configPath.Custom == "" {
		for {
			fmt.Print("Enter WeatherAPI Key: ")
			byteAPIKey, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return config, err
			}
			config.APIKey = strings.TrimSpace(string(byteAPIKey))
			fmt.Println()
			if config.APIKey != "" {
				break
			}
			fmt.Println("API Key cannot be empty for default configuration.")
		}
	}

	for {
		fmt.Print("Enter location: ")
		location, err := reader.ReadString('\n')
		if err != nil {
			return config, err
		}
		config.Location = strings.TrimSpace(location)
		if config.Location != "" {
			break
		}
		fmt.Println("Location cannot be empty.")
	}

	config.Logger = false

	file, err := os.Create(path)
	if err != nil {
		return config, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return config, err
	}

	fmt.Printf("Created configuration file: %s\n", path)
	return config, nil
}

