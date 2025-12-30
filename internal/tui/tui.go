package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/otiai10/ghostconfig/internal/config"
	"github.com/otiai10/ghostconfig/internal/schema"
)

type mode int

const (
	modeList mode = iota
	modeEdit
	modeSearch
)

// ListItem represents either a section header or an option
type ListItem struct {
	IsSection    bool
	SectionIndex int
	OptionIndex  int
}

type Model struct {
	sections    []schema.Section
	config      *config.Config
	cursor      int
	offset      int
	height      int
	mode        mode
	textInput   textinput.Model
	searchQuery string
	items       []ListItem // flattened list of visible items
	message     string
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212"))

	sectionSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("57"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57"))

	keyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("228"))

	defaultStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			MarginTop(1).
			MarginLeft(2)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("82")).
			Bold(true)

	countStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)

func New(options []schema.Option, cfg *config.Config) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter value..."
	ti.CharLimit = 256
	ti.Width = 50

	sections := schema.GroupBySection(options)

	m := Model{
		sections:  sections,
		config:    cfg,
		height:    20,
		textInput: ti,
	}

	m.rebuildItems()
	return m
}

func (m *Model) rebuildItems() {
	m.items = nil
	query := strings.ToLower(m.searchQuery)

	for si, section := range m.sections {
		// Filter options if searching
		var matchingOpts []int
		for oi, opt := range section.Options {
			if query == "" ||
				strings.Contains(strings.ToLower(opt.Key), query) ||
				strings.Contains(strings.ToLower(opt.Description), query) {
				matchingOpts = append(matchingOpts, oi)
			}
		}

		// Skip section if no matching options
		if len(matchingOpts) == 0 && query != "" {
			continue
		}

		// Add section header
		m.items = append(m.items, ListItem{
			IsSection:    true,
			SectionIndex: si,
		})

		// Add options if expanded (or searching)
		if section.Expanded || query != "" {
			for _, oi := range matchingOpts {
				m.items = append(m.items, ListItem{
					IsSection:    false,
					SectionIndex: si,
					OptionIndex:  oi,
				})
			}
		}
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case modeList:
			return m.updateList(msg)
		case modeEdit:
			return m.updateEdit(msg)
		case modeSearch:
			return m.updateSearch(msg)
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height - 10
		if m.height < 5 {
			m.height = 5
		}
	}

	return m, nil
}

func (m Model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			if m.cursor < m.offset {
				m.offset = m.cursor
			}
		}

	case "down", "j":
		if m.cursor < len(m.items)-1 {
			m.cursor++
			if m.cursor >= m.offset+m.height {
				m.offset = m.cursor - m.height + 1
			}
		}

	case "enter", " ":
		if len(m.items) > 0 && m.cursor < len(m.items) {
			item := m.items[m.cursor]
			if item.IsSection {
				// Toggle section expansion
				m.sections[item.SectionIndex].Expanded = !m.sections[item.SectionIndex].Expanded
				m.rebuildItems()
				// Keep cursor in bounds
				if m.cursor >= len(m.items) {
					m.cursor = len(m.items) - 1
				}
			} else {
				// Edit option
				m.mode = modeEdit
				opt := m.sections[item.SectionIndex].Options[item.OptionIndex]
				currentVal := m.config.Get(opt.Key)
				if currentVal == "" {
					currentVal = opt.DefaultValue
				}
				m.textInput.SetValue(currentVal)
				m.textInput.Focus()
				return m, textinput.Blink
			}
		}

	case "tab":
		// Expand/collapse all
		allExpanded := true
		for _, s := range m.sections {
			if !s.Expanded {
				allExpanded = false
				break
			}
		}
		for i := range m.sections {
			m.sections[i].Expanded = !allExpanded
		}
		m.rebuildItems()

	case "/":
		m.mode = modeSearch
		m.textInput.SetValue(m.searchQuery)
		m.textInput.Placeholder = "Search..."
		m.textInput.Focus()
		return m, textinput.Blink

	case "esc":
		if m.searchQuery != "" {
			m.searchQuery = ""
			m.rebuildItems()
			m.cursor = 0
			m.offset = 0
		}
	}

	return m, nil
}

