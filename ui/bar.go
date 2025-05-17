package ui

import (
	"fmt"

	"github.com/schollz/progressbar/v3"
)

type ProgressBar struct {
	Len int
	Description string
}

func (progressBar *ProgressBar) Init() *progressbar.ProgressBar {
	bar := progressbar.NewOptions(progressBar.Len,
		progressbar.OptionSetDescription(progressBar.Description),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]━[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: "━",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionOnCompletion(func() { fmt.Println() }),
	)

	bar.RenderBlank()
	return bar
} 