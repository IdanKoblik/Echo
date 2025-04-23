package main

import (
	"echo/fileproto"
	"io"
	"log"
	"net"
	"os"
	"time"
	"google.golang.org/protobuf/proto"
)

const VERSION = 1

func Send(filename, address string) error {
	file, err := os.Open(filename); if err != nil {
		return err
	}

	defer file.Close()
	connection, err := net.Dial("udp", address); if err != nil {
		return err
	}

	defer connection.Close()

	const chunkSize = 1024
	buffer := make([]byte, chunkSize)
	var chunks [][]byte

	for {
		num, err := file.Read(buffer); if err == io.EOF {
			break
		}

		chunk := make([]byte, num)
		copy(chunk, buffer[:num])
		chunks = append(chunks, chunk)
	}

	total := uint32(len(chunks))
	for i, chunk := range chunks {
		msg := &fileproto.FileChunk {
			Version: uint32(VERSION),
			Filename: filename,
			ChunkIndex: uint32(i),
			TotalChunks: total,
			Data: chunk,	
			IsLastChunk: (i == len(chunks)-1),
		}

		if msg.IsLastChunk {
			checksum, err := GetFileChecksum(file); if err != nil {
				return err
			}

			file.Seek(0, 0)
			msg.Checksum = checksum 
		}

		encoded, err := proto.Marshal(msg); if err != nil {
			return err
		}

		_, err = connection.Write(encoded); if err != nil {
			return err
		}

		ok, err := handleAck(connection, uint32(i)); if err != nil || !ok {
			log.Printf("ACK failed for chunk %d: %v. Retrying...", i, err)
			i--
			continue
		}

		log.Printf("Received ACK for chunk %d/%d", i+1, total)
	}

	log.Println("File transfer complete!")
	return nil;
}

func handleAck(connection net.Conn, expectedIndex uint32) (bool, error) {
	const timeout = 5 * time.Second // 5 seconds

	ackBuffer := make([]byte, 128)
	connection.SetReadDeadline(time.Now().Add(timeout))

	num, err := connection.Read(ackBuffer); if err != nil {
		return false, err
	}

	var ackMsg fileproto.FileAck
	err = proto.Unmarshal(ackBuffer[:num], &ackMsg); if err != nil {
		return false, err
	}

	return ackMsg.ChunkIndex == expectedIndex, nil
}