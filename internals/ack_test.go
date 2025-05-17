package internals_test

import (
	"echo/fileproto"
	"echo/internals"
	"net"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
)

func TestHandleAck(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	expectedIndex := uint32(42)
	ackMsg := &fileproto.FileAck{
		ChunkIndex: expectedIndex,
	}

	data, err := proto.Marshal(ackMsg)
	if err != nil {
		t.Fatalf("failed to marshal protobuf: %v", err)
	}

	go func() {
		_, err := clientConn.Write(data)
		if err != nil {
			t.Errorf("failed to write to clientConn: %v", err)
		}
	}()

	ok, err := internals.HandleAck(serverConn, expectedIndex)
	if err != nil {
		t.Fatalf("HandleAck returned error: %v", err)
	}

	if !ok {
		t.Fatalf("expected ack for chunk %d, but got false", expectedIndex)
	}
}

func TestHandleAck_WrongIndex(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	expectedIndex := uint32(42)
	wrongIndex := uint32(43)

	ackMsg := &fileproto.FileAck{
		ChunkIndex: wrongIndex,
	}
	data, err := proto.Marshal(ackMsg)
	if err != nil {
		t.Fatalf("failed to marshal protobuf: %v", err)
	}

	go func() {
		_, err := clientConn.Write(data)
		if err != nil {
			t.Errorf("failed to write to clientConn: %v", err)
		}
	}()

	ok, err := internals.HandleAck(serverConn, expectedIndex)
	if err != nil {
		t.Fatalf("HandleAck returned error: %v", err)
	}
	if ok {
		t.Fatalf("expected false for wrong chunk index, but got true")
	}
}

func TestHandleAck_Timeout(t *testing.T) {
	serverConn, _ := net.Pipe()
	defer serverConn.Close()

	expectedIndex := uint32(42)

	start := time.Now()
	_, err := internals.HandleAck(serverConn, expectedIndex)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected timeout error but got nil")
	}

	if elapsed < 4*time.Second {
		t.Fatalf("expected timeout after ~5s, but returned early after %v", elapsed)
	}
}
