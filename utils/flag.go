package utils

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Mode       string `json:"mode"`
	LocalPort  string `json:"port"`
	RemoteAddr string `json:"remote"`
	HelpMode   bool   `json:"help,omitempty"`
	FilePath   string `json:"file,omitempty"`
	BenchMark  bool   `json:"bench,omitempty"`
	Web        bool   `json:"web,omitempty"`
	Dest       string `json:"dest,omitempty"`
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
	flags.BoolVar(&cfg.Web, "web", false, "Opens echo-ft web ui")
	flags.StringVar(&cfg.Dest, "dest", "", "Path to file to send (required if mode is send)")

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
