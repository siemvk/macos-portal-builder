package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindSteamLibraries(t *testing.T) {
	tempDir := t.TempDir()

	t.Setenv("HOME", tempDir)

	steamPath := filepath.Join(tempDir, "Library", "Application Support", "Steam", "steamapps")
	err := os.MkdirAll(steamPath, 0755)
	if err != nil {
		t.Fatal(err)
	}

	vdfPath := filepath.Join(steamPath, "libraryfolders.vdf")
	content := []byte(`"libraryfolders" {`)
	for i := 0; i < 100; i++ {
		content = append(content, []byte("\n\"path\" \"/path/to/steam/loop\"")...)
	}
	content = append(content, []byte(`}`)...)

	err = os.WriteFile(vdfPath, content, 0644)
	if err != nil {
		t.Fatal(err)
	}

	libraries := findSteamLibraries()
	if len(libraries) != 2 {
		t.Fatalf("Expected 2 libraries, got %d", len(libraries))
	}
}

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
			name:     "valid hl2 with spaces",
			input:    "  hl2  ",
			expected: true,
		},
		{
			name:     "invalid game",
			input:    "tf2",
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
