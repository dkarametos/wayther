package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mattn/go-isatty"
)

// Mock implementations for GetWeather and LoadConfig for testing runApp
var mockWeatherResponse *WeatherAPIResponse
var mockConfig *Config

func mockGetWeather(location, apiKey string) (*WeatherAPIResponse, error) {
	return mockWeatherResponse, nil
}



// isTerminalMock is a variable that holds the function to check if a file descriptor is a terminal.
// It can be overridden in tests for deterministic behavior.
var isTerminalMock = isatty.IsTerminal

func TestDisplayHelp(t *testing.T) {
	// Keep old stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	displayHelp()

	// Restore old stdout
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	actualOutput := buf.String()

	// Construct expected output
	var expectedOutputBuilder strings.Builder
	expectedOutputBuilder.WriteString("Usage: wayther [location] [flags]\n")
	expectedOutputBuilder.WriteString("\n")
	expectedOutputBuilder.WriteString("A command-line weather application.\n")
	expectedOutputBuilder.WriteString("\n")
	expectedOutputBuilder.WriteString("Arguments:\n")
	expectedOutputBuilder.WriteString("  [location]    Optional. The city or location to get weather for. If not provided,\n")
	expectedOutputBuilder.WriteString("                the default location from the configuration file will be used.\n")
	expectedOutputBuilder.WriteString("\n")
	expectedOutputBuilder.WriteString("Flags:\n")
	expectedOutputBuilder.WriteString("  -c, --config <path>  Specify a custom path for the configuration file.\n")
	expectedOutputBuilder.WriteString("  --json               Output weather data in JSON format.\n")
	expectedOutputBuilder.WriteString("  -h, --help           Display this help message and exit.\n")
	expectedOutputBuilder.WriteString("\n")

	configDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("Error getting user config directory: %v", err)
	}
	defaultConfigPath := filepath.Join(configDir, "wayther", "config.json")
	expectedOutputBuilder.WriteString(fmt.Sprintf("Default configuration path: %s\n", defaultConfigPath))
	expectedOutputBuilder.WriteString("\n")
	expectedOutputBuilder.WriteString("Configuration:\n")
	expectedOutputBuilder.WriteString("  The application uses a configuration file to store your WeatherAPI key and default location.\n")
	expectedOutputBuilder.WriteString("  If no configuration file is found, you will be prompted to create one interactively.\n")
	expectedOutputBuilder.WriteString("  The 'logger' key in the config (boolean, defaults to false) enables syslog output if true.\n")

	expectedOutput := expectedOutputBuilder.String()

	if actualOutput != expectedOutput {
		t.Errorf("Help output mismatch:\nExpected:\n%s\nActual:\n%s", expectedOutput, actualOutput)
	}
}

