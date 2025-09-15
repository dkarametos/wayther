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
	currentLine, err := renderTemplateToString("table-current", config.CurrentTmpl, weather.Current)
	if err != nil {
		// Log the error, but don't stop execution for formatting errors
		fmt.Printf("Error rendering current template: %v\n", err)
	}
	t.AppendRow(table.Row{currentLine})

	// Location section
	locationLine, err := renderTemplateToString("table-location", config.LocationTmpl, weather.Location)
	if err != nil {
		fmt.Printf("Error rendering location template: %v\n", err)
	}
	t.AppendRow(table.Row{locationLine})

	// Hourly Forecast section
	t.AppendSeparator()
	t.AppendRow(table.Row{"Hourly Forecast:"})
	t.AppendSeparator()

	// Hourly forecast details
	if len(weather.HourlyForecast) > 0 {
		for _, hour := range weather.HourlyForecast {
			timeVal := time.Unix(hour.TimeEpoch, 0)
			if timeVal.Before(nowFunc()) {
				continue
			}

			hourlyLineContent, err := renderTemplateToString("table-hourly", config.ForecastTmpl, hour)
			if err != nil {
				fmt.Printf("Error rendering hourly template: %v\n", err)
			}
			hourlyLine := fmt.Sprintf("%s : %s", timeVal.Format("15:04"), hourlyLineContent)
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
	text, err := renderTemplateToString("json-text", config.CurrentTmpl, weather.Current)
	if err != nil {
		fmt.Printf("Error rendering json text template: %v\n", err)
		text = "" // Fallback to empty string on error
	}

	// Construct the 'tooltip' field
	tooltip := []string{}
	if len(weather.HourlyForecast) > 0 {
		for _, hour := range weather.HourlyForecast {
			timeVal := time.Unix(hour.TimeEpoch, 0)
			if timeVal.Before(nowFunc()) {
				continue
			}

			tooltipLineContent, err := renderTemplateToString("json-tooltip", config.ForecastTmpl, hour)
			if err != nil {
				fmt.Printf("Error rendering json tooltip template: %v\n", err)
				tooltipLineContent = "" // Fallback to empty string on error
			}
			tooltipLine := fmt.Sprintf("%s: %s", timeVal.Format("15:04"), tooltipLineContent)
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
		Text:    text,
		Tooltip: strings.Join(tooltip, "\r"),
	}

	// Marshal to JSON
	jsonOutput, err := json.Marshal(outputStruct)
	if err != nil {
		fmt.Printf("Error marshalling JSON output: %v\n", err)
		return "{}" // Fallback to empty JSON object on error
	}

	return string(jsonOutput)
}

// renderTemplateToString parses and executes a template, returning the result as a string.
func renderTemplateToString(
	templateName string,
	templateString string,
	data interface{},
) (string, error) {
	tmpl, err := template.New(templateName).Parse(templateString)
	if err != nil {
		return "", fmt.Errorf("error creating template %s: %w", templateName, err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("error executing template %s: %w", templateName, err)
	}
	return buf.String(), nil
}
