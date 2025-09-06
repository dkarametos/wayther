package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
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

	loadedConfig, err := LoadConfig(configPath)
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
	if !reflect.DeepEqual(loadedConfig.CurrentFields, defaultConfig.CurrentFields) {
		t.Errorf("Expected CurrentFields %v, got %v", defaultConfig.CurrentFields, loadedConfig.CurrentFields)
	}
	if !reflect.DeepEqual(loadedConfig.ForecastFields, defaultConfig.ForecastFields) {
		t.Errorf("Expected ForecastFields %v, got %v", defaultConfig.ForecastFields, loadedConfig.ForecastFields)
	}
}

func TestLoadConfig_CustomConfigOverridesDefault(t *testing.T) {
	tempDir := t.TempDir()
	defaultConfig := &Config{APIKey: "default_key", Location: "DefaultCity", Logger: false}
	customConfig := &Config{APIKey: "custom_key", Location: "CustomCity", Logger: true}

	defaultConfigPath := createTempConfigFile(t, tempDir, "wayther/config.json", defaultConfig)
	customConfigPath := createTempConfigFile(t, tempDir, "custom/config.json", customConfig)

	configPath := ConfigPath{
		DefConf: defaultConfigPath,
		Custom:  customConfigPath,
	}

	loadedConfig, err := LoadConfig(configPath)
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

	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loadedConfig.Location != "CustomLocation" {
		t.Errorf("Expected Location CustomLocation, got %s", loadedConfig.Location)
	}
}