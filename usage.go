package main

import "runtime"

func GetMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Sys
}

func getCpuUsage() float64 {
	return float64(runtime.NumCPU()) * 100
}
