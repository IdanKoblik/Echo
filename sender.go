package main

import (
	"echo/fileproto"
	"echo/utils"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const VERSION = 1

func Send(filename string, conn *net.UDPConn, remoteAddr string) error {
	file, err := os.Open(filename); if err != nil {
		return err
	}

	defer file.Close()

	raddr, err := net.ResolveUDPAddr("udp", remoteAddr); if err != nil {
		return err
	}

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
		msg := &fileproto.FileChunk{
			Version:     uint32(VERSION),
			Filename:    filename,
			ChunkIndex:  uint32(i),
			TotalChunks: total,
			Data:        chunk,
			IsLastChunk: (i == len(chunks)-1),
		}

		if msg.IsLastChunk {
			checksum, err := utils.GetFileChecksum(file); if err != nil {
				return err
			}

			msg.Checksum = checksum
		}

		encoded, err := proto.Marshal(msg); if err != nil {
			return err
		}

		_, err = conn.WriteToUDP(encoded, raddr); if err != nil {
			return err
		}

		ok, err := handleAck(conn, uint32(i))
		if err != nil || !ok {
			log.Printf("ACK failed for chunk %d: %v. Retrying...", i, err)
			i--
			continue
		}

		log.Printf("Received ACK for chunk %d/%d", i+1, total)
	}

	log.Println("File transfer complete!")
	return nil
}

func handleAck(connection net.Conn, expectedIndex uint32) (bool, error) {
	const timeout = 5 * time.Second // 5 seconds

	ackBuffer := make([]byte, 128)
	err := connection.SetReadDeadline(time.Now().Add(timeout)); if err != nil {
		return false, err
	}

	num, err := connection.Read(ackBuffer); if err != nil {
		return false, err
	}

	var ackMsg fileproto.FileAck
	err = proto.Unmarshal(ackBuffer[:num], &ackMsg); if err != nil {
		return false, err
	}

	return ackMsg.ChunkIndex == expectedIndex, nil
}
