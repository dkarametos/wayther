package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

// formatTable formats the weather data into a human-readable table
func formatTable(weather *Weather, config *Config, nowFunc func() time.Time) string {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)

	// Current section
	t.AppendRow(table.Row{"Current:"})
	t.AppendSeparator()
	// Create a new template and parse the template string
	currentTmpl, err := template.New("table-current").Parse(config.CurrentTmpl)
	if err != nil {
		fmt.Errorf("error creating table template: %w", err)
	}

	// Create a buffer to store the executed template
	var currentLine bytes.Buffer
	err = currentTmpl.Execute(&currentLine, weather)
	if err != nil {
		fmt.Errorf("error executing table template: %w", err)
	}
	t.AppendRow(table.Row{currentLine.String()})

	// Location section
	locationTmpl, err := template.New("table-location").Parse(config.LocationTmpl)
	if err != nil {
		fmt.Errorf("error creating table template: %w", err)
	}
	var locationLine bytes.Buffer
	err = locationTmpl.Execute(&locationLine, weather)
	if err != nil {
		fmt.Errorf("error executing table template: %w", err)
	}
	t.AppendRow(table.Row{locationLine.String()})

	// Hourly Forecast section
	t.AppendSeparator()
	t.AppendRow(table.Row{"Hourly Forecast:"})
	t.AppendSeparator()

	// Hourly forecast details
	if len(weather.HourlyForecast) > 0 {
		// Create a new template and parse the template string
		hourlyTmpl, err := template.New("table-hourly").Parse(config.ForecastTmpl)
		if err != nil {
			fmt.Errorf("error creating table template: %w", err)
		}

		for _, hour := range weather.HourlyForecast {
			timeVal := time.Unix(hour.TimeEpoch, 0)
			if timeVal.Before(nowFunc()) {
				continue
			}

			// Create a buffer to store the executed template
			var hourlyLineBts bytes.Buffer
			err = hourlyTmpl.Execute(&hourlyLineBts, hour)
			if err != nil {
				fmt.Errorf("error executing table template: %w", err)
			}

			hourlyLine := fmt.Sprintf("%s: %s", timeVal.Format("15:04"), hourlyLineBts.String())
			t.AppendRow(table.Row{hourlyLine})

			//we need this to restrict the results to 24hours
			if timeVal.After(nowFunc().Add(time.Hour * 23)) {
				break
			}
		}
	}

	return t.Render()
}

// formatJSON formats the weather data into a JSON string
func formatJSON(weather *Weather, config *Config, nowFunc func() time.Time) string {

	// Create a new template and parse the template string
	t, err := template.New("json-text").Parse(config.CurrentTmpl)
	if err != nil {
		fmt.Errorf("error creating json template: %w", err)
	}

	// Create a buffer to store the executed template
	var text bytes.Buffer
	err = t.Execute(&text, weather)
	if err != nil {
		fmt.Errorf("error executing json template: %w", err)
	}

	// Construct the 'tooltip' field
	tooltip := []string{}
	if len(weather.HourlyForecast) > 0 {

		// Create a new template and parse the template string
		t, err := template.New("json-tooltip").Parse(config.ForecastTmpl)
		if err != nil {
			fmt.Errorf("error creating json template: %w", err)
		}

		for _, hour := range weather.HourlyForecast {
			timeVal := time.Unix(hour.TimeEpoch, 0)
			if timeVal.Before(nowFunc()) {
				continue
			}

			// Create a buffer to store the executed template
			var tooltipLineBts bytes.Buffer
			err = t.Execute(&tooltipLineBts, hour)
			if err != nil {
				fmt.Errorf("error executing json template: %w", err)
			}

	tooltipLine := fmt.Sprintf("%s: %s", timeVal.Format("15:04"), tooltipLineBts.String())
			tooltip = append(tooltip, tooltipLine)

			if timeVal.After(nowFunc().Add(time.Hour * 23)) {
				break
			}
		}
	}

	// Create the final output struct
	outputStruct := struct {
		Text    string `json:"text"`
		Tooltip string `json:"tooltip"`
	}{
		Text:    text.String(),
		Tooltip: strings.Join(tooltip, "\r"),
	}

	// Marshal to JSON
	jsonOutput, err := json.Marshal(outputStruct)
	if err != nil {
		fmt.Errorf("error marshalling JSON output: %w", err)
	}

	return string(jsonOutput)
}
