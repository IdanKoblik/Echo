package main

import (
	"echo/utils"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"log"
	"strings"
	"net"
)

const HELP = `
Usage:
  echo [flags]

Flags:
  -mode string
        Mode of operation: send or receive (optional if using interactive mode)
  -local string
        Local port to listen on (e.g. 9000)
  -remote string
        Remote peer address (e.g. 127.0.0.1:9001)
  -file string
        File path to send (required in send mode)
  -help, -h
        Show this help message and exit

Interactive mode will start if no flags are provided.
`

func main() {
	cfg, err := utils.ParseFlags()
	if err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		return
	}

	if cfg.HelpMode {
		printHelpBox()
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

	blue := color.New(color.FgBlue).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	fmt.Printf("\n%s Please choose your settings.\n", bold("CONFIGURATION"))

	survey.AskOne(&survey.Input{
		Message: fmt.Sprintf("%s Enter your local port to listen on (e.g. 9000):", blue(">>")),
		Default: "9000",
	}, &cfg.LocalPort, opts...)

	survey.AskOne(&survey.Input{
		Message: fmt.Sprintf("%s Enter peer's address (e.g. 127.0.0.1:9001):", blue(">>")),
	}, &cfg.RemoteAddr, opts...)

	if cfg.Mode == "send" {
		survey.AskOne(&survey.Input{
			Message: fmt.Sprintf("%s Enter path to the file you want to send:", blue(">>")),
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

func printHelpBox() {
	boxColor := color.New(color.FgHiBlue, color.Bold)
	textColor := color.New(color.FgHiWhite)

	lines := strings.Split(HELP, "\n")
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	banner := color.New(color.FgGreen, color.Bold).Sprint(" Echo File Transfer ")

	boxColor.Println("╔" + strings.Repeat("═", maxWidth+2) + "╗")
	fmt.Printf("║%s%s║\n", banner, strings.Repeat(" ", maxWidth-len("Echo File Transfer")))
	boxColor.Println("╠" + strings.Repeat("═", maxWidth+2) + "╣")
	for _, line := range lines {
		textColor.Printf("║ %-*s ║\n", maxWidth, line)
	}

	boxColor.Println("╚" + strings.Repeat("═", maxWidth+2) + "╝")
}