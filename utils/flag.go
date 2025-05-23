package utils

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HelpMode   bool
	Mode       string
	LocalPort  string
	RemoteAddr string
	FilePath   string
	BenchMark  bool
}

func ParseFlags() (*Config, error) {
	flags := flag.NewFlagSet("echo", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	cfg := &Config{}
	flags.BoolVar(&cfg.HelpMode, "help", false, "Help mode")
	flags.StringVar(&cfg.Mode, "mode", "", "Mode: send or receive")
	flags.StringVar(&cfg.LocalPort, "port", "9000", "Local port to listen on")
	flags.StringVar(&cfg.RemoteAddr, "remote", "", "Peer's address (e.g. 127.0.0.1:9001)")
	flags.StringVar(&cfg.FilePath, "file", "", "Path to file to send (required if mode is send)")
	flags.BoolVar(&cfg.BenchMark, "bench", false, "Run benchmarking")

	if err := flags.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	return cfg, nil
}

func ValidateFlags(cfg *Config) error {
	if cfg.Mode != "send" && cfg.Mode != "receive" {
		return errors.New("invalid --mode, must be 'send' or 'receive'")
	}

	if cfg.RemoteAddr == "" || !strings.Contains(cfg.RemoteAddr, ":") {
		return errors.New("invalid or missing --remote, expected format 'host:port'")
	}

	if _, err := strconv.Atoi(cfg.LocalPort); err != nil {
		return errors.New("invalid --port, must be numeric")
	}

	if cfg.Mode == "send" {
		info, err := os.Stat(cfg.FilePath)
		if err != nil || info.IsDir() {
			return errors.New("invalid --file, must point to a file")
		}
	}

	return nil
}
