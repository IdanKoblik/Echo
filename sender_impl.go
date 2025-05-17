package main

import (
	"echo/internals"
	"echo/ui"
	"io"
	"net"
	"os"
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
	var chunks [][]byte

	for {
		num, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		chunk := make([]byte, num)
		copy(chunk, buffer[:num])
		chunks = append(chunks, chunk)
	}

	total := uint32(len(chunks))
	progress := ui.ProgressBar {
		Len: int(total),
		Description: "Sending file",
	}

	bar := progress.Init()
	for i, chunk := range chunks {
		isLastChunk := (i == len(chunks)-1)
		err := internals.SendPacket(conn, raddr, filename, chunk, uint32(i), total, isLastChunk, file)
		if err != nil {
			return err
		}

		bar.Add(1)
	}

	return nil
}