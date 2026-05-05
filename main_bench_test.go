package main

import "testing"

func BenchmarkFindSteamLibraries(b *testing.B) {
	for i := 0; i < b.N; i++ {
		findSteamLibraries()
	}
}
