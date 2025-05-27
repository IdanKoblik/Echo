package main

import (
	"echo/internals"
	"echo/utils"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

func Send(conn *net.UDPConn, cfg *utils.Config) error {
	start := time.Now()

	file, err := os.Open(cfg.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	raddr, err := net.ResolveUDPAddr("udp", cfg.RemoteAddr)
	if err != nil {
		return err
	}

	const chunkSize = 1400
	buffer := make([]byte, chunkSize)
	var chunks []internals.Chunk
	index := 0

	stats := &BenchmarkStats{}

	for {
		num, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		data := make([]byte, num)
		copy(data, buffer[:num])
		index++

		chunks = append(chunks, internals.Chunk{
			Data:  data,
			Index: index,
		})
	}

	chunkCount := len(chunks)
	uploadMbps, rttMs, err := MeasureUpload(); if err != nil {
		return err
	}

	chunksPerSecond := (uploadMbps * 125000) / float64(chunkSize)
	acksPerWorker := 1000.0 / rttMs
	optimalWorkersFloat := chunksPerSecond / acksPerWorker
	optimalWorkers := int(optimalWorkersFloat)
	if optimalWorkers < 1 {
		optimalWorkers = 1
	}
	
	chunksPerWorker := (chunkCount + optimalWorkers - 1) / optimalWorkers
	ackManager := internals.NewAckManager()
	go ackManager.Listen(conn)

	
	fmt.Println("Number of workers: ", optimalWorkers)
	var wg sync.WaitGroup
	for w := 0; w < optimalWorkers; w++ {
		start := w * chunksPerWorker
		end := (w + 1) * chunksPerWorker
		if end > chunkCount {
			end = chunkCount
		}

		if start >= chunkCount {
			break
		}

		assignedChunks := chunks[start:end]
		wg.Add(1)
		go func(workerChunks []internals.Chunk) {
			defer wg.Done()
			var totalBytesTransferred int64
			var totalPacketsSent int

			for _, chunk := range workerChunks {
				start := time.Now()
				err := internals.SendPacket(conn, raddr, &chunk, uint32(chunkCount), file, VERSION, ackManager)
				if err != nil {
					fmt.Println("Send error:", err)
					return
				}

				stats.ChunkTimings = append(stats.ChunkTimings, time.Since(start))
				totalBytesTransferred += int64(len(chunk.Data))
				totalPacketsSent++
			}

			stats.TotalBytes += totalBytesTransferred
			stats.PacketsSent += totalPacketsSent
		}(assignedChunks)
	}

	wg.Wait()

	duration := time.Since(start)
	stats.TotalTime = duration
	stats.TransferSpeed = float64(stats.TotalBytes) / duration.Seconds()
	stats.MemoryUsage = GetMemoryUsage()
	stats.CpuUsage = getCpuUsage()
	stats.PrintStats(cfg.BenchMark)

	return nil
}
