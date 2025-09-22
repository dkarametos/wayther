# Configuration

This application uses the [weatherapi.com](https://www.weatherapi.com/) for weather data. You will need to obtain a free API key from their website.

The application now supports robust configuration merging:
*   **Default Configuration:** A default configuration file is located at `XDG_CONFIG_HOME/wayther/config.json` (typically `~/.config/wayther/config.json` on Linux).
*   **Custom Configurations:** You can specify a custom configuration file using the `-c` or `--config` flag.
*   **Merging Logic:** When a custom configuration is used, its values will override those in the in the default configuration. If the API key is missing in the custom configuration, the value from the default configuration will be used. This means keys in custom configurations are optional.

The first time you run the application without an existing configuration, it will prompt you to enter details interactively. If a default configuration with an API key already exists, the interactive setup will only ask for the location, skipping the API key prompt.

Additionally, the configuration file now supports an optional `logger` key (boolean, defaults to `false`). If set to `true`, the application will output logs to syslog.

## Sample `config.json`

```json
{
  "apiKey": "XXXXXX",
  "location": "auto:ip",
  "logger": false,
  "output": "table",
  "short_template": "{{.Emoji}} {{printf \"%.1f\" .TempC}}°",
  "current_template": "{{.Location}} - {{.Country}}",
  "forecast_template": "{{.Emoji}} {{printf \"%5.1f\" .TempC}}° [{{printf \"%5.1f\" .FeelslikeC}}°]",
  "forecastHours": 23,
  "noCache": false
}
```

## Configuration Entries

*   `apiKey`: Your weatherapi.com API key.
*   `location`: The default location to get the weather for. Can be a city name, a zip code, or `auto:ip` to use the IP address of the machine.
*   `logger`: If set to `true`, the application will output logs to syslog.
*   `output`: The default output format. Can be `table` or `json`.
*   `short_template`: The Go template used to format the `text` field when `output` is set to `json`.
*   `current_template`: The Go template for the location.
*   `forecast_template`: The Go template for the hourly forecast.
*   `forecastHours`: The number of forecast hours to display.
*   `noCache`: If set to `true`, the application will not use the cache.
