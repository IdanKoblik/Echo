package main

import (
	"os"
	"crypto/sha256"
	"io"
	"encoding/hex"
)

func GetFileChecksum(file *os.File) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum, nil
}