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
			name:     "already normalized",
			input:    "portal",
			expected: "portal",
		},
		{
			name:     "uppercase to lowercase",
			input:    "PORTAL",
			expected: "portal",
		},
		{
			name:     "mixed case",
			input:    "pOrTaL",
			expected: "portal",
		},
		{
			name:     "leading and trailing spaces",
			input:    "  portal  ",
			expected: "portal",
		},
		{
			name:     "spaces and uppercase",
			input:    "  HL2  ",
			expected: "hl2",
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
