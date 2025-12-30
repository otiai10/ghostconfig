package schema

import (
	"bufio"
	"os/exec"
	"strings"

	"github.com/otiai10/ghostconfig/internal/i18n"
)

// Option represents a single Ghostty configuration option
type Option struct {
	Key          string
	DefaultValue string
	Description  string
	Section      string
}

// Section represents a group of related options
type Section struct {
	Name     string
	Options  []Option
	Expanded bool
}

// Semantic categories for Ghostty configuration (internal keys)
const (
	CategoryFont       = "font"
	CategoryAppearance = "appearance"
	CategoryWindow     = "window"
	CategoryInput      = "input"
	CategoryShell      = "shell"
	CategoryPlatform   = "platform"
	CategoryAdvanced   = "advanced"
)

// CategoryName returns the translated display name for a category
func CategoryName(key string) string {
	return i18n.T("category." + key)
}

// categoryOrder defines the display order of categories
var categoryOrder = []string{
	CategoryAppearance,
	CategoryFont,
	CategoryWindow,
	CategoryInput,
	CategoryShell,
	CategoryPlatform,
	CategoryAdvanced,
}

// ExtractSection extracts the semantic category from a key
func ExtractSection(key string) string {
	// Font: font-*, adjust-*, grapheme-*, freetype-*, alpha-blending
	if strings.HasPrefix(key, "font-") ||
		strings.HasPrefix(key, "adjust-") ||
		strings.HasPrefix(key, "grapheme-") ||
		strings.HasPrefix(key, "freetype-") ||
		key == "alpha-blending" {
		return CategoryFont
	}

	// Appearance: theme, colors, background*, cursor-*, selection-*, palette
	if key == "theme" ||
		key == "background" ||
		key == "foreground" ||
		key == "bold-color" ||
		key == "palette" ||
		key == "minimum-contrast" ||
		key == "faint-opacity" ||
		key == "split-divider-color" ||
		strings.HasPrefix(key, "background-") ||
		strings.HasPrefix(key, "cursor-") ||
		strings.HasPrefix(key, "selection-") {
		return CategoryAppearance
	}

	// Window: window-*, title*, split-*, resize-*, quick-terminal-*, etc.
	if strings.HasPrefix(key, "window-") ||
		strings.HasPrefix(key, "title") ||
		strings.HasPrefix(key, "quick-terminal-") ||
		strings.HasPrefix(key, "resize-") ||
		strings.HasPrefix(key, "unfocused-split-") ||
		key == "class" ||
		key == "fullscreen" ||
		key == "maximize" ||
		key == "initial-window" ||
		key == "confirm-close-surface" {
		return CategoryWindow
	}

	// Input: keybind, mouse-*, clipboard-*, etc.
	if key == "keybind" ||
		key == "input" ||
		key == "copy-on-select" ||
		key == "right-click-action" ||
		key == "focus-follows-mouse" ||
		strings.HasPrefix(key, "mouse-") ||
		strings.HasPrefix(key, "clipboard-") ||
		strings.HasPrefix(key, "click-") {
		return CategoryInput
	}

	// Shell: command, shell-*, working-directory, env, term, scrollback-*
	if key == "command" ||
		key == "initial-command" ||
		key == "working-directory" ||
		key == "env" ||
		key == "term" ||
		key == "wait-after-command" ||
		key == "abnormal-command-exit-runtime" ||
		key == "enquiry-response" ||
		key == "scroll-to-bottom" ||
		strings.HasPrefix(key, "command-") ||
		strings.HasPrefix(key, "shell-") ||
		strings.HasPrefix(key, "scrollback-") {
		return CategoryShell
	}

	// Platform: macos-*, linux-*, gtk-*, x11-*
	if strings.HasPrefix(key, "macos-") ||
		strings.HasPrefix(key, "linux-") ||
		strings.HasPrefix(key, "gtk-") ||
		strings.HasPrefix(key, "x11-") {
		return CategoryPlatform
	}

	// Advanced: everything else
	return CategoryAdvanced
}

// GroupBySection groups options into sections
func GroupBySection(options []Option) []Section {
	sectionMap := make(map[string][]Option)

	for _, opt := range options {
		section := ExtractSection(opt.Key)
		opt.Section = section
		sectionMap[section] = append(sectionMap[section], opt)
	}

	var sections []Section
	for _, key := range categoryOrder {
		if opts, exists := sectionMap[key]; exists {
			sections = append(sections, Section{
				Name:     key, // Return key, translate on client side
				Options:  opts,
				Expanded: false,
			})
		}
	}

	return sections
}

// Parse runs `ghostty +show-config --default --docs` and parses the output
func Parse() ([]Option, error) {
	cmd := exec.Command("ghostty", "+show-config", "--default", "--docs")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parseOutput(string(output)), nil
}

func parseOutput(output string) []Option {
	var options []Option
	var currentDesc strings.Builder
	var currentKey string
	var currentValue string

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "# ") {
			// Comment line - part of description
			currentDesc.WriteString(strings.TrimPrefix(line, "# "))
			currentDesc.WriteString("\n")
		} else if line == "#" {
			// Empty comment line
			currentDesc.WriteString("\n")
		} else if line == "" {
			// Empty line - might be end of an option block
			if currentKey != "" {
				options = append(options, Option{
					Key:          currentKey,
					DefaultValue: currentValue,
					Description:  strings.TrimSpace(currentDesc.String()),
				})
				currentKey = ""
				currentValue = ""
				currentDesc.Reset()
			}
		} else if !strings.HasPrefix(line, "#") {
			// Key = Value line
			parts := strings.SplitN(line, " = ", 2)
			if len(parts) == 2 {
				currentKey = parts[0]
				currentValue = parts[1]
			} else if len(parts) == 1 && strings.HasSuffix(line, " = ") {
				currentKey = strings.TrimSuffix(line, " = ")
				currentValue = ""
			} else {
				// Handle "key = " format
				parts = strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					currentKey = strings.TrimSpace(parts[0])
					currentValue = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	// Don't forget the last option
	if currentKey != "" {
		options = append(options, Option{
			Key:          currentKey,
			DefaultValue: currentValue,
			Description:  strings.TrimSpace(currentDesc.String()),
		})
	}

	return options
}
