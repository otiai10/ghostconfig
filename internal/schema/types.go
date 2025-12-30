package schema

import (
	"os/exec"
	"strings"
)

// OptionType represents the type of a configuration option
type OptionType int

const (
	TypeText OptionType = iota
	TypeColor
	TypeFont
	TypeBool
	TypeNumber
)

// GetOptionType returns the type of a configuration option
func GetOptionType(key string) OptionType {
	// Color options
	colorKeys := map[string]bool{
		"background":               true,
		"foreground":               true,
		"bold-color":               true,
		"cursor-color":             true,
		"cursor-text":              true,
		"selection-background":     true,
		"selection-foreground":     true,
		"split-divider-color":      true,
		"window-padding-color":     true,
		"window-titlebar-background": true,
		"window-titlebar-foreground": true,
	}
	if colorKeys[key] || key == "palette" {
		return TypeColor
	}

	// Font options
	if strings.HasPrefix(key, "font-family") {
		return TypeFont
	}

	return TypeText
}

// ListFonts returns available fonts from ghostty
func ListFonts() ([]string, error) {
	cmd := exec.Command("ghostty", "+list-fonts")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var fonts []string
	seen := make(map[string]bool)

	for _, line := range strings.Split(string(output), "\n") {
		// Font family names are not indented
		if line != "" && !strings.HasPrefix(line, " ") {
			name := strings.TrimSpace(line)
			if name != "" && !seen[name] {
				fonts = append(fonts, name)
				seen[name] = true
			}
		}
	}

	return fonts, nil
}

// Common colors for quick selection
var CommonColors = []struct {
	Name  string
	Value string
}{
	{"Black", "000000"},
	{"White", "ffffff"},
	{"Red", "ff0000"},
	{"Green", "00ff00"},
	{"Blue", "0000ff"},
	{"Yellow", "ffff00"},
	{"Cyan", "00ffff"},
	{"Magenta", "ff00ff"},
	{"Orange", "ff8800"},
	{"Purple", "8800ff"},
	{"Gray", "888888"},
	{"Dark Gray", "444444"},
	{"Light Gray", "cccccc"},
}
