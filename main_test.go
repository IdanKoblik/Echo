package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

const (
	ADDR = "127.0.0.1:8080"
)

func outputHelper(args []string) (string) {
	origStdout := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = args
	main()

	w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	return output
}

func TestInvalidMethod(t *testing.T) {
	output := outputHelper([]string{"cmd", ADDR, "invalidMethod"}, )
	assert.Equal(t, "Invalid method invalidMethod", output)
}

func TestNullFile(t *testing.T) {
	expectedOutput := "Sending file 404.txt to addr: 127.0.0.1:8080\nCannot send file 404.txt to addr: 127.0.0.1:8080\nopen 404.txt: no such file or directory\n"

	output := outputHelper([]string{"cmd", ADDR, "server", "404.txt"})
	assert.Equal(t, expectedOutput, output)
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
		originalFile := fixture
		go func() {
			os.Args = []string{"cmd", ADDR, "client"}
			main()
		}()

		time.Sleep(1 * time.Second)
		os.Args = []string{"cmd", ADDR, "server", fixture}
		main()

		time.Sleep(10 * time.Second)
		homeDir, err := os.UserHomeDir(); if err != nil {
			t.Errorf("failed to get home dir: %v", err)
		}

		receivedFile := filepath.Join(homeDir, filepath.Base(originalFile))
		if !compareFiles(originalFile, receivedFile) {
			t.Errorf("File mismatch for %s\n", originalFile)
			t.Errorf("test: %s\n, %s\n", receivedFile, originalFile)
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

	return GetFileChecksum(file)
}