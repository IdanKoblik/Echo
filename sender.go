package main

import (
	"echo/fileproto"
	"echo/utils"
	"io"
	"net"
	"os"
	"time"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/protobuf/proto"
	"fmt"
	"github.com/fatih/color"
)

const VERSION = 1

func Send(filename string, conn *net.UDPConn, remoteAddr string) error {
	startTime := time.Now()
	success := color.New(color.FgGreen).SprintFunc()
	info := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("%s Opening file: %s\n", info("ℹ"), filename)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	totalSize := fileInfo.Size()

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
	fmt.Printf("%s Total chunks: %d (%.2f MB)\n", info("ℹ"), total, float64(totalSize)/(1024*1024))

	bar := progressbar.NewOptions(len(chunks),
		progressbar.OptionSetDescription("Sending file"),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]━[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: "━",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionOnCompletion(func() { fmt.Println() }),
	)

	retries := 0
	for i, chunk := range chunks {
		msg := &fileproto.FileChunk{
			Version:     uint32(VERSION),
			Filename:    filename,
			ChunkIndex:  uint32(i),
			TotalChunks: total,
			Data:        chunk,
			IsLastChunk: (i == len(chunks)-1),
		}

		if msg.IsLastChunk {
			checksum, err := utils.GetFileChecksum(file)
			if err != nil {
				return err
			}
			msg.Checksum = checksum
		}

		encoded, err := proto.Marshal(msg)
		if err != nil {
			return err
		}

		_, err = conn.WriteToUDP(encoded, raddr)
		if err != nil {
			return err
		}

		ok, err := handleAck(conn, uint32(i))
		if err != nil || !ok {
			if retries < 3 {
				retries++
				i--
				continue
			}
			return fmt.Errorf("failed to receive ACK after 3 retries for chunk %d: %v", i, err)
		}
		retries = 0

		bar.Add(1)
	}

	duration := time.Since(startTime)
	speed := float64(totalSize) / duration.Seconds() / (1024 * 1024) // MB/s

	fmt.Printf("\n%s Transfer complete!\n", success("✓"))
	fmt.Printf("%s Time taken: %s\n", info("ℹ"), duration.Round(time.Second))
	fmt.Printf("%s Average speed: %.2f MB/s\n", info("ℹ"), speed)

	return nil
}

func handleAck(connection net.Conn, expectedIndex uint32) (bool, error) {
	const timeout = 5 * time.Second // 5 seconds

	ackBuffer := make([]byte, 128)
	err := connection.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return false, err
	}

	num, err := connection.Read(ackBuffer)
	if err != nil {
		return false, err
	}

	var ackMsg fileproto.FileAck
	err = proto.Unmarshal(ackBuffer[:num], &ackMsg)
	if err != nil {
		return false, err
	}

	return ackMsg.ChunkIndex == expectedIndex, nil
}
