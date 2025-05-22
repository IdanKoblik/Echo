package main

import (
	"echo/internals"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

func Send(filename string, conn *net.UDPConn, remoteAddr string, workerCount int, benchmark bool) error {
	start := time.Now()
	
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	raddr, err := net.ResolveUDPAddr("udp", remoteAddr)
	if err != nil {
		return err
	}

	const chunkSize = 1024
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
	chunksPerWorker := (chunkCount + workerCount - 1) / workerCount

	var wg sync.WaitGroup
	for w := 0; w < workerCount; w++ {
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
		go fixedWorker(assignedChunks, conn, raddr, uint32(chunkCount), file, &wg, stats)
	}

	wg.Wait()

	duration := time.Since(start)
	stats.TotalTime = duration
	stats.TransferSpeed = float64(stats.TotalBytes) / duration.Seconds()

	stats.MemoryUsage = GetMemoryUsage()
	stats.CpuUsage = getCpuUsage()

	stats.PrintStats(benchmark)

	return nil
}

func fixedWorker(chunks []internals.Chunk, conn *net.UDPConn, raddr *net.UDPAddr, total uint32, file *os.File, wg *sync.WaitGroup, stats *BenchmarkStats) {
	defer wg.Done()

	var totalBytesTransferred int64
	var totalPacketsSent int

	for _, chunk := range chunks {
		start := time.Now()
		err := internals.SendPacket(conn, raddr, &chunk, total, file, VERSION)
		if err != nil {
			fmt.Println("Cannot send packet: ", err)
			return
		}

		stats.ChunkTimings = append(stats.ChunkTimings, time.Since(start))
		totalBytesTransferred += int64(len(chunk.Data))
		totalPacketsSent++
	}

	stats.TotalBytes += totalBytesTransferred
	stats.PacketsSent += totalPacketsSent
}
