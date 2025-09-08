package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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

	// Temporarily override the weatherAPIURL for testing
	oldWeatherAPIURL := weatherAPIURL
	weatherAPIURL = server.URL + "/v1/forecast.json" // Adjust to match the mock server's path
	defer func() { weatherAPIURL = oldWeatherAPIURL }()

	// Call the function under test
	provider := &weatherapiProvider{}
	weather, err := provider.GetWeather("London", "test_api_key")
	if err != nil {
		t.Fatalf("GetWeather returned an error: %v", err)
	}

	// Assertions
	if weather.Location.Name != "London" {
		t.Errorf("Expected location name 'London', got: %s", weather.Location.Name)
	}
	if weather.Current.TempC != 10.0 {
		t.Errorf("Expected temperature 10.0, got: %.1f", weather.Current.TempC)
	}
	if weather.Current.Humidity != 70 {
		t.Errorf("Expected humidity 70, got: %d", weather.Current.Humidity)
	}
	if weather.Current.Condition.Text != "Partly cloudy" {
		t.Errorf("Expected condition 'Partly cloudy', got: %s", weather.Current.Condition.Text)
	}
	if weather.Current.Condition.Code != 1003 {
		t.Errorf("Expected condition code 1003, got: %d", weather.Current.Condition.Code)
	}
}
