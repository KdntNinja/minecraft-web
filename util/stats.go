package util

import "runtime"

// GetGCPercent returns the current GC percent (heap in use / next GC)
func GetGCPercent() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	if m.NextGC == 0 {
		return 0
	}
	return float64(m.HeapAlloc) / float64(m.NextGC) * 100.0
}

// GetTickTimes returns a slice of recent tick times in ms (dummy, replace with real times)
func GetTickTimes(n int) []float64 {
	// TODO: Replace with real tick times from your game loop
	times := make([]float64, n)
	for i := range times {
		times[i] = 16 + 4*float64(i%10)/10 // Simulated ms
	}
	return times
}
