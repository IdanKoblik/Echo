package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateFlags(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "test*.txt"); if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test"); if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer os.RemoveAll(tempDir)

	tests := []struct {
		name    string
		mode    string
		port    string
		remote  string
		file    string
		wantErr string
	}{
		{
			name:    "Valid send configuration",
			mode:    "send",
			port:    "8080",
			remote:  "localhost:9090",
			file:    tempFile.Name(),
			wantErr: "",
		},
		{
			name:    "Valid receive configuration",
			mode:    "receive",
			port:    "8080",
			remote:  "localhost:9090",
			file:    "",
			wantErr: "",
		},
		{
			name:    "Invalid mode",
			mode:    "invalid",
			port:    "8080",
			remote:  "localhost:9090",
			file:    tempFile.Name(),
			wantErr: "invalid --mode, must be 'send' or 'receive'",
		},
		{
			name:    "Empty mode",
			mode:    "",
			port:    "8080",
			remote:  "localhost:9090",
			file:    tempFile.Name(),
			wantErr: "invalid --mode, must be 'send' or 'receive'",
		},
		{
			name:    "Missing remote",
			mode:    "send",
			port:    "8080",
			remote:  "",
			file:    tempFile.Name(),
			wantErr: "invalid or missing --remote, expected format 'host:port'",
		},
		{
			name:    "Invalid remote format",
			mode:    "send",
			port:    "8080",
			remote:  "localhost",
			file:    tempFile.Name(),
			wantErr: "invalid or missing --remote, expected format 'host:port'",
		},
		{
			name:    "Invalid port (non-numeric)",
			mode:    "send",
			port:    "abc",
			remote:  "localhost:9090",
			file:    tempFile.Name(),
			wantErr: "invalid --port, must be numeric",
		},
		{
			name:    "Missing file in send mode",
			mode:    "send",
			port:    "8080",
			remote:  "localhost:9090",
			file:    "",
			wantErr: "invalid --file, must point to a file",
		},
		{
			name:    "Non-existent file in send mode",
			mode:    "send",
			port:    "8080",
			remote:  "localhost:9090",
			file:    filepath.Join(tempDir, "nonexistent.txt"),
			wantErr: "invalid --file, must point to a file",
		},
		{
			name:    "Directory instead of file in send mode",
			mode:    "send",
			port:    "8080",
			remote:  "localhost:9090",
			file:    tempDir,
			wantErr: "invalid --file, must point to a file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Mode: tt.mode,
				LocalPort: tt.port,
				RemoteAddr: tt.remote,
				FilePath: tt.file,
			}

			err := ValidateFlags(&cfg)

			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("ValidateFlags() error = %v, want no error", err)
				}
			} else {
				if err == nil {
					t.Errorf("ValidateFlags() expected error = %v, got no error", tt.wantErr)
				} else if err.Error() != tt.wantErr {
					t.Errorf("ValidateFlags() error = %v, want %v", err, tt.wantErr)
				}
			}
		})
	}
}
