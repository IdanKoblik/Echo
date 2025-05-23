package main

import (
	"echo/internals"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"
)

func Receive(conn *net.UDPConn, benchmark bool) error {
	var outputFile *os.File
	var fileName string

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	chunks := make(map[int][]byte)
	buffer := make([]byte, 2048)
	count := 0
	var expectedChunks int

	stats := &BenchmarkStats{}
	var start time.Time
	flag := false

	for {
		msg, _, err := internals.ReceivePacket(conn, buffer)
		if err != nil {
			return err
		}

		if msg.Version != VERSION {
			return fmt.Errorf("protocol version mismatch :(")
		}

		if !flag {
			flag = true
			start = time.Now()
		}

		if outputFile == nil {
			fileName = msg.Filename
			filePath := filepath.Join(homeDir, filepath.Base(fileName))
			outputFile, err = os.Create(filePath)
			if err != nil {
				return fmt.Errorf("failed to create file: %v", err)
			}
			defer outputFile.Close()

			expectedChunks = int(msg.TotalChunks)
		}

		if _, exists := chunks[int(msg.ChunkIndex)]; !exists {
			chunks[int(msg.ChunkIndex)] = msg.Data
			stats.ChunkTimings = append(stats.ChunkTimings, time.Since(start))
			count++
			stats.PacketsReceived++
		}

		if count == expectedChunks {
			break
		}
	}

	flag = false
	duration := time.Since(start)

	for i := 1; i <= expectedChunks; i++ {
		data, exists := chunks[i]
		if !exists {
			return fmt.Errorf("missing chunk at index %d", i)
		}
		if _, err := outputFile.Write(data); err != nil {
			return fmt.Errorf("failed to write chunk %d: %v", i, err)
		}

		stats.TotalBytes += int64(len(data))
	}

	stats.TotalTime = duration
	stats.PacketLoss = (1 - float64(stats.PacketsReceived)/float64(expectedChunks)) * 100
	stats.CpuUsage = getCpuUsage()
	stats.MemoryUsage = GetMemoryUsage()

	stats.PrintStats(benchmark)

	return nil
}
