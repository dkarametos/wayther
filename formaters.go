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

// FormatOutput formats the weather data based on the specified output type in the config.
func FormatOutput(weather *Weather, config *Config, nowFunc func() time.Time) (string, error) {

	if config.Output == "json" {
		return formatJSON(weather, config, nowFunc)
	}
	return formatTable(weather, config, nowFunc)
}

// formatTable formats the weather data into a human-readable table.
func formatTable(weather *Weather, config *Config, nowFunc func() time.Time) (string, error) {

	t := table.NewWriter()
	t.SetStyle(table.StyleLight)

	// Current section
	t.AppendRow(table.Row{"Current:"})
	t.AppendSeparator()

	currentLine, err := renderTemplateToString("table-current", config.CurrentTmpl, weather.Current)
	if err != nil {
		return "", fmt.Errorf("error rendering location template: %w", err)
	}
	t.AppendRow(table.Row{currentLine})

	// Hourly Forecast section
	if err := renderHourlyForecast(t, weather, config, nowFunc); err != nil {
		return "", err
	}

	return t.Render(), nil
}

// renderHourlyForecast renders the hourly forecast section of the table.
func renderHourlyForecast(t table.Writer, weather *Weather, config *Config, nowFunc func() time.Time) error {

	if config.ForecastHours > 0 {
		t.AppendSeparator()
		t.AppendRow(table.Row{"Hourly Forecast:"})
		t.AppendSeparator()

		hoursCount := 0
		for _, hour := range weather.HourlyForecast {
			if hoursCount >= config.ForecastHours {
				break
			}
			timeVal := time.Unix(hour.TimeEpoch, 0)
			if timeVal.Before(nowFunc()) {
				continue
			}

			hourlyLineContent, err := renderTemplateToString("table-hourly", config.ForecastTmpl, hour)
			if err != nil {
				return fmt.Errorf("error rendering hourly template: %w", err)
			}
			hourlyLine := fmt.Sprintf("%s : %s", timeVal.Format("15:04"), hourlyLineContent)
			t.AppendRow(table.Row{hourlyLine})
			hoursCount++

			if timeVal.After(nowFunc().Add(time.Hour * 23)) {
				break
			}
		}
	}
	return nil
}

// formatJSON formats the weather data into a JSON string.
func formatJSON(weather *Weather, config *Config, nowFunc func() time.Time) (string, error) {

	text, err := renderTemplateToString("json-text", config.ShortTmpl, weather.Current)
	if err != nil {
		return "", fmt.Errorf("error rendering json text template: %w", err)
	}

	tooltipContent, err := renderJSONTooltip(weather, config, nowFunc)
	if err != nil {
		return "", err
	}

	// Create the final output struct
	outputStruct := struct {
		Text    string `json:"text"`
		Tooltip string `json:"tooltip"`
	}{
		Text:    text,
		Tooltip: tooltipContent,
	}

	// Marshal to JSON
	jsonOutput, err := json.Marshal(outputStruct)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON output: %w", err)
	}

	return string(jsonOutput), nil
}

// renderJSONTooltip renders the JSON tooltip field.
func renderJSONTooltip(weather *Weather, config *Config, nowFunc func() time.Time) (string, error) {

	tooltip := []string{}
	if config.ForecastHours > 0 {
		hoursCount := 0
		for _, hour := range weather.HourlyForecast {
			if hoursCount >= config.ForecastHours {
				break
			}
			timeVal := time.Unix(hour.TimeEpoch, 0)
			if timeVal.Before(nowFunc()) {
				continue
			}

			tooltipLineContent, err := renderTemplateToString("json-tooltip", config.ForecastTmpl, hour)
			if err != nil {
				return "", fmt.Errorf("error rendering json tooltip template: %w", err)
			}
			tooltipLine := fmt.Sprintf(" %s: %s ", timeVal.Format("15:04"), tooltipLineContent)
			tooltip = append(tooltip, tooltipLine)
			hoursCount++


			if timeVal.After(nowFunc().Add(time.Hour * 23)) {
				break
			}
		}
	}
	return strings.Join(tooltip, "\r"), nil
}

// renderTemplateToString parses and executes a template, returning the result as a string.
func renderTemplateToString(templateName string, templateString string, data interface{},) (string, error) {

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
