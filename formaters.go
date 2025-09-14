package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

// formatTable formats the weather data into a human-readable table
func formatTable(weather *WeatherAPIResponse, config *Config, nowFunc func() time.Time) string {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)

	// Current section
	t.AppendRow(table.Row{"Current:"})
	t.AppendSeparator()
	currentLine := fmt.Sprintf(config.Outputs.Table.CurrentFmt, weather.Current.Condition.Emoji, weather.Current.TempC, weather.Location.Name, weather.Location.Country)
	t.AppendRow(table.Row{currentLine})

	// Hourly Forecast section
	t.AppendSeparator()
	t.AppendRow(table.Row{"Hourly Forecast:"})
	t.AppendSeparator()

	// Hourly forecast details
	if len(weather.Forecast.Forecastday) > 0 {
		for _, forecastday := range weather.Forecast.Forecastday {
			for _, hour := range forecastday.Hour {
				timeVal := time.Unix(hour.TimeEpoch, 0)
				if timeVal.Before(nowFunc()) {
					continue
				}

				hourlyLine := fmt.Sprintf(config.Outputs.Table.ForecastFmt, timeVal.Format("15:04"), hour.Condition.Emoji, hour.TempC, hour.FeelslikeC)
				t.AppendRow(table.Row{hourlyLine})

				//we need this to restrict the results to 24hours
				if timeVal.After(nowFunc().Add(time.Hour * 23)) {
					break
				}
			}
		}
	}

	return t.Render()
}

// formatJSON formats the weather data into a JSON string
func formatJSON(weather *WeatherAPIResponse, config *Config, nowFunc func() time.Time) string {
	// Construct the 'text' field
	text := fmt.Sprintf(config.Outputs.JSON.CurrentFmt, weather.Current.Condition.Emoji, weather.Current.TempC)

	// Construct the 'tooltip' field
	tooltip := []string{}
	if len(weather.Forecast.Forecastday) > 0 {
		for _, forecastday := range weather.Forecast.Forecastday {
			for _, hour := range forecastday.Hour {
				timeVal := time.Unix(hour.TimeEpoch, 0)
					if timeVal.Before(nowFunc()) {
						continue
					}

					tooltip = append(tooltip, fmt.Sprintf(config.Outputs.JSON.ForecastFmt, timeVal.Format("15:04"), hour.Condition.Emoji, hour.TempC, hour.FeelslikeC))

					if timeVal.After(nowFunc().Add(time.Hour * 23)) {
						break
					}
			}
		}
	}

	// Create the final output struct
	outputStruct := struct {
		Text    string `json:"text"`
		Tooltip string `json:"tooltip"`
	}{
		Text:    text,
		Tooltip: strings.Join(tooltip, "\r"),
	}

	// Marshal to JSON
	jsonOutput, err := json.Marshal(outputStruct)
	if err != nil {
		fmt.Errorf("error marshalling JSON output: %w", err)
	}

	return string(jsonOutput)
}
