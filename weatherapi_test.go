package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeatherProvider_GetWeather(t *testing.T) {


	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request path and query parameters
		if r.URL.Path != "/v1/forecast.json" {
			t.Errorf("Expected to request '/v1/forecast.json', got: %s", r.URL.Path)
		}
		if r.URL.Query().Get("q") != "London" {
			t.Errorf("Expected query parameter 'q' to be 'London', got: %s", r.URL.Query().Get("q"))
		}
		if r.URL.Query().Get("key") != "test_api_key" {
			t.Errorf("Expected query parameter 'key' to be 'test_api_key', got: %s", r.URL.Query().Get("key"))
		}

		// Provide a sample JSON response
		sampleResponse := WeatherAPIResponse{
			Location: Location{Name: "London"},
			Current: Current{
				TempC:    10.0,
				Humidity: 70,
				Condition: Condition{
					Text: "Partly cloudy",
					Code: 1003,
				},
			},
		}
		json.NewEncoder(w).Encode(sampleResponse)
	}))
	defer server.Close()

	// Call the function under test
	originalURL := weatherAPIURL
	weatherAPIURL = server.URL + "/v1/forecast.json"
	defer func() { weatherAPIURL = originalURL }()

	tempDir, err := os.MkdirTemp("", "cache-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cache, err := NewCache(filepath.Join(tempDir, "cache.json"))
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	provider := &weatherapiProvider{cache: cache}
	config := &Config{
		Location: "London",
		APIKey:   "test_api_key",
	}
	weather, err := provider.GetWeather(config)
	if err != nil {
		t.Fatalf("GetWeather returned an error: %v", err)
	}

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "London", weather.Location.Name)
	assert.Equal(t, 10.0, weather.Current.TempC)
	assert.Equal(t, 70, weather.Current.Humidity)
	assert.Equal(t, "Partly cloudy", weather.Current.Condition.Text)
	assert.Equal(t, 1003, weather.Current.Condition.Code)
}
