package main

import (
	"testing"
)

func TestGetEmojiForWeatherCode(t *testing.T) {
	// Manually set weatherCodeToEmojiMap for testing getEmojiForWeatherCode
	weatherCodeToEmojiMap = map[int]string{
		1000: "󰖙",
		1003: "☁️",
	}

	// Test known code
	if getEmojiForWeatherCode(1000) != "󰖙" {
		t.Errorf("Expected emoji for 1000 to be '󰖙', got: %s", getEmojiForWeatherCode(1000))
	}

	// Test unknown code
	if getEmojiForWeatherCode(9999) != "❓" {
		t.Errorf("Expected emoji for 9999 to be '❓', got: %s", getEmojiForWeatherCode(9999))
	}
}
