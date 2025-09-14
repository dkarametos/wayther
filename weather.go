package main

// WeatherProvider is an interface for fetching weather data.
type WeatherProvider interface {
	GetWeather(config *Config) (*WeatherAPIResponse, error)
	toWeather(w *WeatherAPIResponse) (*Weather)
}

// Weather holds the simplified weather data for formatting.
type Weather struct {
	LocationName    string
	LocationCountry string
	CurrentEmoji    string
	CurrentTempC    float64
	HourlyForecast  []HourlyForecast
}

// HourlyForecast holds the simplified hourly forecast data.
type HourlyForecast struct {
	TimeEpoch  int64
	Emoji      string
	TempC      float64
	FeelslikeC float64
}

// NewWeather creates a new Weather struct from the provider and config.
func NewWeather(provider WeatherProvider, config *Config) (*Weather, error) {
	weatherAPIResponse, err := provider.GetWeather(config)
	if err != nil {
		return nil, err
	}
	return provider.toWeather(weatherAPIResponse), nil
}
