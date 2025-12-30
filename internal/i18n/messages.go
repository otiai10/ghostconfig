package i18n

// English messages (default)
var messagesEN = map[string]string{
	// App
	"app.title": "Ghostty Config Editor",

	// Main errors
	"error.parse_schema":     "Error parsing ghostty config schema: %v",
	"error.ghostty_not_found": "Make sure ghostty is installed and available in PATH",
	"error.no_options":       "No configuration options found",
	"error.load_config":      "Error loading config: %v",
	"error.tui":              "Error running TUI: %v",
	"error.gui":              "Error running GUI server: %v",

	// TUI
	"tui.search":          "Search: ",
	"tui.filter":          "Filter: %s (ESC to clear)",
	"tui.new_value":       "New value: ",
	"tui.placeholder":     "Enter value...",
	"tui.search_placeholder": "Search...",
	"tui.select_color":    "Select Color: %s",
	"tui.select_font":     "Select Font: %s",
	"tui.current":         "Current: %s %s",
	"tui.custom":          "Custom: ",
	"tui.custom_color":    "Custom color...",
	"tui.no_fonts":        "No fonts match filter",
	"tui.default":         "(default)",

	// TUI help
	"help.main":        "j/k: move | enter/space: toggle/edit | tab: expand all | /: search | q: quit",
	"help.edit":        "enter: save | esc: cancel",
	"help.search":      "enter: apply | esc: cancel",
	"help.color":       "j/k: move | enter: select | esc: cancel",
	"help.font":        "j/k: move | enter: select | type to filter | esc: cancel (%d fonts)",

	// Messages
	"msg.saved":         "Saved: %s = %s",
	"msg.error":         "Error: %v",
	"msg.loading_fonts": "Error loading fonts: %v",

	// GUI server
	"gui.browser_failed":  "Failed to open browser: %v",
	"gui.open_manually":   "Please open %s manually",
	"gui.starting":        "Starting Ghostty Config GUI at %s",
	"gui.press_ctrl_c":    "Press Ctrl+C to stop",
	"gui.shutting_down":   "Shutting down server...",
	"gui.exit_requested":  "Exit requested, shutting down server...",
	"gui.method_not_allowed": "Method not allowed",

	// GUI frontend
	"gui.search_placeholder": "Search options...",
	"gui.exit":               "Exit",
	"gui.all":                "All",
	"gui.cancel":             "Cancel",
	"gui.save":               "Save",
	"gui.loading":            "Loading options...",
	"gui.no_options":         "No options found",
	"gui.loading_fonts":      "Loading fonts...",
	"gui.filter_fonts":       "Filter fonts...",
	"gui.custom_label":       "Custom:",
	"gui.modified":           "modified",
	"gui.default":            "default",
	"gui.exit_confirm":       "Exit Ghostty Config Editor?",
	"gui.server_stopped":     "Server stopped. You can close this tab.",

	// GUI errors
	"gui.error.load_options": "Failed to load options. Please refresh the page.",
	"gui.error.load_options_api": "Failed to load options",
	"gui.error.load_colors":  "Failed to load colors",
	"gui.error.load_fonts":   "Failed to load fonts",
	"gui.error.save":         "Failed to save config",
	"gui.error.save_prefix":  "Failed to save: ",
	"gui.error.exit":         "Failed to exit",
	"gui.error.init":         "Failed to initialize:",

	// Categories
	"category.font":       "Font",
	"category.appearance": "Appearance",
	"category.window":     "Window",
	"category.input":      "Input",
	"category.shell":      "Shell",
	"category.platform":   "Platform",
	"category.advanced":   "Advanced",
}

// Japanese messages
var messagesJA = map[string]string{
	// App
	"app.title": "Ghostty 設定エディタ",

	// Main errors
	"error.parse_schema":     "Ghostty設定スキーマの解析エラー: %v",
	"error.ghostty_not_found": "ghosttyがインストールされ、PATHに含まれていることを確認してください",
	"error.no_options":       "設定オプションが見つかりません",
	"error.load_config":      "設定の読み込みエラー: %v",
	"error.tui":              "TUI実行エラー: %v",
	"error.gui":              "GUIサーバー実行エラー: %v",

	// TUI
	"tui.search":          "検索: ",
	"tui.filter":          "フィルター: %s (ESCでクリア)",
	"tui.new_value":       "新しい値: ",
	"tui.placeholder":     "値を入力...",
	"tui.search_placeholder": "検索...",
	"tui.select_color":    "色を選択: %s",
	"tui.select_font":     "フォントを選択: %s",
	"tui.current":         "現在: %s %s",
	"tui.custom":          "カスタム: ",
	"tui.custom_color":    "カスタムカラー...",
	"tui.no_fonts":        "一致するフォントがありません",
	"tui.default":         "(デフォルト)",

	// TUI help
	"help.main":        "j/k: 移動 | enter/space: 切替/編集 | tab: 全展開 | /: 検索 | q: 終了",
	"help.edit":        "enter: 保存 | esc: キャンセル",
	"help.search":      "enter: 適用 | esc: キャンセル",
	"help.color":       "j/k: 移動 | enter: 選択 | esc: キャンセル",
	"help.font":        "j/k: 移動 | enter: 選択 | 入力でフィルター | esc: キャンセル (%d フォント)",

	// Messages
	"msg.saved":         "保存しました: %s = %s",
	"msg.error":         "エラー: %v",
	"msg.loading_fonts": "フォント読み込みエラー: %v",

	// GUI server
	"gui.browser_failed":  "ブラウザを開けませんでした: %v",
	"gui.open_manually":   "%s を手動で開いてください",
	"gui.starting":        "Ghostty設定GUIを起動中: %s",
	"gui.press_ctrl_c":    "Ctrl+Cで停止",
	"gui.shutting_down":   "サーバーを停止中...",
	"gui.exit_requested":  "終了リクエスト、サーバーを停止中...",
	"gui.method_not_allowed": "許可されていないメソッドです",

	// GUI frontend
	"gui.search_placeholder": "オプションを検索...",
	"gui.exit":               "終了",
	"gui.all":                "すべて",
	"gui.cancel":             "キャンセル",
	"gui.save":               "保存",
	"gui.loading":            "オプションを読み込み中...",
	"gui.no_options":         "オプションが見つかりません",
	"gui.loading_fonts":      "フォントを読み込み中...",
	"gui.filter_fonts":       "フォントを検索...",
	"gui.custom_label":       "カスタム:",
	"gui.modified":           "変更済",
	"gui.default":            "デフォルト",
	"gui.exit_confirm":       "Ghostty設定エディタを終了しますか？",
	"gui.server_stopped":     "サーバーが停止しました。このタブを閉じてください。",

	// GUI errors
	"gui.error.load_options": "オプションの読み込みに失敗しました。ページを更新してください。",
	"gui.error.load_options_api": "オプションの読み込みに失敗",
	"gui.error.load_colors":  "色の読み込みに失敗",
	"gui.error.load_fonts":   "フォントの読み込みに失敗",
	"gui.error.save":         "設定の保存に失敗",
	"gui.error.save_prefix":  "保存に失敗: ",
	"gui.error.exit":         "終了に失敗",
	"gui.error.init":         "初期化に失敗:",

	// Categories
	"category.font":       "フォント",
	"category.appearance": "外観",
	"category.window":     "ウィンドウ",
	"category.input":      "入力",
	"category.shell":      "シェル",
	"category.platform":   "プラットフォーム",
	"category.advanced":   "詳細設定",
}
