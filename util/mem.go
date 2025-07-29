package util

import "runtime"

// GetMemoryUsageMB returns the current memory usage in MB
func GetMemoryUsageMB() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.Alloc) / 1024.0 / 1024.0
}
