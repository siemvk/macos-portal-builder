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

func BenchmarkFindSteamLibraries(b *testing.B) {
	tempDir := b.TempDir()

	b.Setenv("HOME", tempDir)

	steamPath := filepath.Join(tempDir, "Library", "Application Support", "Steam", "steamapps")
	err := os.MkdirAll(steamPath, 0755)
	if err != nil {
		b.Fatal(err)
	}

	vdfPath := filepath.Join(steamPath, "libraryfolders.vdf")
	content := []byte(`"libraryfolders" {`)
	for i := 0; i < 100; i++ {
		content = append(content, []byte("\n\"path\" \"/path/to/steam/loop\"")...)
	}
	content = append(content, []byte(`}`)...)

	err = os.WriteFile(vdfPath, content, 0644)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findSteamLibraries()
	}
}
