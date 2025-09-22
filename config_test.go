package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	
	"testing"

	"github.com/spf13/cobra"
)

func TestNewConfigPath(t *testing.T) {
	configPath, err := NewConfigPath()
	if err != nil {
		t.Fatalf("UserConfigPath() returned an error: %v", err)
	}

	uConfigDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("os.UserConfigDir() returned an error: %v", err)
	}

	expectedDefConf := filepath.Join(uConfigDir, "wayther", "config.json")
	if configPath.DefConf != expectedDefConf {
		t.Errorf("Expected DefConf to be '%s', got '%s'", expectedDefConf, configPath.DefConf)
	}

	if configPath.Custom != "" {
		t.Errorf("Expected Custom to be an empty string, got '%s'", configPath.Custom)
	}
}

// Helper function to create a temporary config file
func createTempConfigFile(t *testing.T, dir, filename string, config *Config) string {
	t.Helper()
	path := filepath.Join(dir, filename)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(config); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	return path
}

func TestLoadConfig_DefaultConfigExists(t *testing.T) {
	tempDir := t.TempDir()
	defaultConfig := &Config{APIKey: "default_key", Location: "DefaultCity", Logger: false}
	defaultConfigPath := createTempConfigFile(t, tempDir, "wayther/config.json", defaultConfig)

	configPath := ConfigPath{
		DefConf: defaultConfigPath,
	}

	loadedConfig, err := (&FileConfigProvider{}).LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loadedConfig.APIKey != defaultConfig.APIKey {
		t.Errorf("Expected APIKey %v, got %v", defaultConfig.APIKey, loadedConfig.APIKey)
	}
	if loadedConfig.Location != defaultConfig.Location {
		t.Errorf("Expected Location %v, got %v", defaultConfig.Location, loadedConfig.Location)
	}
	if loadedConfig.Logger != defaultConfig.Logger {
		t.Errorf("Expected Logger %v, got %v", defaultConfig.Logger, loadedConfig.Logger)
	}
}

func TestLoadConfig_CustomConfigOverridesDefault(t *testing.T) {
	tempDir := t.TempDir()
	defaultConfig := &Config{APIKey: "default_key", Location: "DefaultCity", Logger: false}
	customConfig := &Config{APIKey: "custom_key", Location: "CustomCity", Logger: true, Output: "json", ShortTmpl: "custom_json_template", ForecastTmpl: "custom_json_forecast"}

	defaultConfigPath := createTempConfigFile(t, tempDir, "wayther/config.json", defaultConfig)
	customConfigPath := createTempConfigFile(t, tempDir, "custom/config.json", customConfig)

	configPath := ConfigPath{
		DefConf: defaultConfigPath,
		Custom:  customConfigPath,
	}

	loadedConfig, err := (&FileConfigProvider{}).LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loadedConfig.APIKey != customConfig.APIKey {
		t.Errorf("Expected APIKey %s, got %s", customConfig.APIKey, loadedConfig.APIKey)
	}
	if loadedConfig.Location != customConfig.Location {
		t.Errorf("Expected Location %s, got %s", customConfig.Location, loadedConfig.Location)
	}
	if loadedConfig.Logger != customConfig.Logger {
		t.Errorf("Expected Logger %t, got %t", customConfig.Logger, loadedConfig.Logger)
	}
}



// Helper function to simulate user input
func simulateUserInput(t *testing.T, input string) *os.File {
	t.Helper()
	tempFile, err := os.CreateTemp(t.TempDir(), "stdin")
	if err != nil {
		t.Fatalf("Failed to create temp file for stdin: %v", err)
	}

	if _, err := tempFile.WriteString(input); err != nil {
		t.Fatalf("Failed to write to temp stdin file: %v", err)
	}

	if _, err := tempFile.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek temp stdin file: %v", err)
	}

	return tempFile
}


func TestLoadConfig_CreateCustomConfig(t *testing.T) {
	tempDir := t.TempDir()
	defaultConfig := &Config{APIKey: "default_key", Location: "DefaultCity", Logger: false}
	defaultConfigPath := createTempConfigFile(t, tempDir, "wayther/config.json", defaultConfig)
	customConfigPath := filepath.Join(tempDir, "custom/config.json")

	configPath := ConfigPath{
		DefConf: defaultConfigPath,
		Custom:  customConfigPath,
	}

	// Simulate user input
	input := "CustomLocation\n"
	tempStdin := simulateUserInput(t, input)
	oldStdin := os.Stdin
	os.Stdin = tempStdin
	defer func() {
		os.Stdin = oldStdin
		tempStdin.Close()
	}()

	loadedConfig, err := (&FileConfigProvider{}).LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loadedConfig.Location != "CustomLocation" {
		t.Errorf("Expected Location CustomLocation, got %s", loadedConfig.Location)
	}
}

func TestSetDefaults(t *testing.T) {
	config := &Config{}
	config.SetDefaults()

	if config.Output != "table" {
		t.Errorf("Expected Output to be 'table', got '%s'", config.Output)
	}

	// Verify JSON defaults
	if config.ShortTmpl == "" {
		t.Errorf("Expected ShortTmpl to have a default value")
	}
	if config.ForecastTmpl == "" {
		t.Errorf("Expected ForecastTmpl to have a default value")
	}

	// Verify Table defaults
	
}

func TestMergeConfigs(t *testing.T) {
	baseConfig := &Config{
		APIKey:     "base_key",
		Location:   "BaseCity",
		Logger:     false,
		Output: "table",
		ShortTmpl: "base_json_template",
		ForecastTmpl: "base_json_forecast",
	}

	customConfig := &Config{
		APIKey:     "custom_key",
		Location:   "CustomCity",
		Logger:     true,
		Output: "json",
		ShortTmpl: "custom_json_template",
		ForecastTmpl: "custom_json_forecast",
	}

	baseConfig.MergeConfigs(customConfig)

	if baseConfig.APIKey != "custom_key" {
		t.Errorf("Expected APIKey to be 'custom_key', got '%s'", baseConfig.APIKey)
	}
	if baseConfig.Location != "CustomCity" {
		t.Errorf("Expected Location to be 'CustomCity', got '%s'", baseConfig.Location)
	}
	if baseConfig.Logger != true {
		t.Errorf("Expected Logger to be true, got %v", baseConfig.Logger)
	}
	if baseConfig.Output != "json" {
		t.Errorf("Expected Output to be 'json', got '%s'", baseConfig.Output)
	}
	if baseConfig.ShortTmpl != "custom_json_template" {
		t.Errorf("Expected ShortTmpl to be 'custom_json_template', got '%s'", baseConfig.ShortTmpl)
	}
	if baseConfig.ForecastTmpl != "custom_json_forecast" {
		t.Errorf("Expected ForecastTmpl to be 'custom_json_forecast', got '%s'", baseConfig.ForecastTmpl)
	}
}

func TestParseCommand(t *testing.T) {
	config := &Config{}
	cmd := &cobra.Command{}
	cmd.Flags().StringP("output", "o", "table", "Output format (json, table)")

	// Test with args
	args := []string{"London"}
	config.ParseCommand(cmd, args, true)
	if config.Location != "London" {
		t.Errorf("Expected Location to be 'London', got '%s'", config.Location)
	}

	// Test with output flag
	cmd.Flags().Set("output", "json")
	config.ParseCommand(cmd, args, false)
	if config.Output != "json" {
		t.Errorf("Expected Output to be 'json', got '%s'", config.Output)
	}
}