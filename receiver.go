package main

import (
	"echo/fileproto"
	"echo/utils"
	"fmt"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/protobuf/proto"
	"net"
	"os"
	"path/filepath"
	"time"
)

func Receive(conn *net.UDPConn) error {
	startTime := time.Now()
	var outputFile *os.File
	var fileName string
	var totalChunks uint32
	var receivedSize int64

	success := color.New(color.FgGreen).SprintFunc()
	info := color.New(color.FgCyan).SprintFunc()
	warn := color.New(color.FgYellow).SprintFunc()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	var bar *progressbar.ProgressBar
	buffer := make([]byte, 2048)

	fmt.Printf("%s Waiting for incoming file...\n", info("ℹ"))

	for {
		num, client, err := conn.ReadFromUDP(buffer)
		if err != nil {
			return err
		}

		var msg fileproto.FileChunk
		err = proto.Unmarshal(buffer[:num], &msg)
		if err != nil {
			return err
		}

		if outputFile == nil {
			fileName = msg.Filename
			totalChunks = msg.TotalChunks
			filePath := filepath.Join(homeDir, filepath.Base(fileName))
			outputFile, err = os.Create(filePath)
			if err != nil {
				return fmt.Errorf("failed to create file: %v", err)
			}
			defer outputFile.Close()

			fmt.Printf("%s Receiving: %s\n", info("ℹ"), fileName)
			bar = progressbar.NewOptions(int(totalChunks),
				progressbar.OptionSetDescription("Receiving file"),
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
		}

		written, err := outputFile.Write(msg.Data)
		if err != nil {
			return fmt.Errorf("failed to write chunk to file: %v", err)
		}
		receivedSize += int64(written)

		ack := &fileproto.FileAck{
			ChunkIndex: msg.ChunkIndex,
		}

		encodedAck, err := proto.Marshal(ack)
		if err != nil {
			return err
		}

		_, err = conn.WriteToUDP(encodedAck, client)
		if err != nil {
			return err
		}

		bar.Add(1)

		if msg.IsLastChunk {
			checksum, err := utils.GetFileChecksum(outputFile)
			if err != nil {
				return err
			}

			if checksum != msg.Checksum {
				fmt.Printf("\n%s Checksum verification failed!\n", warn("⚠"))
				return fmt.Errorf("invalid checksums")
			}

			duration := time.Since(startTime)
			speed := float64(receivedSize) / duration.Seconds() / (1024 * 1024) // MB/s

			fmt.Printf("\n%s Transfer complete!\n", success("✓"))
			fmt.Printf("%s Saved to: %s\n", info("ℹ"), filepath.Join(homeDir, filepath.Base(fileName)))
			fmt.Printf("%s File size: %.2f MB\n", info("ℹ"), float64(receivedSize)/(1024*1024))
			fmt.Printf("%s Time taken: %s\n", info("ℹ"), duration.Round(time.Second))
			fmt.Printf("%s Average speed: %.2f MB/s\n", info("ℹ"), speed)
			break
		}
	}

	return nil
}
