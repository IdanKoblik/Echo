package main

import (
	"echo/utils"
	"flag"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"log"
	"net"
	"os"
)

// Build information
var (
	BuildTime = "2025-04-30 10:16:48"
	BuildUser = "IdanKoblik"
)

type Config struct {
	mode       string
	localPort  string
	remoteAddr string
	filePath   string
}

func parseFlags() (*Config, error) {
	flags := flag.NewFlagSet("echo", flag.ExitOnError)

	cfg := &Config{}
	flags.StringVar(&cfg.mode, "mode", "", "Mode: send or receive")
	flags.StringVar(&cfg.localPort, "port", "9000", "Local port to listen on")
	flags.StringVar(&cfg.remoteAddr, "remote", "", "Peer's address (e.g. 127.0.0.1:9001)")
	flags.StringVar(&cfg.filePath, "file", "", "Path to file to send (required if mode is send)")

	if err := flags.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	return cfg, nil
}

func main() {
	cfg, err := parseFlags()
	if err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		return
	}

	useSurvey := (cfg.mode == "" || cfg.remoteAddr == "")

	if !useSurvey {
		err = utils.ValidateFlags(cfg.mode, cfg.localPort, cfg.remoteAddr, cfg.filePath)
		if err != nil {
			fmt.Printf("Invalid input: %v\n", err)
			return
		}
	}

	if useSurvey {
		var selectedMode string
		prompt := &survey.Select{
			Message: "Choose mode:",
			Options: []string{"Send a file", "Receive a file"},
		}
		survey.AskOne(prompt, &selectedMode)

		cfg.mode = map[string]string{
			"Send a file":    "send",
			"Receive a file": "receive",
		}[selectedMode]

		survey.AskOne(&survey.Input{
			Message: "Enter your local port to listen on (e.g. 9000):",
			Default: "9000",
		}, &cfg.localPort)

		survey.AskOne(&survey.Input{
			Message: "Enter peer's address (e.g. 127.0.0.1:9001):",
		}, &cfg.remoteAddr)

		if cfg.mode == "send" {
			survey.AskOne(&survey.Input{
				Message: "Enter path to the file you want to send:",
			}, &cfg.filePath)
		}
	}

	localAddr := fmt.Sprintf(":%s", cfg.localPort)

	err = RunPeer(localAddr, cfg.remoteAddr, cfg.filePath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func RunPeer(localAddr, remoteAddr, sendFile string) error {
	laddr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return err
	}

	defer conn.Close()

	if sendFile != "" {
		return Send(sendFile, conn, remoteAddr)
	} else {
		return Receive(conn)
	}
}
