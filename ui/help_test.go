package ui

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestPrintHelpBox(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintHelpBox()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "Echo File Transfer") {
		t.Errorf("Expected output to contain banner title 'Echo File Transfer'")
	}
	if !strings.Contains(output, "--mode string") {
		t.Errorf("Expected help flags to be printed")
	}
	if !strings.Contains(output, "╔") || !strings.Contains(output, "╝") {
		t.Errorf("Expected box drawing characters in output")
	}

	expectedLines := strings.Count(HELP, "\n") + 4
	actualLines := strings.Count(output, "\n")
	if actualLines != expectedLines {
		t.Errorf("Expected %d lines, got %d", expectedLines, actualLines)
	}
}