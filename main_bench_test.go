package main

import "testing"

func BenchmarkFindSteamLibraries(b *testing.B) {
	findSteamLibraries() // Trigger cache initialization
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findSteamLibraries()
	}
}
