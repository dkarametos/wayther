package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
	"github.com/spf13/cobra"
)

// ConfigPath holds paths related to application configuration files.
type ConfigPath struct {
	DefConf string
	Custom  string
}

// NewConfigPath returns a ConfigPath struct with the default configuration file path pre-filled.
func NewConfigPath() (ConfigPath, error) {
	configPath := ConfigPath{}

	uConfigDir, err := os.UserConfigDir()
	if err != nil {
		return configPath, err
	}
	configPath.DefConf = filepath.Join(uConfigDir, "wayther", "config.json")
	configPath.Custom = ""

	return configPath, nil
}

// isCustom checks if a custom configuration file path is provided.
// It returns true if a custom path is set, false otherwise.
func (cp *ConfigPath) isCustom() bool {
	return (cp.Custom != "")
}


// Config holds the application configuration
type Config struct {
	APIKey         string   `json:"apiKey,omitempty"`
	Location       string   `json:"location"`
	Logger         bool     `json:"logger"`
	CurrentFmt     string   `json:"currentFmt,omitempty"`
	CurrentFields  []any    `json:"current,omitempty"`
	ForecastFmt    string   `json:"forecastFmt,omitempty"`
	ForecastFields []any    `json:"forecast,omitempty"`
	Output         string   `json:"output,omitempty"`
}

// SetDefaults sets the default values for the configuration.
func (c *Config) SetDefaults() {
	if c.CurrentFmt == "" {
		c.CurrentFmt = ""
	}

	if len(c.CurrentFields) == 0 {
		c.CurrentFields = []any{}
	}

	if c.ForecastFmt == "" {
		c.ForecastFmt = "" 
	}

	if len(c.ForecastFields) == 0 {
		c.ForecastFields = []any{}
	}

	if c.Output == "" {
		c.Output = "table"
	}

}


// MergeConfigs merges the custom configuration into the current configuration.
func (c *Config) MergeConfigs(customConfig *Config) {
	if customConfig.APIKey != "" {
		c.APIKey = customConfig.APIKey
	}

	if customConfig.Location != "" {
		c.Location = customConfig.Location
	}

	c.Logger = customConfig.Logger
	
	if customConfig.CurrentFmt != "" {
		c.CurrentFmt = customConfig.CurrentFmt
	}

	if len(customConfig.CurrentFields) > 0 {
		c.CurrentFields = customConfig.CurrentFields
	}

	if customConfig.ForecastFmt != "" {
		c.ForecastFmt = customConfig.ForecastFmt
	}

	if len(customConfig.ForecastFields) > 0 {
		c.ForecastFields = customConfig.ForecastFields
	}

	if customConfig.Output != "" {
		c.Output = customConfig.Output
	}
}

// ParseCommand parses the command-line arguments and flags and updates the Config struct.
func (c *Config) ParseCommand (cmd *cobra.Command, args []string) {

	//put a switch here.. 
	c.Output, _ =cmd.Flags().GetString("output")	
	if !isTerminal(os.Stdout.Fd()) {
		c.Output = "json"
	}

	if len(args) > 0 {
		c.Location = strings.Join(args, " ")
	}

	// Configure syslog if enabled
	if c.Logger {
		syslogWriter, err := syslog.New(syslog.LOG_NOTICE|syslog.LOG_DAEMON, "wayther")
		if err != nil {
			log.Printf("Failed to connect to syslog: %v", err)
		} else {
			log.SetOutput(syslogWriter)
			log.SetFlags(0) // Syslog adds its own timestamp and hostname
		}
	}
}

// LoadConfig loads the application configuration from the default path or a custom path.
// It merges configurations if a custom path is provided.
// It returns a Config struct or an error if loading/creating fails.
func LoadConfig(configPath ConfigPath) (*Config, error) {

	//Load or Create the defaultConfig
	config, err := LoadOrCreateConfig(configPath.DefConf, true)
	if err != nil {
		return nil, err
	}
	config.SetDefaults()

	// If there is a custom config specified
	if configPath.isCustom() {
		customConfig, err := LoadOrCreateConfig(configPath.Custom, false)
		if err != nil {
			return nil, err
		}		
		config.MergeConfigs(customConfig)	
	}
	
	return config, nil
}


// LoadOrCreateConfig loads a configuration from a given path if it exists,
// otherwise it creates a new one.
// path: The path to the configuration file.
// isDefault: A boolean to indicate if this is the default configuration.
// It returns a Config struct or an error if loading/creating fails.
func LoadOrCreateConfig(path string, isDefault bool) (*Config, error) {

	var err error
	config := &Config{}

	if PathExists(path) {
		config, err = LoadConfigFromFile(path)
		if err != nil {
			return nil, err
		}
	} else {
		if isDefault {
			fmt.Println("Default config not found.\n We will create a default config first:")
		}
		config, err = CreateConfig(path, isDefault)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

// PathExists checks if a given file or directory path exists.
// It returns true if the path exists, false otherwise.
func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
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

// CreateConfig interactively prompts the user for configuration details
// and creates a new configuration file at the specified path.
// It returns the created Config struct or an error.
func CreateConfig(path string, isDefault bool) (*Config, error) {

	var err error

	config := &Config{}

	//start the dialog
	if isDefault {
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

		fmt.Printf("%s is required\n", prompt)
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
