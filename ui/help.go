package ui

import (
	"fmt"
	"strings"

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

func PrintHelpBox() {
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