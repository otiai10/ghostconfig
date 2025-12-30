package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/otiai10/ghostconfig/internal/config"
	"github.com/otiai10/ghostconfig/internal/gui"
	"github.com/otiai10/ghostconfig/internal/schema"
	"github.com/otiai10/ghostconfig/internal/tui"
)

func main() {
	tuiMode := flag.Bool("tui", false, "Use TUI mode (terminal interface)")
	guiMode := flag.Bool("gui", false, "Use GUI mode (web browser interface)")
	port := flag.Int("port", 9999, "Port for GUI server")
	flag.Parse()

	// Parse Ghostty schema
	options, err := schema.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing ghostty config schema: %v\n", err)
		fmt.Fprintf(os.Stderr, "Make sure ghostty is installed and available in PATH\n")
		os.Exit(1)
	}

	if len(options) == 0 {
		fmt.Fprintf(os.Stderr, "No configuration options found\n")
		os.Exit(1)
	}

	// Load current config
	cfg, err := config.Load("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Mode selection: --tui uses TUI, otherwise GUI (default)
	if *tuiMode && !*guiMode {
		runTUI(options, cfg)
	} else {
		runGUI(options, cfg, *port)
	}
}

func runTUI(options []schema.Option, cfg *config.Config) {
	model := tui.New(options, cfg)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

func runGUI(options []schema.Option, cfg *config.Config, port int) {
	server := gui.NewServer(options, cfg, port)
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running GUI server: %v\n", err)
		os.Exit(1)
	}
}