func (m Model) updateEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeList
		m.textInput.Blur()
		return m, nil

	case "enter":
		item := m.items[m.cursor]
		opt := m.sections[item.SectionIndex].Options[item.OptionIndex]
		newValue := m.textInput.Value()
		m.config.Set(opt.Key, newValue)
		if err := m.config.Save(); err != nil {
			m.message = fmt.Sprintf("Error: %v", err)
		} else {
			m.message = fmt.Sprintf("Saved: %s = %s", opt.Key, newValue)
		}
		m.mode = modeList
		m.textInput.Blur()
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeList
		m.textInput.Blur()
		return m, nil

	case "enter":
		m.searchQuery = m.textInput.Value()
		m.rebuildItems()
		m.cursor = 0
		m.offset = 0
		m.mode = modeList
		m.textInput.Blur()
		m.textInput.Placeholder = "Enter value..."
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Ghostty Config Editor"))
	b.WriteString("\n")

	if m.mode == modeSearch {
		b.WriteString("Search: ")
		b.WriteString(m.textInput.View())
		b.WriteString("\n\n")
	} else if m.searchQuery != "" {
		b.WriteString(fmt.Sprintf("Filter: %s (ESC to clear)\n\n", m.searchQuery))
	} else {
		b.WriteString("\n")
	}

	// List items
	end := m.offset + m.height
	if end > len(m.items) {
		end = len(m.items)
	}

	for i := m.offset; i < end; i++ {
		item := m.items[i]
		isSelected := i == m.cursor

		if item.IsSection {
			section := m.sections[item.SectionIndex]
			icon := "▶"
			if section.Expanded || m.searchQuery != "" {
				icon = "▼"
			}
			count := countStyle.Render(fmt.Sprintf("(%d)", len(section.Options)))
			line := fmt.Sprintf("%s %s %s", icon, section.Name, count)

			if isSelected {
				b.WriteString(sectionSelectedStyle.Render(line))
			} else {
				b.WriteString(sectionStyle.Render(line))
			}
		} else {
			opt := m.sections[item.SectionIndex].Options[item.OptionIndex]
			currentVal := m.config.Get(opt.Key)

			var line string
			if isSelected {
				val := currentVal
				if val == "" {
					val = opt.DefaultValue + " (default)"
				}
				line = selectedStyle.Render(fmt.Sprintf("  > %s = %s", opt.Key, val))
			} else {
				if currentVal != "" {
					line = fmt.Sprintf("    %s = %s", keyStyle.Render(opt.Key), valueStyle.Render(currentVal))
				} else {
					line = fmt.Sprintf("    %s = %s", keyStyle.Render(opt.Key), defaultStyle.Render(opt.DefaultValue+" (default)"))
				}
			}
			b.WriteString(line)
		}
		b.WriteString("\n")
	}

	// Show description for selected option
	if len(m.items) > 0 && m.cursor < len(m.items) {
		item := m.items[m.cursor]
		if !item.IsSection {
			opt := m.sections[item.SectionIndex].Options[item.OptionIndex]
			desc := opt.Description
			if len(desc) > 200 {
				desc = desc[:200] + "..."
			}
			b.WriteString(descStyle.Render(desc))
			b.WriteString("\n")
		}
	}

	// Edit mode
	if m.mode == modeEdit {
		b.WriteString("\n")
		b.WriteString("New value: ")
		b.WriteString(m.textInput.View())
		b.WriteString("\n")
	}

	// Message
	if m.message != "" {
		b.WriteString("\n")
		b.WriteString(messageStyle.Render(m.message))
	}

	// Help
	help := "j/k: move | enter/space: toggle/edit | tab: expand all | /: search | q: quit"
	if m.mode == modeEdit {
		help = "enter: save | esc: cancel"
	} else if m.mode == modeSearch {
		help = "enter: apply | esc: cancel"
	}
	b.WriteString(helpStyle.Render("\n" + help))

	return b.String()
}
