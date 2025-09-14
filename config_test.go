package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
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
	defaultConfig := &Config{APIKey: "default_key", Location: "DefaultCity", Logger: false, Outputs: Outputs{JSON: OutputConfig{CurrentFields: []any{}, ForecastFields: []any{}}, Table: OutputConfig{CurrentFields: []any{}, ForecastFields: []any{}}}}
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
	if !reflect.DeepEqual(loadedConfig.Outputs.JSON.CurrentFields, defaultConfig.Outputs.JSON.CurrentFields) {
		t.Errorf("Expected JSON CurrentFields %v, got %v", defaultConfig.Outputs.JSON.CurrentFields, loadedConfig.Outputs.JSON.CurrentFields)
	}
	if !reflect.DeepEqual(loadedConfig.Outputs.JSON.ForecastFields, defaultConfig.Outputs.JSON.ForecastFields) {
		t.Errorf("Expected JSON ForecastFields %v, got %v", defaultConfig.Outputs.JSON.ForecastFields, loadedConfig.Outputs.JSON.ForecastFields)
	}
	if !reflect.DeepEqual(loadedConfig.Outputs.Table.CurrentFields, defaultConfig.Outputs.Table.CurrentFields) {
		t.Errorf("Expected Table CurrentFields %v, got %v", defaultConfig.Outputs.Table.CurrentFields, loadedConfig.Outputs.Table.CurrentFields)
	}
	if !reflect.DeepEqual(loadedConfig.Outputs.Table.ForecastFields, defaultConfig.Outputs.Table.ForecastFields) {
		t.Errorf("Expected Table ForecastFields %v, got %v", defaultConfig.Outputs.Table.ForecastFields, loadedConfig.Outputs.Table.ForecastFields)
	}
}

func TestLoadConfig_CustomConfigOverridesDefault(t *testing.T) {
	tempDir := t.TempDir()
	defaultConfig := &Config{APIKey: "default_key", Location: "DefaultCity", Logger: false, Outputs: Outputs{JSON: OutputConfig{CurrentFields: []any{}, ForecastFields: []any{}}, Table: OutputConfig{CurrentFields: []any{}, ForecastFields: []any{}}}}
	customConfig := &Config{APIKey: "custom_key", Location: "CustomCity", Logger: true, OutputType: "json", Outputs: Outputs{JSON: OutputConfig{CurrentFmt: "custom_json_current", ForecastFmt: "custom_json_forecast"}, Table: OutputConfig{CurrentFmt: "custom_table_current", ForecastFmt: "custom_table_forecast"}}}

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
	defaultConfig := &Config{APIKey: "default_key", Location: "DefaultCity", Logger: false, Outputs: Outputs{JSON: OutputConfig{CurrentFields: []any{}, ForecastFields: []any{}}, Table: OutputConfig{CurrentFields: []any{}, ForecastFields: []any{}}}}
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

	if config.OutputType != "table" {
		t.Errorf("Expected OutputType to be 'table', got '%s'", config.OutputType)
	}

	// Verify JSON defaults
	if config.Outputs.JSON.CurrentFmt == "" {
		t.Errorf("Expected JSON.CurrentFmt to have a default value")
	}
	if len(config.Outputs.JSON.CurrentFields) != 0 {
		t.Errorf("Expected JSON.CurrentFields to be empty, got %v", config.Outputs.JSON.CurrentFields)
	}
	if config.Outputs.JSON.ForecastFmt == "" {
		t.Errorf("Expected JSON.ForecastFmt to have a default value")
	}
	if len(config.Outputs.JSON.ForecastFields) != 0 {
		t.Errorf("Expected JSON.ForecastFields to be empty, got %v", config.Outputs.JSON.ForecastFields)
	}

	// Verify Table defaults
	if config.Outputs.Table.CurrentFmt == "" {
		t.Errorf("Expected Table.CurrentFmt to have a default value")
	}
	if len(config.Outputs.Table.CurrentFields) != 0 {
		t.Errorf("Expected Table.CurrentFields to be empty, got %v", config.Outputs.Table.CurrentFields)
	}
	if config.Outputs.Table.ForecastFmt == "" {
		t.Errorf("Expected Table.ForecastFmt to have a default value")
	}
	if len(config.Outputs.Table.ForecastFields) != 0 {
		t.Errorf("Expected Table.ForecastFields to be empty, got %v", config.Outputs.Table.ForecastFields)
	}
}

func TestMergeConfigs(t *testing.T) {
	baseConfig := &Config{
		APIKey:     "base_key",
		Location:   "BaseCity",
		Logger:     false,
		OutputType: "table",
		Outputs: Outputs{
			JSON:  OutputConfig{CurrentFmt: "base_json_current"},
			Table: OutputConfig{CurrentFmt: "base_table_current"},
		},
	}

	customConfig := &Config{
		APIKey:     "custom_key",
		Location:   "CustomCity",
		Logger:     true,
		OutputType: "json",
		Outputs: Outputs{
			JSON:  OutputConfig{CurrentFmt: "custom_json_current"},
			Table: OutputConfig{CurrentFmt: "custom_table_current"},
		},
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
	if baseConfig.OutputType != "json" {
		t.Errorf("Expected OutputType to be 'json', got '%s'", baseConfig.OutputType)
	}
	if baseConfig.Outputs.JSON.CurrentFmt != "custom_json_current" {
		t.Errorf("Expected JSON.CurrentFmt to be 'custom_json_current', got '%s'", baseConfig.Outputs.JSON.CurrentFmt)
	}
	if baseConfig.Outputs.Table.CurrentFmt != "custom_table_current" {
		t.Errorf("Expected Table.CurrentFmt to be 'custom_table_current', got '%s'", baseConfig.Outputs.Table.CurrentFmt)
	}
}

func TestParseCommand(t *testing.T) {
	config := &Config{}
	cmd := &cobra.Command{}
	cmd.Flags().StringP("output", "o", "table", "Output format (json, table)")

	// Test with args
	args := []string{"London"}
	config.ParseCommand(cmd, args, func(uintptr) bool { return true })
	if config.Location != "London" {
		t.Errorf("Expected Location to be 'London', got '%s'", config.Location)
	}

	// Test with output flag
	cmd.Flags().Set("output", "json")
	config.ParseCommand(cmd, args, func(uintptr) bool { return false })
	if config.OutputType != "json" {
		t.Errorf("Expected OutputType to be 'json', got '%s'", config.OutputType)
	}
}