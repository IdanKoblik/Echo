package main

import (
	"echo/internals"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

func Send(filename string, conn *net.UDPConn, remoteAddr string) error {
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
	const workerCount = 100
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
		go fixedWorker(assignedChunks, conn, raddr, uint32(chunkCount), file, &wg)
	}

	wg.Wait()
	return nil
}

func fixedWorker(chunks []internals.Chunk, conn *net.UDPConn, raddr *net.UDPAddr, total uint32, file *os.File, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, chunk := range chunks {
		err := internals.SendPacket(conn, raddr, &chunk, total, file)
		if err != nil {
			fmt.Println("Cannot send packet: ", err)
			return
		}
	}
}
