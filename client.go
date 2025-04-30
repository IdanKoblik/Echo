package main

import (
	"echo/utils"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"log"
	"net"
)

func main() {
	cfg, err := utils.ParseFlags()
	if err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		return
	}

	if cfg.Mode == "" {
		handleSurveyMode(cfg)
	} else {
		err = utils.ValidateFlags(cfg)
		if err != nil {
			fmt.Printf("Invalid input: %v\n", err)
			return
		}
	}

	localAddr := fmt.Sprintf(":%s", cfg.LocalPort)
	err = RunPeer(localAddr, cfg.RemoteAddr, cfg.FilePath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func handleSurveyMode(cfg *utils.Config, opts ...survey.AskOpt) {
	var selectedMode string
	prompt := &survey.Select{
		Message: "Choose mode:",
		Options: []string{"Send a file", "Receive a file"},
	}
	survey.AskOne(prompt, &selectedMode, opts...)

	cfg.Mode = map[string]string{
		"Send a file":    "send",
		"Receive a file": "receive",
	}[selectedMode]

	survey.AskOne(&survey.Input{
		Message: "Enter your local port to listen on (e.g. 9000):",
		Default: "9000",
	}, &cfg.LocalPort, opts...)

	survey.AskOne(&survey.Input{
		Message: "Enter peer's address (e.g. 127.0.0.1:9001):",
	}, &cfg.RemoteAddr, opts...)

	if cfg.Mode == "send" {
		survey.AskOne(&survey.Input{
			Message: "Enter path to the file you want to send:",
		}, &cfg.FilePath, opts...)
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
