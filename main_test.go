package main

import (
	"testing"
)

func TestValidateGameName(t *testing.T) {
	tests := []struct {
		name     string
		gameName string
		want     bool
	}{
		{"Valid lowercase portal", "portal", true},
		{"Valid lowercase hl2", "hl2", true},
		{"Valid uppercase portal", "PORTAL", true},
		{"Valid mixed case portal", "pOrTaL", true},
		{"Valid portal with leading/trailing spaces", "  portal  ", true},
		{"Valid hl2 with leading/trailing spaces", " hl2 ", true},
		{"Invalid game name tf2", "tf2", false},
		{"Invalid game name csgo", "csgo", false},
		{"Empty string", "", false},
		{"String with only spaces", "   ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateGameName(tt.gameName); got != tt.want {
				t.Errorf("validateGameName(%q) = %v, want %v", tt.gameName, got, tt.want)
			}
		})
	}
}
