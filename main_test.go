package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// loadMockResponse is a helper to load and unmarshal the canonical response.json
func loadMockResponse(t *testing.T) *WeatherAPIResponse {
	t.Helper()
	responsePath := filepath.Join("samples", "response.json")
	byteValue, err := os.ReadFile(responsePath)
	if err != nil {
		t.Fatalf("Failed to read mock response file: %v", err)
	}

	var response WeatherAPIResponse
	if err := json.Unmarshal(byteValue, &response); err != nil {
		t.Fatalf("Failed to unmarshal mock response: %v", err)
	}
	return &response
}

// executeCommand captures the stdout and stderr of a cobra command.
func executeCommand(args ...string) (string, error) {
	// Redirect stdout and stderr
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w // Redirect stderr to the same pipe

	// Execute the command
	rootCmd.SetArgs(args)
	// Reset flags to default state before each execution
	rootCmd.Flags().Set("json", "false")
	err := rootCmd.Execute()

	// Restore stdout and stderr and read the captured output
	w.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String(), err
}

func TestAppOutput(t *testing.T) {
	// --- Setup Mocks ---
	mockResponse := loadMockResponse(t)

	// Keep original functions
	originalGetWeather := GetWeatherAPI
	originalLoadConfig := LoadConfig
	originalNowFunc := nowFunc
	originalIsTerminal := isTerminal
	defer func() {
		GetWeatherAPI = originalGetWeather
		LoadConfig = originalLoadConfig
		nowFunc = originalNowFunc
		isTerminal = originalIsTerminal
	}()

	// Override with mock implementations
	LoadConfig = func(configPath ConfigPath) (*Config, error) {
		return &Config{APIKey: "mock-key", Location: "Brussels"}, nil
	}
	GetWeatherAPI = func(location, apiKey string) (*WeatherAPIResponse, error) {
		return mockResponse, nil
	}
	nowFunc = func() time.Time {
		return time.Unix(mockResponse.Location.LocaltimeEpoch, 0)
	}

	t.Run("JSON Output", func(t *testing.T) {
		isTerminal = func(fd uintptr) bool { return false } // Force JSON

		actualOutput, err := executeCommand("Brussels")
		assert.NoError(t, err)

		// The app produces valid JSON, so we check for key substrings.
		assert.Contains(t, actualOutput, "\"text\":", "Output should contain the JSON key 'text'")
		assert.Contains(t, actualOutput, "1.3°", "JSON output should contain the current temperature for Brussels")
		assert.Contains(t, actualOutput, "\"tooltip\":", "Output should contain the JSON key 'tooltip'")
		assert.Contains(t, actualOutput, "-1.2°", "JSON tooltip should contain a forecast temperature for Brussels")
	})

	t.Run("Table Output", func(t *testing.T) {
		isTerminal = func(fd uintptr) bool { return true } // Force Table

		actualOutput, err := executeCommand("Brussels")
		assert.NoError(t, err)

		// Assert that key elements for the Brussels response are present
		assert.Contains(t, actualOutput, "Current:", "Table should have a 'Current' section")
		assert.Contains(t, actualOutput, "Brussels - Belgium", "Table should contain the correct location")
		assert.Contains(t, actualOutput, "1.3°", "Table should contain the current temperature")
		assert.Contains(t, actualOutput, "Hourly Forecast:", "Table should have an 'Hourly Forecast' section")
		assert.Contains(t, actualOutput, "-1.2°", "Table should contain a forecast temperature")
	})
}

func TestExecutionError(t *testing.T) {
	// Keep original function
	originalLoadConfig := LoadConfig
	defer func() {
		LoadConfig = originalLoadConfig
	}()

	// Mock LoadConfig to return a predictable error
	LoadConfig = func(configPath ConfigPath) (*Config, error) {
		return nil, errors.New("mock config load error")
	}

	// Since Cobra prints the error to stderr and returns it, we check both.
	output, err := executeCommand("some-location")

	assert.Error(t, err, "Expected an error to be returned from Execute()")
	assert.Contains(t, output, "mock config load error", "The error message should be printed to the console")
}
