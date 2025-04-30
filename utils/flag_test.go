package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func createTempResources(t *testing.T) (string, string) {
	t.Helper()

	tempFile, err := os.CreateTemp("", "test*.txt"); if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	tempFilePath := tempFile.Name()
	tempFile.Close()

	tempDir, err := os.MkdirTemp("", "test"); if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	t.Cleanup(func() {
		os.Remove(tempFilePath)
		os.RemoveAll(tempDir)
	})

	return tempFilePath, tempDir
}

func runValidationTest(t *testing.T, name string, cfg Config, wantErr string) {
	t.Run(name, func(t *testing.T) {
		err := ValidateFlags(&cfg)
		if wantErr == "" {
			if err != nil {
				t.Errorf("ValidateFlags() error = %v, want no error", err)
			}
		} else {
			if err == nil {
				t.Errorf("ValidateFlags() expected error = %v, got no error", wantErr)
			} else if err.Error() != wantErr {
				t.Errorf("ValidateFlags() error = %v, want %v", err, wantErr)
			}
		}
	})
}

func TestValidateFlags(t *testing.T) {
	tempFile, tempDir := createTempResources(t)

	tests := []struct {
		name    string
		mode    string
		port    string
		remote  string
		file    string
		wantErr string
	}{
		{"Valid send configuration", "send", "8080", "localhost:9090", tempFile, ""},
		{"Valid receive configuration", "receive", "8080", "localhost:9090", "", ""},
		{"Invalid mode", "invalid", "8080", "localhost:9090", tempFile, "invalid --mode, must be 'send' or 'receive'"},
		{"Empty mode", "", "8080", "localhost:9090", tempFile, "invalid --mode, must be 'send' or 'receive'"},
		{"Missing remote", "send", "8080", "", tempFile, "invalid or missing --remote, expected format 'host:port'"},
		{"Invalid remote format", "send", "8080", "localhost", tempFile, "invalid or missing --remote, expected format 'host:port'"},
		{"Invalid port (non-numeric)", "send", "abc", "localhost:9090", tempFile, "invalid --port, must be numeric"},
		{"Missing file in send mode", "send", "8080", "localhost:9090", "", "invalid --file, must point to a file"},
		{"Non-existent file in send mode", "send", "8080", "localhost:9090", filepath.Join(tempDir, "nonexistent.txt"), "invalid --file, must point to a file"},
		{"Directory instead of file in send mode", "send", "8080", "localhost:9090", tempDir, "invalid --file, must point to a file"},
	}

	for _, tt := range tests {
		cfg := Config{
			Mode:       tt.mode,
			LocalPort:  tt.port,
			RemoteAddr: tt.remote,
			FilePath:   tt.file,
		}
		runValidationTest(t, tt.name, cfg, tt.wantErr)
	}
}