func TestRunApp_OutputConsistency_JSON(t *testing.T) {
	// Simulate non-interactive mode
	oldIsTerminal := isTerminalCheck // Use isTerminalCheck from main.go
	isTerminalCheck = func(fd uintptr) bool { return false }
	defer func() { isTerminalCheck = oldIsTerminal }() // Restore original after test

	// Set a fixed time for deterministic testing, matching the start time in samples/output.json
	fixedTime := time.Date(2023, time.March, 15, 20, 0, 0, 0, time.UTC) // 20:00 on March 15, 2023
	oldNowFunc := nowFunc
	nowFunc = func() time.Time { return fixedTime }
	defer func() { nowFunc = oldNowFunc }() // Restore original nowFunc after test

	// Set up mock config
	mockConfig = &Config{
		APIKey:   "mock_api_key",
		Location: "MockCity",
	}

	// Set up mock weather response
	// This mock data should be consistent with samples/output.json
	mockWeatherResponse = &WeatherAPIResponse{
		Location: Location{Name: "MockCity"},
		Current: Current{
			TempC:    28.0,
			Humidity: 60,
			Condition: Condition{
				Text: "Clear",
				Code: 1000, // Assuming 1000 is clear for emoji mapping
			},
		},
		Forecast: Forecast{
			Forecastday: []Forecastday{
				{
					Hour: []Hour{
						{TimeEpoch: fixedTime.Add(time.Hour * 0).Unix(), TempC: 23.6, FeelslikeC: 25.1, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 1).Unix(), TempC: 22.8, FeelslikeC: 24.8, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 2).Unix(), TempC: 22.1, FeelslikeC: 24.6, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 3).Unix(), TempC: 21.5, FeelslikeC: 21.5, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 4).Unix(), TempC: 20.8, FeelslikeC: 20.8, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 5).Unix(), TempC: 20.2, FeelslikeC: 20.2, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 6).Unix(), TempC: 19.7, FeelslikeC: 19.7, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 7).Unix(), TempC: 19.3, FeelslikeC: 19.3, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 8).Unix(), TempC: 18.9, FeelslikeC: 18.9, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 9).Unix(), TempC: 19.1, FeelslikeC: 19.2, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 10).Unix(), TempC: 22.0, FeelslikeC: 22.0, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 11).Unix(), TempC: 24.7, FeelslikeC: 25.3, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 12).Unix(), TempC: 27.5, FeelslikeC: 27.0, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 13).Unix(), TempC: 29.8, FeelslikeC: 28.9, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 14).Unix(), TempC: 31.8, FeelslikeC: 30.7, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 15).Unix(), TempC: 33.4, FeelslikeC: 32.2, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 16).Unix(), TempC: 34.4, FeelslikeC: 33.2, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 17).Unix(), TempC: 35.1, FeelslikeC: 33.8, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 18).Unix(), TempC: 35.5, FeelslikeC: 34.1, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 19).Unix(), TempC: 35.5, FeelslikeC: 33.9, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 20).Unix(), TempC: 35.3, FeelslikeC: 33.8, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 21).Unix(), TempC: 34.1, FeelslikeC: 33.3, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 22).Unix(), TempC: 30.8, FeelslikeC: 30.2, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 23).Unix(), TempC: 27.8, FeelslikeC: 27.4, Condition: Condition{Code: 1000}},
					},
				},
			},
		},
	}

	// Read expected output from samples/output.json
	expectedOutputPath := filepath.Join("samples", "output.json")
	expectedOutputBytes, err := ioutil.ReadFile(expectedOutputPath)
	if err != nil {
		t.Fatalf("Failed to read expected output file: %v", err)
	}
	expectedOutput := strings.TrimSpace(string(expectedOutputBytes))

	// Construct the 'text' field for logging
	emoji := GetEmoji(mockWeatherResponse.Current.Condition.Code)
	t.Logf("Emoji returned by GetEmoji for code %d: %s", mockWeatherResponse.Current.Condition.Code, emoji)

	// Call runApp with mock dependencies
	actualOutput, err := runApp([]string{"cmd", "MockCity"}, mockConfig, mockGetWeather)
	if err != nil {
		t.Fatalf("runApp returned an error: %v", err)
	}

	// Trim trailing newline from actualOutput, as fmt.Println adds one
	actualOutput = strings.TrimSuffix(actualOutput, "\n")
	actualOutput = strings.TrimSuffix(actualOutput, "\r") // Also trim carriage return if present

	// Compare outputs
	if actualOutput != expectedOutput {
		t.Errorf("Output mismatch:\nExpected: %s\nActual:   %s", expectedOutput, actualOutput)
	}
}

