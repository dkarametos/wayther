# Templates

Wayther uses Go templates to allow for flexible and customizable output. There are three main template sections: `short`, `current`, and `forecast`.

## Template Usage

The application uses different templates depending on the chosen output format:

*   **JSON Output:** Uses the `short_template` for the `text` field and the `forecast_template` for the `tooltip` field.
*   **Table Output:** Uses the `current_template` for the current weather summary and the `forecast_template` for the hourly forecast.

## Template Examples

### Table Output Example

When using the `table` output format, the `current_template` and `forecast_template` are utilized.

```
┌───────────────────────────┐
│ Current:                  │
├───────────────────────────┤
│ 🌦️ 24.0°                  │
│ Athens - Greece           │
├───────────────────────────┤
│ Hourly Forecast:          │
├───────────────────────────┤
│ 00:00:  🌦️ 24.0° [ 20.8°] │
│ 01:00:  ☀️ 20.2° [ 20.2°] │
│ 02:00:  ☀️ 19.7° [ 19.7°] │
│ 03:00:  ☀️ 19.3° [ 19.3°] │
│ 04:00:  ☀️ 18.9° [ 18.9°] │
└───────────────────────────┘
```
*   The section marked `Current:` is formatted using the `current_template`.
*   Each line under `Hourly Forecast:` is formatted using the `forecast_template`.

### JSON Output Example

When using the `json` output format, the `short_template` and `forecast_template` are utilized.

```json
{
  "text": "☁️ 1.3°C",
  "tooltip": " 00:00: ❄️ -0.4°C [-3.5°C]\r 01:00: ☁️ -1.1°C [-4.4°C]\r 02:00: ☁️ -1.4°C [-4.8°C]\r 03:00: ☁️ -1.3°C [-5.2°C]"
}
```
*   The `text` field is formatted using the `short_template`.
*   The `tooltip` field is formatted using the `forecast_template` for each hourly entry, joined by `\r`.

## Available Template Elements

You can use the following elements within your templates:

### For `current_template` and `short_template` (based on `Weather.Current`):

*   `.Location`: The name of the location (string).
*   `.Country`: The country of the location (string).
*   `.Emoji`: An emoji representing the current weather condition (string).
*   `.TempC`: The current temperature in Celsius (float64).

### For `forecast_template` (based on `HourlyForecast`):

*   `.TimeEpoch`: Unix timestamp for the forecast hour (int64).
*   `.Emoji`: An emoji representing the hourly weather condition (string).
*   `.TempC`: The temperature in Celsius for the hour (float64).
*   `.FeelslikeC`: The "feels like" temperature in Celsius for the hour (float64).

You can also use Go template functions like `printf` for formatting numbers. For example, `{{printf "%.1f" .TempC}}` will format `TempC` to one decimal place.

```
