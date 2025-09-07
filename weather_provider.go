package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// WeatherProvider is an interface for fetching weather data.
type WeatherProvider interface {
	GetWeatherAPI(location, apiKey string) (*WeatherAPIResponse, error)
}

// APIWeatherProvider is the real implementation of WeatherProvider that uses the weather API.
type APIWeatherProvider struct{}

// GetWeatherAPI fetches weather forecast data from the WeatherAPI for a given location.
// It takes the location (e.g., "London") and an API key as input.
// It returns a pointer to a WeatherAPIResponse struct containing the parsed data,
// or an error if the request fails or the response cannot be decoded.
func (p *APIWeatherProvider) GetWeatherAPI(location, apiKey string) (*WeatherAPIResponse, error) {
	url := fmt.Sprintf("%s?key=%s&q=%s&days=2&aqi=no&alerts=no", weatherAPIURL, apiKey, location)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, resp.Status)
	}

	var weatherResp WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	// Populate emojis
	weatherResp.Current.Condition.Emoji = getEmojiForWeatherCode(weatherResp.Current.Condition.Code)

	for i := range weatherResp.Forecast.Forecastday {
		// Day condition
		weatherResp.Forecast.Forecastday[i].Day.Condition.Emoji = getEmojiForWeatherCode(weatherResp.Forecast.Forecastday[i].Day.Condition.Code)

		// Hourly conditions
		for j := range weatherResp.Forecast.Forecastday[i].Hour {
			weatherResp.Forecast.Forecastday[i].Hour[j].Condition.Emoji = getEmojiForWeatherCode(weatherResp.Forecast.Forecastday[i].Hour[j].Condition.Code)
		}
	}

	return &weatherResp, nil
}