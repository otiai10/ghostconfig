package gui

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/otiai10/ghostconfig/internal/config"
	"github.com/otiai10/ghostconfig/internal/schema"
)

//go:embed static/*
var staticFiles embed.FS

// Server represents the GUI HTTP server
type Server struct {
	options  []schema.Option
	config   *config.Config
	port     int
	server   *http.Server
	shutdown chan struct{}
}

// NewServer creates a new GUI server
func NewServer(options []schema.Option, cfg *config.Config, port int) *Server {
	return &Server{
		options:  options,
		config:   cfg,
		port:     port,
		shutdown: make(chan struct{}, 1),
	}
}

// Start starts the HTTP server and opens the browser
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/options", s.handleGetOptions)
	mux.HandleFunc("/api/config", s.handleConfig)
	mux.HandleFunc("/api/fonts", s.handleGetFonts)
	mux.HandleFunc("/api/colors", s.handleGetColors)
	mux.HandleFunc("/api/exit", s.handleExit)

	// Static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("failed to create static file system: %w", err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	// Open browser after a short delay
	url := fmt.Sprintf("http://localhost:%d", s.port)
	go func() {
		time.Sleep(200 * time.Millisecond)
		if err := OpenBrowser(url); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open browser: %v\n", err)
			fmt.Fprintf(os.Stderr, "Please open %s manually\n", url)
		}
	}()

	fmt.Printf("Starting Ghostty Config GUI at %s\n", url)
	fmt.Println("Press Ctrl+C to stop")

	// Handle graceful shutdown
	go s.waitForShutdown()

	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		fmt.Println("\nShutting down server...")
	case <-s.shutdown:
		fmt.Println("\nExit requested, shutting down server...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.server.Shutdown(ctx)
}
