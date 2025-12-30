package i18n

// Language codes
const (
	LangEN = "en"
	LangJA = "ja"
)

var currentLang = LangEN

// Init initializes the i18n package. Default language is English.
// Language can be changed via GUI language switcher.
func Init() {
	// Default is English (currentLang = LangEN is already set)
}

// SetLang sets the current language
func SetLang(lang string) {
	switch lang {
	case LangJA, "ja_JP", "ja_JP.UTF-8":
		currentLang = LangJA
	default:
		currentLang = LangEN
	}
}

// GetLang returns the current language
func GetLang() string {
	return currentLang
}

// T returns the translated string for the given key
func T(key string) string {
	var msgs map[string]string
	switch currentLang {
	case LangJA:
		msgs = messagesJA
	default:
		msgs = messagesEN
	}

	if msg, ok := msgs[key]; ok {
		return msg
	}
	// Fallback to English
	if msg, ok := messagesEN[key]; ok {
		return msg
	}
	return key
}

// GetMessages returns all messages for the current language (for GUI)
func GetMessages() map[string]string {
	switch currentLang {
	case LangJA:
		return messagesJA
	default:
		return messagesEN
	}
}

// GetAllMessages returns messages for all languages (for GUI language switching)
func GetAllMessages() map[string]map[string]string {
	return map[string]map[string]string{
		LangEN: messagesEN,
		LangJA: messagesJA,
	}
}

// GetAvailableLanguages returns the list of available languages
func GetAvailableLanguages() []string {
	return []string{LangEN, LangJA}
}

// TDesc returns translated description for an option key, with fallback to original
func TDesc(optionKey, originalDesc string) string {
	key := "desc." + optionKey
	translated := T(key)
	// If translation exists (not returning the key itself), use it
	if translated != key {
		return translated
	}
	return originalDesc
}
