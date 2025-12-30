package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/otiai10/ghostconfig/internal/config"
	"github.com/otiai10/ghostconfig/internal/i18n"
	"github.com/otiai10/ghostconfig/internal/schema"
)

type mode int

const (
	modeList mode = iota
	modeEdit
	modeSearch
	modeColorPicker
	modeFontPicker
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
	items       []ListItem
	message     string

	// For color picker
	colorCursor int
	customColor bool

	// For font picker
	fonts       []string
	fontCursor  int
	fontOffset  int
	fontFilter  string
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

	colorSwatchStyle = lipgloss.NewStyle().
				Width(4).
				Height(1)

	pickerSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("57"))

	pickerItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))
)

func New(options []schema.Option, cfg *config.Config) Model {
	ti := textinput.New()
	ti.Placeholder = i18n.T("tui.placeholder")
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
		var matchingOpts []int
		for oi, opt := range section.Options {
			if query == "" ||
				strings.Contains(strings.ToLower(opt.Key), query) ||
				strings.Contains(strings.ToLower(opt.Description), query) {
				matchingOpts = append(matchingOpts, oi)
			}
		}

		if len(matchingOpts) == 0 && query != "" {
			continue
		}

		m.items = append(m.items, ListItem{
			IsSection:    true,
			SectionIndex: si,
		})

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
		case modeColorPicker:
			return m.updateColorPicker(msg)
		case modeFontPicker:
			return m.updateFontPicker(msg)
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
				m.sections[item.SectionIndex].Expanded = !m.sections[item.SectionIndex].Expanded
				m.rebuildItems()
				if m.cursor >= len(m.items) {
					m.cursor = len(m.items) - 1
				}
			} else {
				opt := m.sections[item.SectionIndex].Options[item.OptionIndex]
				optType := schema.GetOptionType(opt.Key)

				switch optType {
				case schema.TypeColor:
					m.mode = modeColorPicker
					m.colorCursor = 0
					m.customColor = false
					currentVal := m.config.Get(opt.Key)
					if currentVal == "" {
						currentVal = opt.DefaultValue
					}
					m.textInput.SetValue(currentVal)
					return m, nil

				case schema.TypeFont:
					fonts, err := schema.ListFonts()
					if err != nil {
						m.message = fmt.Sprintf(i18n.T("msg.loading_fonts"), err)
						return m, nil
					}
					m.fonts = fonts
					m.mode = modeFontPicker
					m.fontCursor = 0
					m.fontOffset = 0
					m.fontFilter = ""
					// Find current font in list
					currentVal := m.config.Get(opt.Key)
					for i, f := range fonts {
						if f == currentVal {
							m.fontCursor = i
							if m.fontCursor >= m.height-2 {
								m.fontOffset = m.fontCursor - m.height/2
							}
							break
						}
					}
					return m, nil

				default:
					m.mode = modeEdit
					currentVal := m.config.Get(opt.Key)
					if currentVal == "" {
						currentVal = opt.DefaultValue
					}
					m.textInput.SetValue(currentVal)
					m.textInput.Focus()
					return m, textinput.Blink
				}
			}
		}

	case "tab":
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
		m.textInput.Placeholder = i18n.T("tui.search_placeholder")
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
			m.message = fmt.Sprintf(i18n.T("msg.error"), err)
		} else {
			m.message = fmt.Sprintf(i18n.T("msg.saved"), opt.Key, newValue)
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
		m.textInput.Placeholder = i18n.T("tui.placeholder")
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) updateColorPicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	totalItems := len(schema.CommonColors) + 1 // +1 for custom

	switch msg.String() {
	case "esc":
		m.mode = modeList
		return m, nil

	case "up", "k":
		if m.colorCursor > 0 {
			m.colorCursor--
			m.customColor = false
		}

	case "down", "j":
		if m.colorCursor < totalItems-1 {
			m.colorCursor++
		}
		if m.colorCursor == totalItems-1 {
			m.customColor = true
			m.textInput.Focus()
			return m, textinput.Blink
		}

	case "enter":
		item := m.items[m.cursor]
		opt := m.sections[item.SectionIndex].Options[item.OptionIndex]

		var newValue string
		if m.customColor {
			newValue = m.textInput.Value()
		} else {
			newValue = schema.CommonColors[m.colorCursor].Value
		}

		m.config.Set(opt.Key, newValue)
		if err := m.config.Save(); err != nil {
			m.message = fmt.Sprintf(i18n.T("msg.error"), err)
		} else {
			m.message = fmt.Sprintf(i18n.T("msg.saved"), opt.Key, newValue)
		}
		m.mode = modeList
		m.textInput.Blur()
		return m, nil
	}

	if m.customColor {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) updateFontPicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	filteredFonts := m.getFilteredFonts()
	maxVisible := m.height - 4

	switch msg.String() {
	case "esc":
		m.mode = modeList
		m.fontFilter = ""
		return m, nil

	case "up", "k":
		if m.fontCursor > 0 {
			m.fontCursor--
			if m.fontCursor < m.fontOffset {
				m.fontOffset = m.fontCursor
			}
		}

	case "down", "j":
		if m.fontCursor < len(filteredFonts)-1 {
			m.fontCursor++
			if m.fontCursor >= m.fontOffset+maxVisible {
				m.fontOffset = m.fontCursor - maxVisible + 1
			}
		}

	case "enter":
		if len(filteredFonts) > 0 && m.fontCursor < len(filteredFonts) {
			item := m.items[m.cursor]
			opt := m.sections[item.SectionIndex].Options[item.OptionIndex]
			newValue := filteredFonts[m.fontCursor]

			m.config.Set(opt.Key, newValue)
			if err := m.config.Save(); err != nil {
				m.message = fmt.Sprintf(i18n.T("msg.error"), err)
			} else {
				m.message = fmt.Sprintf(i18n.T("msg.saved"), opt.Key, newValue)
			}
			m.mode = modeList
			m.fontFilter = ""
			return m, nil
		}

	case "backspace":
		if len(m.fontFilter) > 0 {
			m.fontFilter = m.fontFilter[:len(m.fontFilter)-1]
			m.fontCursor = 0
			m.fontOffset = 0
		}

	default:
		// Add character to filter
		if len(msg.String()) == 1 {
			m.fontFilter += msg.String()
			m.fontCursor = 0
			m.fontOffset = 0
		}
	}

	return m, nil
}

