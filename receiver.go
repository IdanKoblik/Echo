package main

import (
	"datadrop/fileproto"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"os"
	"path/filepath"
)

func Receive(address string) error {
	addr, err := net.ResolveUDPAddr("udp", address); if err != nil {
		return err
	}

	connection, err := net.ListenUDP("udp", addr); if err != nil {
		return err
	}

	defer connection.Close()

	var outputFile *os.File
	var fileName string

	homeDir, err := os.UserHomeDir(); if err != nil {
		return err
	}

	buffer := make([]byte, 2048)
	for {
		num, client, err := connection.ReadFromUDP(buffer); if err != nil {
			return err
		}

		var msg fileproto.FileChunk
		err = proto.Unmarshal(buffer[:num], &msg)
		if err != nil {
			return err
		}

		if outputFile == nil {
			fileName = msg.Filename
			filePath := filepath.Join(homeDir, filepath.Base(fileName))
			outputFile, err = os.Create(filePath); if err != nil {
				return fmt.Errorf("failed to create file: %v", err)
			}

			defer outputFile.Close()
			fmt.Printf("Creating file: %s\n", fileName)
		}

		_, err = outputFile.Write(msg.Data); if err != nil {
			return fmt.Errorf("failed to write chunk to file: %v", err)
		}

		fmt.Printf("Received chunk %d/%d: %s\n", msg.ChunkIndex+1, msg.TotalChunks, msg.Filename)
		ack := &fileproto.FileAck{
			ChunkIndex: msg.ChunkIndex,
		}

		encodedAck, err := proto.Marshal(ack); if err != nil {
			return err
		}

		_, err = connection.WriteToUDP(encodedAck, client); if err != nil {
			return err
		}

		if msg.IsLastChunk {
			checksum, err := GetFileChecksum(outputFile); if err != nil {
				return err
			}

			if checksum != msg.Checksum {
				return fmt.Errorf("invalid checksums")
			}

			fmt.Println("Received last chunk. File transfer complete!")
			break
		}

		log.Printf("Sent ACK for chunk %d\n", msg.ChunkIndex)
	}

	return nil
}