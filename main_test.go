package main

import (
	"testing"
)

func TestNormalizeGameName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase",
			input:    "portal",
			expected: "portal",
		},
		{
			name:     "uppercase",
			input:    "PORTAL",
			expected: "portal",
		},
		{
			name:     "mixed case",
			input:    "PorTal",
			expected: "portal",
		},
		{
			name:     "with leading spaces",
			input:    "  hl2",
			expected: "hl2",
		},
		{
			name:     "with trailing spaces",
			input:    "hl2  ",
			expected: "hl2",
		},
		{
			name:     "with both spaces",
			input:    "  hl2  ",
			expected: "hl2",
		},
		{
			name:     "with spaces and mixed case",
			input:    "  PoRTal  ",
			expected: "portal",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeGameName(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeGameName(%q) = %q; expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateGameName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid portal",
			input:    "portal",
			expected: true,
		},
		{
			name:     "valid hl2",
			input:    "hl2",
			expected: true,
		},
		{
			name:     "valid portal uppercase",
			input:    "PORTAL",
			expected: true,
		},
		{
			name:     "valid hl2 uppercase",
			input:    "HL2",
			expected: true,
		},
		{
			name:     "valid hl2 with spaces",
			input:    "  hl2  ",
			expected: true,
		},
		{
			name:     "valid portal with tabs and newlines",
			input:    "\tportal\n",
			expected: true,
		},
		{
			name:     "invalid game",
			input:    "tf2",
			expected: false,
		},
		{
			name:     "invalid csgo",
			input:    "csgo",
			expected: false,
		},
		{
			name:     "invalid empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "invalid portal2",
			input:    "portal2",
			expected: false,
		},
		{
			name:     "invalid almost correct",
			input:    "halflife2",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateGameName(tt.input)
			if result != tt.expected {
				t.Errorf("validateGameName(%q) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
