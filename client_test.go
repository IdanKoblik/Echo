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

func runWithArgs(args []string) error {
	oldArgs := os.Args
	os.Args = append([]string{"cmd"}, args...)
	defer func() { os.Args = oldArgs }()

	return mainEntry()
}

func outputHelper(args []string) (string, error) {
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runWithArgs(args)

	w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String(), err
}

func TestHelpFlag(t *testing.T) {
	err := runWithArgs([]string{"--help"})
	assert.NoError(t, err)
}

func TestInvalidMode(t *testing.T) {
	err := runWithArgs([]string{"--mode", "invalidMode"})
	assert.Error(t, err)
}

func TestMissingRemoteAddress(t *testing.T) {
	_, err := outputHelper([]string{"--mode", "send", "--file", "test.txt"})
	assert.Error(t, err)
}

func TestMissingFileInSendMode(t *testing.T) {
	err := runWithArgs([]string{
		"--mode", "send",
		"--remote", REMOTE_ADDR,
	})
	assert.Error(t, err)
}

func TestInvalidRemote(t *testing.T) {
	err := runWithArgs([]string{
		"--mode", "send",
		"--remote", "12345456",
		"--file", "fixtures/test.txt",
	})

	assert.Error(t, err)
}

func TestNonExistentFile(t *testing.T) {
	err := runWithArgs([]string{
		"--mode", "send",
		"--remote", REMOTE_ADDR,
		"--file", "404.txt",
	})
	assert.Error(t, err)
}

func TestUnreadableFile(t *testing.T) {
	path := "fixtures/unreadable.txt"
	_ = os.WriteFile(path, []byte("test"), 0000)
	defer os.Remove(path)

	err := runWithArgs([]string{
		"--mode", "send",
		"--remote", REMOTE_ADDR,
		"--file", path,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
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

		go func() {
			defer close(receiverDone)
			receiverArgs := []string{
				"--mode", "receive",
				"--port", LOCAL_PORT,
				"--remote", REMOTE_ADDR,
			}
			_ = runWithArgs(receiverArgs)
		}()

		time.Sleep(1 * time.Second)

		senderArgs := []string{
			"--mode", "send",
			"--port", REMOTE_ADDR[len(REMOTE_ADDR)-4:],
			"--remote", fmt.Sprintf("127.0.0.1:%s", LOCAL_PORT),
			"--file", fixture,
		}
		err := runWithArgs(senderArgs)
		assert.NoError(t, err)

		select {
		case <-receiverDone:
		case <-time.After(10 * time.Second):
			t.Error("Transfer timed out")
			return
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Errorf("failed to get home dir: %v", err)
			return
		}

		receivedFile := filepath.Join(homeDir, filepath.Base(fixture))
		if !compareFiles(fixture, receivedFile) {
			t.Errorf("File mismatch for %s", fixture)
		}
	})
}

func compareFiles(file1, file2 string) bool {
	hash1, err := fileChecksum(file1)
	if err != nil {
		return false
	}

	hash2, err := fileChecksum(file2)
	if err != nil {
		return false
	}

	return hash1 == hash2
}

func fileChecksum(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return utils.GetFileChecksum(file)
}
