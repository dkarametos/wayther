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

	"github.com/spf13/cobra"
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

type MockWeatherProvider struct {
	mockResponse *WeatherAPIResponse
	err          error
}

func (m *MockWeatherProvider) GetWeather(config *Config) (*WeatherAPIResponse, error) {
	return m.mockResponse, m.err
}

func (m *MockWeatherProvider) ToWeather(w *WeatherAPIResponse) *Weather {
	var hourlyForecasts []HourlyForecast
	if len(w.Forecast.Forecastday) > 0 {
		for _, forecastday := range w.Forecast.Forecastday {
			for _, hour := range forecastday.Hour {
				hourlyForecasts = append(hourlyForecasts, HourlyForecast{
					TimeEpoch:  hour.TimeEpoch,
					Emoji:      hour.Condition.Emoji,
					TempC:      hour.TempC,
					FeelslikeC: hour.FeelslikeC,
				})
			}
		}
	}

	return &Weather{
		Current: WeatherCurrent{
			Location: w.Location.Name,
			Country:  w.Location.Country,
			Emoji:    w.Current.Condition.Emoji,
			TempC:    w.Current.TempC,
		},
		HourlyForecast: hourlyForecasts,
	}
}

func (m *MockWeatherProvider) CleanCache(maxAge time.Duration) {
	// Mock implementation - do nothing or log if needed for testing cache cleaning logic
}

type MockConfigProvider struct {
	mockConfig *Config
	err        error
}

func (m *MockConfigProvider) LoadConfig(configPath ConfigPath) (*Config, error) {
	return m.mockConfig, m.err
}

func TestAppOutput(t *testing.T) {
	// --- Setup Mocks ---
	mockResponse := loadMockResponse(t)
	weatherProvider := &MockWeatherProvider{mockResponse: mockResponse}
	configProvider := &MockConfigProvider{mockConfig: &Config{
		APIKey:   "mock-key",
		Location: "Brussels",
		CurrentTmpl:  "{{.Emoji}}  {{.TempC}}°",
		ForecastTmpl: "{{.Emoji}} {{.TempC}}° [{{.FeelslikeC}}°]",
		
	}}

	mockNowFunc := func() time.Time {
		return time.Unix(mockResponse.Location.LocaltimeEpoch, 0)
	}

	t.Run("JSON Output", func(t *testing.T) {
		// Redirect stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		configProvider.mockConfig.OutputType = "json"
		cmd := &cobra.Command{}
		err := runApp(cmd, []string{"Brussels"}, ConfigPath{}, weatherProvider, configProvider, false, mockNowFunc)
		assert.NoError(t, err)

		// Restore stdout and read the captured output
		w.Close()
		os.Stdout = oldStdout
		var buf bytes.Buffer
		io.Copy(&buf, r)
		actualOutput := buf.String()

		// The app produces valid JSON, so we check for key substrings.
		assert.Contains(t, actualOutput, "\"text\":", "Output should contain the JSON key 'text'")
		assert.Contains(t, actualOutput, "1.3°", "JSON output should contain the current temperature for Brussels")
		assert.Contains(t, actualOutput, "\"tooltip\":", "Output should contain the JSON key 'tooltip'")
		assert.Contains(t, actualOutput, "\"tooltip\":\"\"", "JSON tooltip should be empty")
	})

	t.Run("Table Output", func(t *testing.T) {
		configProvider.mockConfig.ForecastHours = 4
		// Redirect stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		configProvider.mockConfig.OutputType = "table"
		cmd := &cobra.Command{}
		cmd.Flags().Int("forecast-hours", 4, "")
		err := runApp(cmd, []string{"Brussels"}, ConfigPath{}, weatherProvider, configProvider, true, mockNowFunc)
		assert.NoError(t, err)

		// Restore stdout and read the captured output
		w.Close()
		os.Stdout = oldStdout
		var buf bytes.Buffer
		io.Copy(&buf, r)
		actualOutput := buf.String()

		// Assert that key elements for the Brussels response are present
		assert.Contains(t, actualOutput, "Current:", "Table should have a 'Current' section")
		
		assert.Contains(t, actualOutput, "1.3°", "Table should contain the current temperature")
		assert.Contains(t, actualOutput, "Hourly Forecast:", "Table should have an 'Hourly Forecast' section")
		
	})
}

