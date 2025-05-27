package main

import (
	"echo/ui"
	"echo/utils"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

const HELP = `
Usage:
  echo [flags]

Flags:
  --mode string
        Mode of operation: send or receive (optional if using interactive mode)
  --local string
        Local port to listen on (e.g. 9000)
  --remote string
        Remote peer address (e.g. 127.0.0.1:9001)
  --file string
        File path to send (required in send mode)
  --help
        Show this help message and exit
  --bench
        Run benchmarking

Interactive mode will start if no flags are provided.
`

const VERSION = 1

func main() {
	if err := mainEntry(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func mainEntry() error {
	cfg, err := utils.ParseFlags()
	if err != nil {
		return fmt.Errorf("flag parsing failed: %w", err)
	}

	if cfg.HelpMode {
		ui.PrintHelpBox()
		return nil
	}

	if cfg.Mode == "" {
		handleSurveyMode(cfg)
	} else {
		if err := utils.ValidateFlags(cfg); err != nil {
			return fmt.Errorf("invalid input: %w", err)
		}
	}

	if err := RunPeer(cfg); err != nil {
		return fmt.Errorf("run failed: %w", err)
	}

	return nil
}

func handleSurveyMode(cfg *utils.Config, opts ...survey.AskOpt) {
	var selectedMode string
	prompt := &survey.Select{
		Message: "Choose mode:",
		Options: []string{"Send a file", "Receive a file"},
	}
	
	if err := survey.AskOne(prompt, &selectedMode, opts...); err != nil {
		os.Exit(0)
	}

	cfg.Mode = map[string]string{
		"Send a file":    "send",
		"Receive a file": "receive",
	}[selectedMode]

	blue := color.New(color.FgBlue).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	fmt.Printf("\n%s Please choose your settings.\n", bold("CONFIGURATION"))

	if err := survey.AskOne(&survey.Input{
		Message: fmt.Sprintf("%s Enter your local port to listen on (e.g. 9000):", blue(">>")),
		Default: "9000",
	}, &cfg.LocalPort, opts...); err != nil {
		os.Exit(0)
	}

	if err := survey.AskOne(&survey.Input{
		Message: fmt.Sprintf("%s Enter peer's address (e.g. 127.0.0.1:9001):", blue(">>")),
	}, &cfg.RemoteAddr, opts...); err != nil {
		os.Exit(0)
	}

	if cfg.Mode == "send" {
		if err := survey.AskOne(&survey.Input{
			Message: fmt.Sprintf("%s Enter path to the file you want to send:", blue(">>")),
		}, &cfg.FilePath, opts...); err != nil {
			os.Exit(0)
		}
	} else {
		if err := survey.AskOne(&survey.Input{
			Message: fmt.Sprintf("%s Enter destenetion path of the output file:", blue(">>")),
		}, &cfg.OutputDest, opts...); err != nil {
			os.Exit(0)
		}
	}
}

func RunPeer(cfg *utils.Config) error {
	laddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", cfg.LocalPort))
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return err
	}

	defer conn.Close()

	if cfg.FilePath != "" {
		return Send(conn, cfg)
	} else {
		return Receive(conn, cfg)
	}
}
