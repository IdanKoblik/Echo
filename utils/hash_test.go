package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCalculateChecksum(t *testing.T) {
	testData := []byte("hello world")
	expectedChecksum := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

	checksum := CalculateChecksum(testData)
	if checksum != expectedChecksum {
		t.Errorf("Checksum mismatch: got %v, want %v", checksum, expectedChecksum)
	}

	emptyData := []byte("")
	expectedEmptyChecksum := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	emptyChecksum := CalculateChecksum(emptyData)
	if emptyChecksum != expectedEmptyChecksum {
		t.Errorf("Checksum for empty data mismatch: got %v, want %v", emptyChecksum, expectedEmptyChecksum)
	}
}

func TestGetFileChecksum(t *testing.T) {
	projectRoot := filepath.Join("..")
	testFilePath := filepath.Join(projectRoot, "fixtures", "test.txt")
	file, err := os.Open(testFilePath)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}

	defer file.Close()

	checksum, err := GetFileChecksum(file)
	if err != nil {
		t.Fatalf("GetFileChecksum failed: %v", err)
	}

	expectedChecksum := "678288ba287310f6f225ef73d36a618c2ca2d1d1f6085e19aafc015f96af98d8"
	if checksum != expectedChecksum {
		t.Errorf("Checksum mismatch: got %v, want %v", checksum, expectedChecksum)
	}

	invalidFile, _ := os.Open("nonexistent-file")
	_, err = GetFileChecksum(invalidFile)
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}
