package main

import (
	"echo/ui"
	"echo/utils"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
)

const VERSION = 3

//go:embed web/dist/**
var embeddedFiles embed.FS

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

	if cfg.Web {
		http.HandleFunc("/ws", WSHandler)
		fmt.Println("WebSocket server started on :8080")

		entries, err := fs.ReadDir(embeddedFiles, "web/dist")
		if err != nil {
			return fmt.Errorf("embed read error: %w", err)
		}

		for _, e := range entries {
			fmt.Println("Embedded file:", e.Name())
		}

		subFS, err := fs.Sub(embeddedFiles, "web/dist")
		if err != nil {
			return fmt.Errorf("failed to create sub filesystem: %w", err)
		}

		staticHandler := http.FileServer(http.FS(subFS))
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			_, err := fs.Stat(subFS, r.URL.Path[1:])
			if err != nil {
				index, err := subFS.Open("index.html")
				if err != nil {
					http.Error(w, "index.html not found", http.StatusInternalServerError)
					return
				}
				defer index.Close()

				data, err := io.ReadAll(index)
				if err != nil {
					http.Error(w, "failed to read index.html", http.StatusInternalServerError)
					return
				}

				w.Write(data)
				return
			}

			staticHandler.ServeHTTP(w, r)
		})

		if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
			return err
		}

		return nil
	}

	if cfg.Mode == "" {
		handleSurveyMode(cfg)
	} else {
		if err := utils.ValidateFlags(cfg); err != nil {
			return fmt.Errorf("invalid input: %w", err)
		}
	}

	localAddr := fmt.Sprintf(":%s", cfg.LocalPort)
	if err := RunPeer(localAddr, cfg.RemoteAddr, cfg.FilePath, cfg.Dest, cfg.BenchMark); err != nil {
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
	survey.AskOne(prompt, &selectedMode, opts...)

	cfg.Mode = map[string]string{
		"Send a file":    "send",
		"Receive a file": "receive",
	}[selectedMode]

	fmt.Printf("\n%s Please choose your settings.\n", "CONFIGURATION")

	survey.AskOne(&survey.Input{
		Message: fmt.Sprintf("%s Enter your local port to listen on (e.g. 9000):", ">>"),
		Default: "9000",
	}, &cfg.LocalPort, opts...)

	survey.AskOne(&survey.Input{
		Message: fmt.Sprintf("%s Enter peer's address (e.g. 127.0.0.1:9001):", ">>"),
	}, &cfg.RemoteAddr, opts...)

	if cfg.Mode == "send" {
		survey.AskOne(&survey.Input{
			Message: fmt.Sprintf("%s Enter path to the file you want to send:", ">>"),
		}, &cfg.FilePath, opts...)
	} else {
		survey.AskOne(&survey.Input{
			Message: fmt.Sprintf("%s Enter output file destention:", ">>"),
		}, &cfg.Dest, opts...)
	}
}

func RunPeer(localAddr, remoteAddr, sendFile, dest string, benchmark bool) error {
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
		return Send(sendFile, conn, remoteAddr, benchmark)
	} else {
		return Receive(conn, benchmark, dest)
	}
}
