# Usage

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

To force a refresh of the data from the API, use the `-f` or `--no-cache` flag:
```bash
./wayther -f
```

To clean the cache, use the `-C` or `--clean-cache` flag. This will clean entries that are older than one hour:
```bash
./wayther -C
```

To specify the number of forecast hours to display, use the `-n` or `--forecast-hours` flag. 0 means no forecast and the max is 23 hours of forecast:
```bash
./wayther -n 5
```

By default, if you are in a terminal, the output will be a human-readable table:

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

If you are not in a terminal (e.g., piping the output to another command or using it with waybar), the output will be in JSON format.

To force JSON output, use the `--output json` flag:

```bash
./wayther "London" --output json
```

The JSON output will be a JSON object with `text` (current weather summary) and `tooltip` (hourly forecast) fields:

For details on customizing the output using templates, see [Templates](templates.md).

```json
{"text":"🌦️  28.0°","tooltip":"20:00: ☀️ 23.6° [ 25.1°]\r21:00: ☀️ 22.8° [ 24.8°]\r22:00: ☀️ 22.1° [ 24.6°]\r23:00: ☀️ 21.5° [ 21.5°]\r00:00: ☀️ 20.8° [ 20.8°]\r01:00: ☀️ 20.2° [ 20.2°]\r02:00: ☀️ 19.7° [ 19.7°]\r03:00: ☀️ 19.3° [ 19.3°]\r04:00: ☀️ 18.9° [ 18.9°]\r05:00: ☀️ 19.1° [ 19.2°]\r06:00:  22.0° [ 22.0°]\r07:00: ☀️ 24.7° [ 25.3°]\r08:00: ☀️ 27.5° [ 27.0°]\r09:00: ☀️ 29.8° [ 28.9°]\r10:00: ☀️ 31.8° [ 30.7°]\r11:00: ☀️ 33.4° [ 32.2°]\r12:00: ☀️ 34.4° [ 33.2°]\r13:00: ☀️ 35.1° [ 33.8°]\r14:00: ☀️ 35.5° [ 34.1°]\r15:00: ☀️ 35.5° [ 33.9°]\r16:00: ☀️ 35.3° [ 33.8°]\r17:00: ☀️ 34.1° [ 33.3°]\r18:00: ☀️ 30.8° [ 30.2°]\r19:00: ☀️ 27.8° [ 27.4°]"}
```

### waybar Integration

For waybar to find the `wayther` executable, ensure it's placed in your system's PATH (e.g., `/usr/local/bin`). We plan to support package managers (RPM, Deb, ebuilds, etc.) in the future for easier installation.

After the initial run in a terminal (to set up the API key and default location), you can integrate Wayther into your waybar configuration. Wayther outputs JSON when not in a terminal, which is ideal for waybar's `custom` module.

Add the following to your waybar `config` file (e.g., `~/.config/waybar/config.jsonc`):

```jsonc
"custom/wayther": {
    "exec": "wayther",
    "return-type": "json",
    "format": "{} ",
    "on-click": "wayther",
    "interval": 3600, // once every day [this is an example]
    "tooltip": true,
},
```

```