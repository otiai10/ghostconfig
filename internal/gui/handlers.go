package gui

import (
	"encoding/json"
	"net/http"

	"github.com/otiai10/ghostconfig/internal/i18n"
	"github.com/otiai10/ghostconfig/internal/schema"
)

// OptionResponse represents an option in the API response
type OptionResponse struct {
	Key          string `json:"key"`
	DefaultValue string `json:"defaultValue"`
	Description  string `json:"description"`
	Section      string `json:"section"`
	Type         string `json:"type"`
	CurrentValue string `json:"currentValue"`
}

// SectionResponse represents a section with its options
type SectionResponse struct {
	Name    string           `json:"name"`
	Options []OptionResponse `json:"options"`
}

// GET /api/options - Get all options grouped by section
func (s *Server) handleGetOptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, i18n.T("gui.method_not_allowed"), http.StatusMethodNotAllowed)
		return
	}

	sections := schema.GroupBySection(s.options)
	var response []SectionResponse

	for _, section := range sections {
		sectionData := SectionResponse{
			Name:    section.Name,
			Options: make([]OptionResponse, 0, len(section.Options)),
		}

		for _, opt := range section.Options {
			optType := schema.GetOptionType(opt.Key)
			typeStr := "text"
			switch optType {
			case schema.TypeColor:
				typeStr = "color"
			case schema.TypeFont:
				typeStr = "font"
			case schema.TypeBool:
				typeStr = "bool"
			case schema.TypeNumber:
				typeStr = "number"
			}

			sectionData.Options = append(sectionData.Options, OptionResponse{
				Key:          opt.Key,
				DefaultValue: opt.DefaultValue,
				Description:  opt.Description,
				Section:      section.Name,
				Type:         typeStr,
				CurrentValue: s.config.Get(opt.Key),
			})
		}
		response = append(response, sectionData)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET/PUT /api/config - Get or update config values
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.config.Values)

	case http.MethodPut:
		var req struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.config.Set(req.Key, req.Value)
		if err := s.config.Save(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

	default:
		http.Error(w, i18n.T("gui.method_not_allowed"), http.StatusMethodNotAllowed)
	}
}

// GET /api/fonts - Get available fonts
func (s *Server) handleGetFonts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, i18n.T("gui.method_not_allowed"), http.StatusMethodNotAllowed)
		return
	}

	fonts, err := schema.ListFonts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fonts)
}

// ColorOption represents a preset color
type ColorOption struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// GET /api/colors - Get preset colors
func (s *Server) handleGetColors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, i18n.T("gui.method_not_allowed"), http.StatusMethodNotAllowed)
		return
	}

	colors := make([]ColorOption, len(schema.CommonColors))
	for i, c := range schema.CommonColors {
		colors[i] = ColorOption{
			Name:  c.Name,
			Value: c.Value,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(colors)
}

// POST /api/exit - Shutdown the server
func (s *Server) handleExit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, i18n.T("gui.method_not_allowed"), http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

	// Trigger shutdown
	go func() {
		s.shutdown <- struct{}{}
	}()
}

// GET /api/i18n - Get i18n messages for all languages
func (s *Server) handleGetI18n(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, i18n.T("gui.method_not_allowed"), http.StatusMethodNotAllowed)
		return
	}

	response := struct {
		DefaultLang string                       `json:"defaultLang"`
		Languages   []string                     `json:"languages"`
		Messages    map[string]map[string]string `json:"messages"`
	}{
		DefaultLang: i18n.GetLang(),
		Languages:   i18n.GetAvailableLanguages(),
		Messages:    i18n.GetAllMessages(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
