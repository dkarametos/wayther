[![Go CI](https://github.com/dkarametos/wayther/actions/workflows/go.yml/badge.svg)](https://github.com/dkarametos/wayther/actions/workflows/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# Wayther

**Wayther** is a command-line application to get the current weather for a given location, providing output in JSON or a human-readable table format. 

The name "Wayther" is a portmanteau of "[waybar](https://github.com/Alexays/waybar)" and "weather", as it was specifically designed to be used with waybar.

## Features

*   **Multiple Output Formats**: Output the weather in JSON or a human-readable table.
*   **Customizable Templates**: Customize the output format using Go templates.
*   **Configuration Merging**: Merge multiple configuration files.
*   **Caching**: Cache weather data to reduce API calls.
*   **Syslog**: Log to syslog.
*   **Interactive Setup**: Interactive setup for the first run.

## Documentation

*   [Installation](docs/installation.md)
*   [Usage](docs/usage.md)
*   [Configuration](docs/configuration.md)
*   [Testing](docs/testing.md)
*   [License](docs/license.md)