func TestRunApp_OutputConsistency_Table(t *testing.T) {
	// Simulate interactive mode
	oldIsTerminal := isTerminalCheck
	isTerminalCheck = func(fd uintptr) bool { return true }
	defer func() { isTerminalCheck = oldIsTerminal }() // Restore original after test

	// Set a fixed time for deterministic testing
	fixedTime := time.Date(2023, time.March, 15, 0, 0, 0, 0, time.UTC)
	oldNowFunc := nowFunc
	nowFunc = func() time.Time { return fixedTime }
	defer func() { nowFunc = oldNowFunc }() // Restore original nowFunc after test

	// Set up mock config
	mockConfig = &Config{
		APIKey:   "mock_api_key",
		Location: "Athens",
	}

	// Set up mock weather response
	mockWeatherResponse = &WeatherAPIResponse{
		Location: Location{Name: "Athens", Country: "Greece"},
		Current: Current{
			TempC:    24.0,
			Humidity: 60,
			Condition: Condition{
				Text: "Clear",
				Code: 1000,
			},
		},
		Forecast: Forecast{
			Forecastday: []Forecastday{
				{
					Hour: []Hour{
						{TimeEpoch: fixedTime.Add(time.Hour * 0).Unix(), TempC: 24.0, FeelslikeC: 20.8, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 1).Unix(), TempC: 20.2, FeelslikeC: 20.2, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 2).Unix(), TempC: 19.7, FeelslikeC: 19.7, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 3).Unix(), TempC: 19.3, FeelslikeC: 19.3, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 4).Unix(), TempC: 18.9, FeelslikeC: 18.9, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 5).Unix(), TempC: 19.1, FeelslikeC: 19.2, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 6).Unix(), TempC: 22.0, FeelslikeC: 22.0, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 7).Unix(), TempC: 24.7, FeelslikeC: 25.3, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 8).Unix(), TempC: 27.5, FeelslikeC: 27.0, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 9).Unix(), TempC: 29.8, FeelslikeC: 28.9, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 10).Unix(), TempC: 31.8, FeelslikeC: 30.7, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 11).Unix(), TempC: 33.4, FeelslikeC: 32.2, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 12).Unix(), TempC: 34.4, FeelslikeC: 33.2, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 13).Unix(), TempC: 35.1, FeelslikeC: 33.8, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 14).Unix(), TempC: 35.5, FeelslikeC: 34.1, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 15).Unix(), TempC: 35.5, FeelslikeC: 33.9, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 16).Unix(), TempC: 35.3, FeelslikeC: 33.8, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 17).Unix(), TempC: 34.1, FeelslikeC: 33.3, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 18).Unix(), TempC: 30.8, FeelslikeC: 30.2, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 19).Unix(), TempC: 27.8, FeelslikeC: 27.4, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 20).Unix(), TempC: 26.5, FeelslikeC: 26.4, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 21).Unix(), TempC: 24.9, FeelslikeC: 25.5, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 22).Unix(), TempC: 23.6, FeelslikeC: 24.9, Condition: Condition{Code: 1000, Text: "Clear"}},
						{TimeEpoch: fixedTime.Add(time.Hour * 23).Unix(), TempC: 22.7, FeelslikeC: 24.7, Condition: Condition{Code: 1000, Text: "Clear"}},
					},
				},
			},
		},
	}

	// Read expected table output from samples/output.text
	expectedTableOutputPath := filepath.Join("samples", "output.text")
	expectedTableOutputBytes, err := ioutil.ReadFile(expectedTableOutputPath)
	if err != nil {
		t.Fatalf("Failed to read expected table output file: %v", err)
	}
	expectedTableOutput := string(expectedTableOutputBytes) // No TrimSpace here, as it's character-perfect

	// Call runApp with mock dependencies
	actualOutput, err := runApp([]string{"cmd", "Sandweiler"}, mockConfig, mockGetWeather)
	if err != nil {
		t.Fatalf("runApp returned an error: %v", err)
	}

	// Compare outputs
	if actualOutput != expectedTableOutput { // Direct comparison
		t.Errorf("Table output mismatch:\nExpected:\n%s\nActual:\n%s", expectedTableOutput, actualOutput)
	}
}

