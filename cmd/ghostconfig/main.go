package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/otiai10/ghostconfig/internal/config"
	"github.com/otiai10/ghostconfig/internal/schema"
	"github.com/otiai10/ghostconfig/internal/tui"
)

func main() {
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

	// Start TUI
	model := tui.New(options, cfg)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