func (m Model) getFilteredFonts() []string {
	if m.fontFilter == "" {
		return m.fonts
	}

	filter := strings.ToLower(m.fontFilter)
	var filtered []string
	for _, f := range m.fonts {
		if strings.Contains(strings.ToLower(f), filter) {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(i18n.T("app.title")))
	b.WriteString("\n")

	switch m.mode {
	case modeColorPicker:
		return m.viewColorPicker()
	case modeFontPicker:
		return m.viewFontPicker()
	}

	if m.mode == modeSearch {
		b.WriteString(i18n.T("tui.search"))
		b.WriteString(m.textInput.View())
		b.WriteString("\n\n")
	} else if m.searchQuery != "" {
		b.WriteString(fmt.Sprintf(i18n.T("tui.filter")+"\n\n", m.searchQuery))
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
			sectionName := schema.CategoryName(section.Name)
			line := fmt.Sprintf("%s %s %s", icon, sectionName, count)

			if isSelected {
				b.WriteString(sectionSelectedStyle.Render(line))
			} else {
				b.WriteString(sectionStyle.Render(line))
			}
		} else {
			opt := m.sections[item.SectionIndex].Options[item.OptionIndex]
			currentVal := m.config.Get(opt.Key)
			optType := schema.GetOptionType(opt.Key)

			var line string
			val := currentVal
			if val == "" {
				val = opt.DefaultValue
			}

			// Add color swatch for color options
			var colorSwatch string
			if optType == schema.TypeColor && val != "" {
				colorSwatch = colorSwatchStyle.Background(lipgloss.Color(val)).Render("  ") + " "
			}

			if isSelected {
				displayVal := val
				if currentVal == "" {
					displayVal = val + " " + i18n.T("tui.default")
				}
				line = selectedStyle.Render(fmt.Sprintf("  > %s = %s", opt.Key, displayVal))
				if colorSwatch != "" {
					line = colorSwatch + line
				}
			} else {
				if currentVal != "" {
					line = fmt.Sprintf("    %s = %s", keyStyle.Render(opt.Key), valueStyle.Render(currentVal))
				} else {
					line = fmt.Sprintf("    %s = %s", keyStyle.Render(opt.Key), defaultStyle.Render(val+" "+i18n.T("tui.default")))
				}
				if colorSwatch != "" {
					line = colorSwatch + line
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
			desc := i18n.TDesc(opt.Key, opt.Description)
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
		b.WriteString(i18n.T("tui.new_value"))
		b.WriteString(m.textInput.View())
		b.WriteString("\n")
	}

	// Message
	if m.message != "" {
		b.WriteString("\n")
		b.WriteString(messageStyle.Render(m.message))
	}

	// Help
	help := i18n.T("help.main")
	if m.mode == modeEdit {
		help = i18n.T("help.edit")
	} else if m.mode == modeSearch {
		help = i18n.T("help.search")
	}
	b.WriteString(helpStyle.Render("\n" + help))

	return b.String()
}

func (m Model) viewColorPicker() string {
	var b strings.Builder

	item := m.items[m.cursor]
	opt := m.sections[item.SectionIndex].Options[item.OptionIndex]

	b.WriteString(titleStyle.Render(fmt.Sprintf(i18n.T("tui.select_color"), opt.Key)))
	b.WriteString("\n\n")

	// Current value preview
	currentVal := m.config.Get(opt.Key)
	if currentVal == "" {
		currentVal = opt.DefaultValue
	}
	if currentVal != "" {
		preview := colorSwatchStyle.Background(lipgloss.Color(currentVal)).Render("    ")
		b.WriteString(fmt.Sprintf("Current: %s %s\n\n", preview, currentVal))
	}

	// Common colors
	for i, c := range schema.CommonColors {
		swatch := colorSwatchStyle.Background(lipgloss.Color(c.Value)).Render("  ")
		line := fmt.Sprintf("%s %s (%s)", swatch, c.Name, c.Value)

		if i == m.colorCursor && !m.customColor {
			b.WriteString(pickerSelectedStyle.Render("> " + line))
		} else {
			b.WriteString("  " + pickerItemStyle.Render(line))
		}
		b.WriteString("\n")
	}

	// Custom color option
	b.WriteString("\n")
	if m.customColor {
		b.WriteString(pickerSelectedStyle.Render("> " + i18n.T("tui.custom")))
		b.WriteString(m.textInput.View())
	} else {
		if m.colorCursor == len(schema.CommonColors) {
			b.WriteString(pickerSelectedStyle.Render("> " + i18n.T("tui.custom_color")))
		} else {
			b.WriteString("  " + pickerItemStyle.Render(i18n.T("tui.custom_color")))
		}
	}
	b.WriteString("\n")

	b.WriteString(helpStyle.Render("\n" + i18n.T("help.color")))

	return b.String()
}

func (m Model) viewFontPicker() string {
	var b strings.Builder

	item := m.items[m.cursor]
	opt := m.sections[item.SectionIndex].Options[item.OptionIndex]

	b.WriteString(titleStyle.Render(fmt.Sprintf(i18n.T("tui.select_font"), opt.Key)))
	b.WriteString("\n")

	// Current value
	currentVal := m.config.Get(opt.Key)
	if currentVal != "" {
		b.WriteString(fmt.Sprintf("Current: %s\n", valueStyle.Render(currentVal)))
	}

	// Filter display
	if m.fontFilter != "" {
		b.WriteString(fmt.Sprintf("Filter: %s\n", m.fontFilter))
	}
	b.WriteString("\n")

	filteredFonts := m.getFilteredFonts()
	maxVisible := m.height - 4

	end := m.fontOffset + maxVisible
	if end > len(filteredFonts) {
		end = len(filteredFonts)
	}

	for i := m.fontOffset; i < end; i++ {
		font := filteredFonts[i]
		if i == m.fontCursor {
			// Show font name with preview using the font itself (if terminal supports it)
			b.WriteString(pickerSelectedStyle.Render(fmt.Sprintf("> %s", font)))
		} else {
			b.WriteString(fmt.Sprintf("  %s", pickerItemStyle.Render(font)))
		}
		b.WriteString("\n")
	}

	if len(filteredFonts) == 0 {
		b.WriteString(defaultStyle.Render("  " + i18n.T("tui.no_fonts") + "\n"))
	}

	b.WriteString(helpStyle.Render(fmt.Sprintf("\n"+i18n.T("help.font"), len(filteredFonts))))

	return b.String()
}