func TestRunApp_JsonFlag(t *testing.T) {
	// Simulate interactive mode but with --json flag
	oldIsTerminal := isTerminalCheck
	isTerminalCheck = func(fd uintptr) bool { return true } // Simulate TTY
	defer func() { isTerminalCheck = oldIsTerminal }()

	// Set a fixed time for deterministic testing
	fixedTime := time.Date(2023, time.March, 15, 20, 0, 0, 0, time.UTC)
	oldNowFunc := nowFunc
	nowFunc = func() time.Time { return fixedTime }
	defer func() { nowFunc = oldNowFunc }()

	// Set up mock config
	mockConfig = &Config{
		APIKey:   "mock_api_key",
		Location: "MockCity",
	}

	// Set up mock weather response (same as JSON test)
	mockWeatherResponse = &WeatherAPIResponse{
		Location: Location{Name: "MockCity"},
		Current: Current{
			TempC:    28.0,
			Humidity: 60,
			Condition: Condition{
				Text: "Clear",
				Code: 1000,
			},
		},
		Forecast: Forecast{
			Forecastday: []Forecastday{
				{
					Hour: []Hour{
						{TimeEpoch: fixedTime.Add(time.Hour * 0).Unix(), TempC: 23.6, FeelslikeC: 25.1, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 1).Unix(), TempC: 22.8, FeelslikeC: 24.8, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 2).Unix(), TempC: 22.1, FeelslikeC: 24.6, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 3).Unix(), TempC: 21.5, FeelslikeC: 21.5, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 4).Unix(), TempC: 20.8, FeelslikeC: 20.8, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 5).Unix(), TempC: 20.2, FeelslikeC: 20.2, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 6).Unix(), TempC: 19.7, FeelslikeC: 19.7, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 7).Unix(), TempC: 19.3, FeelslikeC: 19.3, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 8).Unix(), TempC: 18.9, FeelslikeC: 18.9, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 9).Unix(), TempC: 19.1, FeelslikeC: 19.2, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 10).Unix(), TempC: 22.0, FeelslikeC: 22.0, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 11).Unix(), TempC: 24.7, FeelslikeC: 25.3, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 12).Unix(), TempC: 27.5, FeelslikeC: 27.0, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 13).Unix(), TempC: 29.8, FeelslikeC: 28.9, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 14).Unix(), TempC: 31.8, FeelslikeC: 30.7, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 15).Unix(), TempC: 33.4, FeelslikeC: 32.2, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 16).Unix(), TempC: 34.4, FeelslikeC: 33.2, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 17).Unix(), TempC: 35.1, FeelslikeC: 33.8, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 18).Unix(), TempC: 35.5, FeelslikeC: 34.1, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 19).Unix(), TempC: 35.5, FeelslikeC: 33.9, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 20).Unix(), TempC: 35.3, FeelslikeC: 33.8, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 21).Unix(), TempC: 34.1, FeelslikeC: 33.3, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 22).Unix(), TempC: 30.8, FeelslikeC: 30.2, Condition: Condition{Code: 1000}},
						{TimeEpoch: fixedTime.Add(time.Hour * 23).Unix(), TempC: 27.8, FeelslikeC: 27.4, Condition: Condition{Code: 1000}},
					},
				},
			},
		},
	}

	// Read expected output from samples/output.json
	expectedOutputPath := filepath.Join("samples", "output.json")
	expectedOutputBytes, err := ioutil.ReadFile(expectedOutputPath)
	if err != nil {
		t.Fatalf("Failed to read expected output file: %v", err)
	}
	expectedOutput := strings.TrimSpace(string(expectedOutputBytes))

	// Call runApp with --json flag
	actualOutput, err := runApp([]string{"cmd", "MockCity", "--json"}, mockConfig, mockGetWeather)
	if err != nil {
		t.Fatalf("runApp returned an error: %v", err)
	}

	// Trim trailing newline from actualOutput
	actualOutput = strings.TrimSuffix(actualOutput, "\n")
	actualOutput = strings.TrimSuffix(actualOutput, "\r")

	// Compare outputs
	if actualOutput != expectedOutput {
		t.Errorf("Output mismatch:\nExpected: %s\nActual:   %s", expectedOutput, actualOutput)
	}
}

func TestRunApp_NoLocation(t *testing.T) {
	// Set up mock config with no location
	mockConfig = &Config{
		APIKey:   "mock_api_key",
		Location: "",
	}

	// Call runApp with no location argument
	_, err := runApp([]string{"cmd"}, mockConfig, mockGetWeather)
	if err == nil {
		t.Fatal("runApp should have returned an error, but it didn't")
	}

	// Check the error message
	expectedError := "no location provided. Please provide a location as an argument or set a default in config.json"
	if err.Error() != expectedError {
		t.Errorf("Unexpected error message:\nExpected: %s\nActual:   %s", expectedError, err.Error())
	}
}
