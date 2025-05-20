package main

import (
	"echo/internals"
	"fmt"
	"net"
	"os"
	"path/filepath"
)

func Receive(conn *net.UDPConn) error {
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

	for {
		msg, _, err := internals.ReceivePacket(conn, buffer)
		if err != nil {
			return err
		}

		if msg.Version != VERSION {
			return fmt.Errorf("protocol version mismatch :(")
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
			count++
		}

		if count == expectedChunks {
			break
		}
	}

	for i := 1; i <= expectedChunks; i++ {
		data, exists := chunks[i]
		if !exists {
			return fmt.Errorf("missing chunk at index %d", i)
		}
		if _, err := outputFile.Write(data); err != nil {
			return fmt.Errorf("failed to write chunk %d: %v", i, err)
		}
	}

	return nil
}
