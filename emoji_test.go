package main

import (
	"testing"
)

func TestGetEmoji(t *testing.T) {
	// Manually set emojiMap for testing GetEmoji
	emojiMap = map[int]string{
		1000: "󰖙",
		1003: "☁️",
	}

	// Test known code
	if GetEmoji(1000) != "󰖙" {
		t.Errorf("Expected emoji for 1000 to be '󰖙', got: %s", GetEmoji(1000))
	}

	// Test unknown code
	if GetEmoji(9999) != "❓" {
		t.Errorf("Expected emoji for 9999 to be '❓', got: %s", GetEmoji(9999))
	}
}
