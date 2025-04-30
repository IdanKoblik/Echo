package utils

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func ValidateFlags(mode, port, remote, file string) error {
	if mode != "send" && mode != "receive" {
		return errors.New("invalid --mode, must be 'send' or 'receive'")
	}

	if remote == "" || !strings.Contains(remote, ":") {
		return errors.New("invalid or missing --remote, expected format 'host:port'")
	}

	if _, err := strconv.Atoi(port); err != nil {
		return errors.New("invalid --port, must be numeric")
	}

	if mode == "send" {
		info, err := os.Stat(file)
		if err != nil || info.IsDir() {
			return errors.New("invalid --file, must point to a file")
		}
	}
	return nil
}
