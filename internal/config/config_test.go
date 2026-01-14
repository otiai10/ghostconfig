package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDefaultPath(t *testing.T) {
	path := DefaultPath()
	t.Logf("OS: %s", runtime.GOOS)
	t.Logf("DefaultPath: %s", path)

	// On macOS, if Application Support config exists, it should be returned
	if runtime.GOOS == "darwin" {
		home, _ := os.UserHomeDir()
		appSupportPath := filepath.Join(home, "Library", "Application Support", "com.mitchellh.ghostty", "config")

		if _, err := os.Stat(appSupportPath); err == nil {
			if path != appSupportPath {
				t.Errorf("Expected %s, got %s", appSupportPath, path)
			}
		}
	}
}

func TestGetPrimaryConfigPath(t *testing.T) {
	path := getPrimaryConfigPath()
	t.Logf("OS: %s", runtime.GOOS)
	t.Logf("PrimaryConfigPath: %s", path)

	switch runtime.GOOS {
	case "darwin":
		if path == "" {
			t.Error("Expected non-empty path on macOS")
		}
		home, _ := os.UserHomeDir()
		expected := filepath.Join(home, "Library", "Application Support", "com.mitchellh.ghostty", "config")
		if path != expected {
			t.Errorf("Expected %s, got %s", expected, path)
		}
	case "linux":
		if path != "" {
			t.Error("Expected empty path on Linux")
		}
	}
}
