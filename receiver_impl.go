package main

import (
	"echo/internals"
	"echo/ui"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

func Receive(conn *net.UDPConn) error {
	var outputFile *os.File
	var fileName string

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	buffer := make([]byte, 2048)

	var progressBar *progressbar.ProgressBar
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

			progress := ui.ProgressBar {
				Len: int(msg.TotalChunks),
				Description: "Receiving file",
			}

			progressBar = progress.Init()
		}

		_, err = outputFile.Write(msg.Data)
		if err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}

		progressBar.Add(1)
		if msg.IsLastChunk {
			return nil
		}
	}

}