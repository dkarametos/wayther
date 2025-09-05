package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var weatherAPIURL = "https://api.weatherapi.com/v1/forecast.json"

// WeatherAPIResponse represents the top-level structure of the WeatherAPI forecast.json response
type WeatherAPIResponse struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
	Forecast Forecast `json:"forecast"`
}

// Location represents the location data
type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TzID           string  `json:"tz_id"`
	LocaltimeEpoch int64   `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

// Current represents the current weather conditions
type Current struct {
	LastUpdatedEpoch int64     `json:"last_updated_epoch"`
	LastUpdated      string    `json:"last_updated"`
	TempC            float64   `json:"temp_c"`
	TempF            float64   `json:"temp_f"`
	IsDay            int       `json:"is_day"`
	Condition        Condition `json:"condition"`
	WindMph          float64   `json:"wind_mph"`
	WindKph          float64   `json:"wind_kph"`
	WindDegree       int       `json:"wind_degree"`
	WindDir          string    `json:"wind_dir"`
	PressureMb       float64   `json:"pressure_mb"`
	PressureIn       float64   `json:"pressure_in"`
	PrecipMm         float64   `json:"precip_mm"`
	PrecipIn         float64   `json:"precip_in"`
	Humidity         int       `json:"humidity"`
	Cloud            int       `json:"cloud"`
	FeelslikeC       float64   `json:"feelslike_c"`
	FeelslikeF       float64   `json:"feelslike_f"`
	WindchillC       float64   `json:"windchill_c"`
	WindchillF       float64   `json:"windchill_f"`
	HeatindexC       float64   `json:"heatindex_c"`
	HeatindexF       float64   `json:"heatindex_f"`
	DewpointC        float64   `json:"dewpoint_c"`
	DewpointF        float64   `json:"dewpoint_f"`
	VisKm            float64   `json:"vis_km"`
	VisMiles         float64   `json:"vis_miles"`
	Uv               float64   `json:"uv"`
	GustMph          float64   `json:"gust_mph"`
	GustKph          float64   `json:"gust_kph"`
}

// Condition represents the weather condition details
type Condition struct {
	Text  string `json:"text"`
	Icon  string `json:"icon"`
	Code  int    `json:"code"`
	Emoji string `json:"emoji,omitempty"`
}

// Forecast represents the forecast data
type Forecast struct {
	Forecastday []Forecastday `json:"forecastday"`
}

// Forecastday represents a single day in the forecast
type Forecastday struct {
	Date      string `json:"date"`
	DateEpoch int64  `json:"date_epoch"`
	Day       Day    `json:"day"`
	Astro     Astro  `json:"astro"`
	Hour      []Hour `json:"hour"`
}

// Day represents the daily summary
type Day struct {
	MaxtempC          float64   `json:"maxtemp_c"`
	MaxtempF          float64   `json:"maxtemp_f"`
	MintempC          float64   `json:"mintemp_c"`
	MintempF          float64   `json:"mintemp_f"`
	AvgtempC          float64   `json:"avgtemp_c"`
	AvgtempF          float64   `json:"avgtemp_f"`
	MaxwindMph        float64   `json:"maxwind_mph"`
	MaxwindKph        float64   `json:"maxwind_kph"`
	TotalprecipMm     float64   `json:"totalprecip_mm"`
	TotalprecipIn     float64   `json:"totalprecip_in"`
	TotalsnowCm       float64   `json:"totalsnow_cm"`
	AvgvisKm          float64   `json:"avgvis_km"`
	AvgvisMiles       float64   `json:"avgvis_miles"`
	Avghumidity       int       `json:"avghumidity"`
	DailyWillItRain   int       `json:"daily_will_it_rain"`
	DailyChanceOfRain int       `json:"daily_chance_of_rain"`
	DailyWillItSnow   int       `json:"daily_will_it_snow"`
	DailyChanceOfSnow int       `json:"daily_chance_of_snow"`
	Condition         Condition `json:"condition"`
	Uv                float64   `json:"uv"`
}

// Astro represents astronomical data
type Astro struct {
	Sunrise          string `json:"sunrise"`
	Sunset           string `json:"sunset"`
	Moonrise         string `json:"moonrise"`
	Moonset          string `json:"moonset"`
	MoonPhase        string `json:"moon_phase"`
	MoonIllumination   int    `json:"moon_illumination"`
	IsMoonUp         int    `json:"is_moon_up"`
	IsSunUp          int    `json:"is_sun_up"`
}

// Hour represents hourly forecast data
type Hour struct {
	TimeEpoch    int64     `json:"time_epoch"`
	Time         string    `json:"time"`
	TempC        float64   `json:"temp_c"`
	TempF        float64   `json:"temp_f"`
	IsDay        int       `json:"is_day"`
	Condition    Condition `json:"condition"`
	WindMph      float64   `json:"wind_mph"`
	WindKph      float64   `json:"wind_kph"`
	WindDegree   int       `json:"wind_degree"`
	WindDir      string    `json:"wind_dir"`
	PressureMb   float64   `json:"pressure_mb"`
	PressureIn   float64   `json:"pressure_in"`
	PrecipMm     float64   `json:"precip_mm"`
	PrecipIn     float64   `json:"precip_in"`
	SnowCm       float64   `json:"snow_cm"`
	Humidity     int       `json:"humidity"`
	Cloud        int       `json:"cloud"`
	FeelslikeC   float64   `json:"feelslike_c"`
	FeelslikeF   float64   `json:"feelslike_f"`
	WindchillC   float64   `json:"windchill_c"`
	WindchillF   float64   `json:"windchill_f"`
	HeatindexC   float64   `json:"heatindex_c"`
	HeatindexF   float64   `json:"heatindex_f"`
	DewpointC    float64   `json:"dewpoint_c"`
	DewpointF    float64   `json:"dewpoint_f"`
	WillItRain   int       `json:"will_it_rain"`
	ChanceOfRain int       `json:"chance_of_rain"`
	WillItSnow   int       `json:"will_it_snow"`
	ChanceOfSnow int       `json:"chance_of_snow"`
	VisKm        float64   `json:"vis_km"`
	VisMiles     float64   `json:"vis_miles"`
	GustMph      float64   `json:"gust_mph"`
	GustKph      float64   `json:"gust_kph"`
	Uv           float64   `json:"uv"`
	ShortRad     float64   `json:"short_rad"`
	DiffRad      float64   `json:"diff_rad"`
}

// GetWeather fetches weather forecast data from the WeatherAPI for a given location.
// It takes the location (e.g., "London") and an API key as input.
// It returns a pointer to a WeatherAPIResponse struct containing the parsed data,
// or an error if the request fails or the response cannot be decoded.
var GetWeatherAPI = func(location, apiKey string) (*WeatherAPIResponse, error) {
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
