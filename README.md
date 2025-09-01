# Wayther

## Why wayther?
The name "Wayther" is a portmanteau of "Waybar" and "Weather". This application was specifically designed with the intention of being used to display weather information on a Waybar, a highly customizable status bar for Wayland compositors.

A command-line application to get the current weather for a given location, providing output in JSON format.

## Configuration

This application uses the WeatherAPI.com for weather data. You will need to obtain a free API key from their website.

The application now supports robust configuration merging:
*   **Default Configuration:** A default configuration file is located at `XDG_CONFIG_HOME/wayther/config.json`.
*   **Custom Configurations:** You can specify a custom configuration file using the `-c` or `--config` flag.
*   **Merging Logic:** When a custom configuration is used, its values will override those in the default configuration. If a key is missing in the custom configuration, the value from the default configuration will be used. This means keys in custom configurations are optional if they are defined in the default.

The first time you run the application without an existing configuration, it will prompt you to enter details interactively. If a default configuration with an API key already exists, the interactive setup will only ask for the location, skipping the API key prompt.

Additionally, the configuration file now supports an optional `logger` key (boolean, defaults to `false`). If set to `true`, the application will output logs to syslog.

## Building

To build the executable:
```bash
go build -o wayther
```
This will create an executable named `wayther` in the current directory.

## Usage

To get weather for a default location (from `config.json`):
```bash
./wayther
```

To get weather for a specific location:
```bash
./wayther "London"
```

To display help message and usage information:
```bash
./wayther --help
# or
./wayther -h
```

By default, if you are in a terminal, the output will be a human-readable table:

```
┌───────────────────────────┐
│ Current:                  │
├───────────────────────────┤
│ 🌦️ 24.0°                  │
│ Sandweiler - Luxembourg   │
├───────────────────────────┤
│ Hourly Forecast:          │
├───────────────────────────┤
│ 00:00:  🌦️ 24.0° [ 20.8°] │
│ 01:00:  ☀️ 20.2° [ 20.2°] │
│ 02:00:  ☀️ 19.7° [ 19.7°] │
│ 03:00:  ☀️ 19.3° [ 19.3°] │
│ 04:00:  ☀️ 18.9° [ 18.9°] │
│ 05:00:  ☀️ 19.1° [ 19.2°] │
│ 06:00:  ☀️ 22.0° [ 22.0°] │
│ 07:00:  ☀️ 24.7° [ 25.3°] │
│ 08:00:  ☀️ 27.5° [ 27.0°] │
│ 09:00:  ☀️ 29.8° [ 28.9°] │
│ 10:00:  ☀️ 31.8° [ 30.7°] │
│ 11:00:  ☀️ 33.4° [ 32.2°] │
│ 12:00:  ☀️ 34.4° [ 33.2°] │
│ 13:00:  ☀️ 35.1° [ 33.8°] │
│ 14:00:  ☀️ 35.5° [ 34.1°] │
│ 15:00:  ☀️ 35.5° [ 33.9°] │
│ 16:00:  ☀️ 35.3° [ 33.8°] │
│ 17:00:  ☀️ 34.1° [ 33.3°] │
│ 18:00:  ☀️ 30.8° [ 30.2°] │
│ 19:00:  ☀️ 27.8° [ 27.4°] │
│ 20:00:  ☀️ 26.5° [ 26.4°] │
│ 21:00:  ☀️ 24.9° [ 25.5°] │
│ 22:00:  ☀️ 23.6° [ 24.9°] │
│ 23:00:  ☀️ 22.7° [ 24.7°] │
└───────────────────────────┘
```

If you are not in a terminal (e.g., piping the output to another command), the output will be in JSON format.

To force JSON output, use the `--json` flag:

```bash
./wayther "London" --json
```

The JSON output will be a JSON object with `text` (current weather summary) and `tooltip` (hourly forecast) fields:

```json
{"text":"🌦️  28.0°","tooltip":"20:00: ☀️ 23.6° [ 25.1°]\r21:00: ☀️ 22.8° [ 24.8°]\r22:00: ☀️ 22.1° [ 24.6°]\r23:00: ☀️ 21.5° [ 21.5°]\r00:00: ☀️ 20.8° [ 20.8°]\r01:00: ☀️ 20.2° [ 20.2°]\r02:00: ☀️ 19.7° [ 19.7°]\r03:00: ☀️ 19.3° [ 19.3°]\r04:00: ☀️ 18.9° [ 18.9°]\r05:00: ☀️ 19.1° [ 19.2°]\r06:00:  22.0° [ 22.0°]\r07:00: ☀️ 24.7° [ 25.3°]\r08:00: ☀️ 27.5° [ 27.0°]\r09:00: ☀️ 29.8° [ 28.9°]\r10:00: ☀️ 31.8° [ 30.7°]\r11:00: ☀️ 33.4° [ 32.2°]\r12:00: ☀️ 34.4° [ 33.2°]\r13:00: ☀️ 35.1° [ 33.8°]\r14:00: ☀️ 35.5° [ 34.1°]\r15:00: ☀️ 35.5° [ 33.9°]\r16:00: ☀️ 35.3° [ 33.8°]\r17:00: ☀️ 34.1° [ 33.3°]\r18:00: ☀️ 30.8° [ 30.2°]\r19:00: ☀️ 27.8° [ 27.4°]"}"
```
## Testing

To run the tests:
```bash
go test .
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

