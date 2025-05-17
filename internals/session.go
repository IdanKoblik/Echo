package internals

import (
	"echo/fileproto"
	"echo/utils"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net"
	"os"
)

const maxRetries = 3

func SendPacket(conn *net.UDPConn, raddr *net.UDPAddr, filename string, chunk []byte, chunkIndex uint32, totalChunks uint32, isLastChunk bool, file *os.File) error {
	msg := &fileproto.FileChunk{
		Version:     uint32(1),
		Filename:    filename,
		ChunkIndex:  chunkIndex,
		TotalChunks: totalChunks,
		Data:        chunk,
		IsLastChunk: isLastChunk,
	}

	if isLastChunk {
		checksum, err := utils.GetFileChecksum(file)
		if err != nil {
			return err
		}
		
		if msg.Checksum != checksum {
			return fmt.Errorf("checksum mismatch")
		}
	}

	encoded, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	retries := 0
	_, err = conn.WriteToUDP(encoded, raddr)
	if err != nil {
		return err
	}

	ok, err := HandleAck(conn, chunkIndex)
	if err != nil || !ok {
		if retries < maxRetries {
			retries++
			return SendPacket(conn, raddr, filename, chunk, chunkIndex, totalChunks, isLastChunk, file)
		}

		return fmt.Errorf("failed to receive ACK after %d retries for chunk %d: %v", maxRetries, chunkIndex, err)
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