package internals

import (
	"echo/fileproto"
	"echo/utils"
	"fmt"
	"net"
	"os"
	"time"

	"google.golang.org/protobuf/proto"
)

const maxRetries = 5

type Chunk struct {
	Index int
	Data  []byte
}

func SendPacket(conn *net.UDPConn, raddr *net.UDPAddr, chunk *Chunk, totalChunks uint32, file *os.File, version uint32, ackManager *AckManager) error {
	checksum := utils.CalculateChecksum(chunk.Data)

	msg := &fileproto.FileChunk{
		Version:     version,
		Filename:    file.Name(),
		ChunkIndex:  uint32(chunk.Index),
		TotalChunks: totalChunks,
		Data:        chunk.Data,
		Checksum:    checksum,
	}

	encoded, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	ch := ackManager.Register(uint32(chunk.Index))
	for retries := 0; retries < maxRetries; retries++ {
		_, err := conn.WriteToUDP(encoded, raddr)
		if err != nil {
			return err
		}

		select {
		case <-ch:
			return nil
		case <-time.After(2 * time.Second):
			fmt.Printf("Retry chunk %d (attempt %d)\n", chunk.Index, retries+2)
		}
	}

	return nil
}

func ReceivePacket(conn *net.UDPConn, buffer []byte) (*fileproto.FileChunk, *net.UDPAddr, error) {
	num, client, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, nil, err
	}

	var msg fileproto.FileChunk
	err = proto.Unmarshal(buffer[:num], &msg)
	if err != nil {
		return nil, nil, err
	}

	checksum := utils.CalculateChecksum(msg.Data)
	if checksum != msg.Checksum {
		return nil, nil, fmt.Errorf("checksum mismatch on chunk %d: expected %s, got %s", msg.ChunkIndex, msg.Checksum, checksum)
	}

	ack := &fileproto.FileAck{
		ChunkIndex: msg.ChunkIndex,
	}
	encodedAck, err := proto.Marshal(ack)
	if err != nil {
		return nil, nil, err
	}

	_, err = conn.WriteToUDP(encodedAck, client)
	if err != nil {
		return nil, nil, err
	}

	return &msg, client, nil
}
