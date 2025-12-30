package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the current Ghostty configuration
type Config struct {
	Path   string
	Values map[string]string
}

// DefaultPath returns the default config file path
func DefaultPath() string {
	// Check XDG config first
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "ghostty", "config")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "ghostty", "config")
}

// Load reads the config file and returns the configuration
func Load(path string) (*Config, error) {
	if path == "" {
		path = DefaultPath()
	}

	cfg := &Config{
		Path:   path,
		Values: make(map[string]string),
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // Return empty config if file doesn't exist
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			cfg.Values[key] = value
		}
	}

	return cfg, scanner.Err()
}

// Save writes the configuration to the file
func (c *Config) Save() error {
	// Read existing file to preserve comments
	existingLines, err := c.readExistingLines()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	file, err := os.Create(c.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	written := make(map[string]bool)

	// Write existing lines, updating values
	for _, line := range existingLines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			fmt.Fprintln(file, line)
			continue
		}

		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			if value, ok := c.Values[key]; ok {
				fmt.Fprintf(file, "%s = %s\n", key, value)
				written[key] = true
			} else {
				fmt.Fprintln(file, line)
			}
		} else {
			fmt.Fprintln(file, line)
		}
	}

	// Append new values
	for key, value := range c.Values {
		if !written[key] {
			fmt.Fprintf(file, "%s = %s\n", key, value)
		}
	}

	return nil
}

func (c *Config) readExistingLines() ([]string, error) {
	file, err := os.Open(c.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Set updates a configuration value
func (c *Config) Set(key, value string) {
	c.Values[key] = value
}

// Get returns a configuration value
func (c *Config) Get(key string) string {
	return c.Values[key]
}
