package main

import (
	"bytes"
	"echo/utils"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	LOCAL_PORT  = "9000"
	REMOTE_ADDR = "0.0.0.0:9001"
)

func runWithArgs(args []string) {
	oldArgs := os.Args
	os.Args = append([]string{"cmd"}, args...)
	defer func() { os.Args = oldArgs }()

	main()
}

func outputHelper(args []string) string {
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	runWithArgs(args)

	w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestInvalidMode(t *testing.T) {
	output := outputHelper([]string{"--mode", "invalidMode"})
	assert.Contains(t, output, "Invalid input")
}

func TestMissingRemoteAddress(t *testing.T) {
	output := outputHelper([]string{"--mode", "send", "--file", "test.txt"})
	assert.Contains(t, output, "missing --remote")
}

func TestMissingFileInSendMode(t *testing.T) {
	output := outputHelper([]string{
		"--mode", "send",
		"--remote", REMOTE_ADDR,
	})
	assert.Contains(t, output, "Invalid input: invalid --file")
}

func TestNonExistentFile(t *testing.T) {
	output := outputHelper([]string{
		"--mode", "send",
		"--remote", REMOTE_ADDR,
		"--file", "404.txt",
	})
	assert.Contains(t, output, "Invalid input: invalid --file")
}

func TestFileTransfer(t *testing.T) {
	fixtures := []string{
		"fixtures/test.txt",
		"fixtures/large_file.bin",
	}

	for _, fixture := range fixtures {
		testFixture(fixture, t)
	}
}

func testFixture(fixture string, t *testing.T) {
	t.Run(fmt.Sprintf("Testing transfer of %s", filepath.Base(fixture)), func(t *testing.T) {
		receiverDone := make(chan struct{})

		// Start receiver in a goroutine
		go func() {
			defer close(receiverDone)
			receiverArgs := []string{
				"--mode", "receive",
				"--port", LOCAL_PORT,
				"--remote", REMOTE_ADDR,
			}
			runWithArgs(receiverArgs)
		}()

		// Wait for receiver to start
		time.Sleep(1 * time.Second)

		// Run sender
		senderArgs := []string{
			"--mode", "send",
			"--port", REMOTE_ADDR[len(REMOTE_ADDR)-4:],
			"--remote", fmt.Sprintf("127.0.0.1:%s", LOCAL_PORT),
			"--file", fixture,
		}
		runWithArgs(senderArgs)

		// Wait for transfer to complete
		select {
		case <-receiverDone:
			// Transfer completed
		case <-time.After(10 * time.Second):
			t.Error("Transfer timed out")
			return
		}

		homeDir, err := os.UserHomeDir(); if err != nil {
			t.Errorf("failed to get home dir: %v", err)
			return
		}

		receivedFile := filepath.Join(homeDir, filepath.Base(fixture)); if !compareFiles(fixture, receivedFile) {
			t.Errorf("File mismatch for %s", fixture)
			t.Errorf("test: %s, %s", receivedFile, fixture)
		}
	})
}

func compareFiles(file1, file2 string) bool {
	hash1, err := fileChecksum(file1); if err != nil {
		return false
	}

	hash2, err := fileChecksum(file2); if err != nil {
		return false
	}

	return hash1 == hash2
}

func fileChecksum(filename string) (string, error) {
	file, err := os.Open(filename); if err != nil {
		return "", err
	}

	defer file.Close()

	return utils.GetFileChecksum(file)
}