func TestExecutionError(t *testing.T) {
	t.Run("Weather Provider Error - Terminal Output", func(t *testing.T) {
		weatherProvider := &MockWeatherProvider{err: errors.New("mock weather error")}
		configProvider := &MockConfigProvider{mockConfig: &Config{OutputType: "table"}} // Simulate terminal output

		err := runApp(&cobra.Command{}, []string{"some-location"}, ConfigPath{}, weatherProvider, configProvider, true, time.Now)
		assert.Error(t, err)
		assert.EqualError(t, err, "mock weather error")
	})

	t.Run("Weather Provider Error - JSON Output", func(t *testing.T) {
		// Redirect stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		weatherProvider := &MockWeatherProvider{err: errors.New("mock weather error")}
		configProvider := &MockConfigProvider{mockConfig: &Config{OutputType: "json"}} // Simulate JSON output

		err := runApp(&cobra.Command{}, []string{"some-location"}, ConfigPath{}, weatherProvider, configProvider, false, time.Now)
		assert.NoError(t, err)

		// Restore stdout and read the captured output
		w.Close()
		os.Stdout = oldStdout
		var buf bytes.Buffer
		io.Copy(&buf, r)
		actualOutput := buf.String()

		assert.Contains(t, actualOutput, `{"text":"N/A ☢","tooltip":" error fetching weather: mock weather error "}`)
	})

	t.Run("Config Load Error - Terminal Output", func(t *testing.T) {
		weatherProvider := &MockWeatherProvider{}
		configProvider := &MockConfigProvider{err: errors.New("mock config load error"), mockConfig: &Config{OutputType: "table"}} // Simulate terminal output

		err := runApp(&cobra.Command{}, []string{"some-location"}, ConfigPath{}, weatherProvider, configProvider, true, time.Now)
		assert.Error(t, err)
		assert.EqualError(t, err, "mock config load error")
	})

	t.Run("Config Load Error - JSON Output", func(t *testing.T) {
		// Redirect stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		weatherProvider := &MockWeatherProvider{}
		configProvider := &MockConfigProvider{err: errors.New("mock config load error"), mockConfig: &Config{OutputType: "json"}} // Simulate JSON output

		err := runApp(&cobra.Command{}, []string{"some-location"}, ConfigPath{}, weatherProvider, configProvider, false, time.Now)
		assert.NoError(t, err)

		// Restore stdout and read the captured output
		w.Close()
		os.Stdout = oldStdout
		var buf bytes.Buffer
		io.Copy(&buf, r)
		actualOutput := buf.String()

		assert.Contains(t, actualOutput, `{"text":"N/A ☢","tooltip":" error fetching weather: mock config load error "}`)
	})

	t.Run("Config Load Error - Nil Config - Terminal Output", func(t *testing.T) {
		weatherProvider := &MockWeatherProvider{}
		configProvider := &MockConfigProvider{err: errors.New("mock config load error"), mockConfig: nil} // Simulate nil config

		err := runApp(&cobra.Command{}, []string{"some-location"}, ConfigPath{}, weatherProvider, configProvider, true, time.Now)
		assert.Error(t, err)
		assert.EqualError(t, err, "mock config load error")
	})

	t.Run("Config Load Error - Nil Config - Non-Terminal Output", func(t *testing.T) {
		// Redirect stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		weatherProvider := &MockWeatherProvider{}
		configProvider := &MockConfigProvider{err: errors.New("mock config load error"), mockConfig: nil} // Simulate nil config

		err := runApp(&cobra.Command{}, []string{"some-location"}, ConfigPath{}, weatherProvider, configProvider, false, time.Now)
		assert.NoError(t, err)

		// Restore stdout and read the captured output
		w.Close()
		os.Stdout = oldStdout
		var buf bytes.Buffer
		io.Copy(&buf, r)
		actualOutput := buf.String()

		assert.Contains(t, actualOutput, `{"text":"N/A ☢","tooltip":" error fetching weather: mock config load error "}`)
	})
}
