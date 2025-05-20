package internals

import (
	"echo/fileproto"
	"fmt"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
)

func HandleAck(connection net.Conn, expectedIndex uint32) (bool, error) {
	const timeout = 5 * time.Second // 5 seconds

	ackBuffer := make([]byte, 128)
	err := connection.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		fmt.Println("test1: ", err)
		return false, err
	}

	num, err := connection.Read(ackBuffer)
	if err != nil {
		fmt.Println("test2: ", err)
		return false, err
	}

	var ackMsg fileproto.FileAck
	err = proto.Unmarshal(ackBuffer[:num], &ackMsg)
	if err != nil {
		fmt.Println("test3: ", err)
		return false, err
	}

	return ackMsg.ChunkIndex == expectedIndex, nil
}
