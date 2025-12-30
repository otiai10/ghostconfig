package i18n

// English messages (default)
var messagesEN = map[string]string{
	// App
	"app.title": "Ghostty Config Editor",

	// Main errors
	"error.parse_schema":      "Error parsing ghostty config schema: %v",
	"error.ghostty_not_found": "Make sure ghostty is installed and available in PATH",
	"error.no_options":        "No configuration options found",
	"error.load_config":       "Error loading config: %v",
	"error.tui":               "Error running TUI: %v",
	"error.gui":               "Error running GUI server: %v",

	// TUI
	"tui.search":             "Search: ",
	"tui.filter":             "Filter: %s (ESC to clear)",
	"tui.new_value":          "New value: ",
	"tui.placeholder":        "Enter value...",
	"tui.search_placeholder": "Search...",
	"tui.select_color":       "Select Color: %s",
	"tui.select_font":        "Select Font: %s",
	"tui.current":            "Current: %s %s",
	"tui.custom":             "Custom: ",
	"tui.custom_color":       "Custom color...",
	"tui.no_fonts":           "No fonts match filter",
	"tui.default":            "(default)",

	// TUI help
	"help.main":   "j/k: move | enter/space: toggle/edit | tab: expand all | /: search | q: quit",
	"help.edit":   "enter: save | esc: cancel",
	"help.search": "enter: apply | esc: cancel",
	"help.color":  "j/k: move | enter: select | esc: cancel",
	"help.font":   "j/k: move | enter: select | type to filter | esc: cancel (%d fonts)",

	// Messages
	"msg.saved":         "Saved: %s = %s",
	"msg.error":         "Error: %v",
	"msg.loading_fonts": "Error loading fonts: %v",

	// GUI server
	"gui.browser_failed":     "Failed to open browser: %v",
	"gui.open_manually":      "Please open %s manually",
	"gui.starting":           "Starting Ghostty Config GUI at %s",
	"gui.press_ctrl_c":       "Press Ctrl+C to stop",
	"gui.shutting_down":      "Shutting down server...",
	"gui.exit_requested":     "Exit requested, shutting down server...",
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
	"gui.server_stopped":     "Server stopped. You can close this tab.",

	// GUI errors
	"gui.error.load_options":     "Failed to load options. Please refresh the page.",
	"gui.error.load_options_api": "Failed to load options",
	"gui.error.load_colors":      "Failed to load colors",
	"gui.error.load_fonts":       "Failed to load fonts",
	"gui.error.save":             "Failed to save config",
	"gui.error.save_prefix":      "Failed to save: ",
	"gui.error.exit":             "Failed to exit",
	"gui.error.init":             "Failed to initialize:",

	// Categories
	"category.font":       "Font",
	"category.appearance": "Appearance",
	"category.window":     "Window",
	"category.input":      "Input",
	"category.shell":      "Shell",
	"category.platform":   "Platform",
	"category.advanced":   "Advanced",

	// Option descriptions (key options only)
	"desc.font-family": `The font families to use.

You can generate the list of valid values using the CLI: ghostty +list-fonts

This configuration can be repeated multiple times to specify preferred fallback fonts when the requested codepoint is not available in the primary font. This is particularly useful for multiple languages, symbolic fonts, etc.

Notes on emoji: On macOS, Ghostty by default will always use Apple Color Emoji and on Linux will always use Noto Emoji. You can override this behavior by specifying a font family here that contains emoji glyphs.

The specific styles (bold, italic, bold italic) do not need to be explicitly set. If a style is not set, the regular style will be searched for stylistic variants. Some styles may be synthesized if they are not supported.`,

	"desc.font-size": `Font size in points. This value can be a non-integer and the nearest integer pixel size will be selected.

If you have a high dpi display where 1pt = 2px then you can get an odd numbered pixel size by specifying a half point.
Example: 13.5pt @ 2px/pt = 27px

Changing this configuration at runtime will only affect existing terminals that have NOT manually adjusted their font size.

On Linux with GTK, font size is scaled according to both display-wide and text-specific scaling factors.`,

	"desc.theme": `The color theme to use.

To see a list of available themes, run ghostty +list-themes.

A theme file is simply another Ghostty configuration file. Please do not use a theme file from an untrusted source.

To specify a different theme for light and dark mode, use:
light:theme-name,dark:theme-name
Example: light:Rose Pine Dawn,dark:Rose Pine

Any additional colors specified via background, foreground, palette, etc. will override the colors specified in the theme.`,

	"desc.background": `Background color for the window.

Specified as either hex (#RRGGBB or RRGGBB) or a named X11 color.`,

	"desc.foreground": `Foreground color for the window.

Specified as either hex (#RRGGBB or RRGGBB) or a named X11 color.`,

	"desc.cursor-color": `The color of the cursor. If this is not set, a default will be chosen.

Direct colors can be specified as either hex (#RRGGBB or RRGGBB) or a named X11 color.

Additionally, special values can be used:
  * cell-foreground - Match the cell foreground color.
  * cell-background - Match the cell background color.`,

	"desc.cursor-style": `The style of the cursor. This sets the default style. A running program can still request an explicit cursor style using escape sequences.

Note that shell integration will automatically set the cursor to a bar at a prompt. You can disable that behavior by specifying shell-integration-features = no-cursor.

Valid values:
  * block
  * bar
  * underline
  * block_hollow`,

	"desc.window-padding-x": `Horizontal window padding (left and right). It will be scaled appropriately for screen DPI.

If this value is set too large, the screen will render nothing, because the grid will be completely squished by the padding.

To set different left and right padding, specify two values separated by a comma.
Example: window-padding-x = 2,4 (left 2, right 4)

Changing this at runtime will only affect new terminals.`,

	"desc.window-padding-y": `Vertical window padding (top and bottom). It will be scaled appropriately for screen DPI.

If this value is set too large, the screen will render nothing, because the grid will be completely squished by the padding.

To set different top and bottom padding, specify two values separated by a comma.
Example: window-padding-y = 2,4 (top 2, bottom 4)

Changing this at runtime will only affect new terminals.`,

	"desc.window-decoration": `Window decoration display mode.

Valid values:
  * auto - Automatic based on system
  * none - No decorations
  * client - Client-side decorations
  * server - Server-side decorations

For backwards compatibility, true (= auto) and false (= none) are also accepted.

The toggle_window_decorations keybind action can toggle this at runtime.

macOS: To hide the titlebar only, use macos-titlebar-style = hidden instead.`,

	"desc.scrollback-limit": `The size of the scrollback buffer in bytes. This also includes the active screen.

When this limit is reached, the oldest lines are removed from the scrollback.

Scrollback currently exists completely in memory. The larger this value, the larger potential memory usage. Scrollback is allocated lazily up to this limit.

This size is per terminal surface, not for the entire application.

This can be changed at runtime but will only affect new terminal surfaces.`,

	"desc.copy-on-select": `Whether to automatically copy selected text to the clipboard.

true will prefer to copy to the selection clipboard, otherwise it will copy to the system clipboard.

The value clipboard will always copy text to both the selection clipboard and the system clipboard.

Middle-click paste will always use the selection clipboard and is always enabled even if this is false.

The default value is true on Linux and macOS.`,

	"desc.confirm-close-surface": `Confirms that a surface should be closed before closing it.

This defaults to true. If set to false, surfaces will close without any confirmation.

This can also be set to always, which will always confirm closing a surface, even if shell integration says a process isn't running.`,

	"desc.shell-integration": `Whether to enable shell integration auto-injection. Shell integration greatly enhances the terminal experience by enabling:

  * Working directory reporting so new tabs, splits inherit the previous terminal's working directory.
  * Prompt marking that enables the jump_to_prompt keybinding.
  * If you're sitting at a prompt, closing a terminal will not ask for confirmation.
  * Resizing the window with a complex prompt usually paints much better.

Allowable values:
  * none - Do not do any automatic injection. You can still manually configure your shell.
  * detect - Detect the shell based on the filename.
  * bash, elvish, fish, zsh - Use this specific shell injection scheme.

The default value is detect.`,

	"desc.command": `The command to run, usually a shell. If this is not an absolute path, it'll be looked up in the PATH.

If this is not set, a default will be looked up from:
  1. SHELL environment variable
  2. passwd entry (user information)

This can contain additional arguments. If arguments are provided, the command will be executed using /bin/sh -c for shell argument expansion.

To avoid shell expansion altogether, prefix the command with direct:, e.g. direct:nvim foo

This command will be used for all new terminal surfaces. For the first surface only when Ghostty starts, use initial-command.`,
}

// Japanese messages
var messagesJA = map[string]string{
	// App
	"app.title": "Ghostty Config Editor",

	// Main errors
	"error.parse_schema":      "Ghostty設定スキーマの解析エラー: %v",
	"error.ghostty_not_found": "ghosttyがインストールされ、PATHに含まれていることを確認してください",
	"error.no_options":        "設定オプションが見つかりません",
	"error.load_config":       "設定の読み込みエラー: %v",
	"error.tui":               "TUI実行エラー: %v",
	"error.gui":               "GUIサーバー実行エラー: %v",

	// TUI
	"tui.search":             "検索: ",
	"tui.filter":             "フィルター: %s (ESCでクリア)",
	"tui.new_value":          "新しい値: ",
	"tui.placeholder":        "値を入力...",
	"tui.search_placeholder": "検索...",
	"tui.select_color":       "色を選択: %s",
	"tui.select_font":        "フォントを選択: %s",
	"tui.current":            "現在: %s %s",
	"tui.custom":             "カスタム: ",
	"tui.custom_color":       "カスタムカラー...",
	"tui.no_fonts":           "一致するフォントがありません",
	"tui.default":            "(デフォルト)",

	// TUI help
	"help.main":   "j/k: 移動 | enter/space: 切替/編集 | tab: 全展開 | /: 検索 | q: 終了",
	"help.edit":   "enter: 保存 | esc: キャンセル",
	"help.search": "enter: 適用 | esc: キャンセル",
	"help.color":  "j/k: 移動 | enter: 選択 | esc: キャンセル",
	"help.font":   "j/k: 移動 | enter: 選択 | 入力でフィルター | esc: キャンセル (%d フォント)",

	// Messages
	"msg.saved":         "保存しました: %s = %s",
	"msg.error":         "エラー: %v",
	"msg.loading_fonts": "フォント読み込みエラー: %v",

	// GUI server
	"gui.browser_failed":     "ブラウザを開けませんでした: %v",
	"gui.open_manually":      "%s を手動で開いてください",
	"gui.starting":           "Ghostty設定GUIを起動中: %s",
	"gui.press_ctrl_c":       "Ctrl+Cで停止",
	"gui.shutting_down":      "サーバーを停止中...",
	"gui.exit_requested":     "終了リクエスト、サーバーを停止中...",
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
	"gui.server_stopped":     "サーバーが停止しました。このタブを閉じてください。",

	// GUI errors
	"gui.error.load_options":     "オプションの読み込みに失敗しました。ページを更新してください。",
	"gui.error.load_options_api": "オプションの読み込みに失敗",
	"gui.error.load_colors":      "色の読み込みに失敗",
	"gui.error.load_fonts":       "フォントの読み込みに失敗",
	"gui.error.save":             "設定の保存に失敗",
	"gui.error.save_prefix":      "保存に失敗: ",
	"gui.error.exit":             "終了に失敗",
	"gui.error.init":             "初期化に失敗:",

	// Categories
	"category.font":       "フォント",
	"category.appearance": "外観",
	"category.window":     "ウィンドウ",
	"category.input":      "入力",
	"category.shell":      "シェル",
	"category.platform":   "プラットフォーム",
	"category.advanced":   "詳細設定",

	// Option descriptions (key options only)
	"desc.font-family": `使用するフォントファミリー。

利用可能なフォント一覧は ghostty +list-fonts で確認できます。

複数のフォントをフォールバックとして指定可能です。プライマリフォントにないコードポイントは次のフォントで描画されます。多言語や記号フォントに便利です。

絵文字について: macOSではApple Color Emoji、LinuxではNoto Emojiがデフォルトで使用されます。絵文字グリフを含むフォントを指定することでオーバーライド可能です。

太字・イタリック等のスタイルは自動的に検索されます。見つからない場合は通常スタイルが使用されるか、合成されます。`,

	"desc.font-size": `フォントサイズ（ポイント単位）。小数点も指定可能で、最も近い整数ピクセルサイズが選択されます。

高DPIディスプレイ（1pt = 2px）の場合、半ポイントで奇数ピクセルサイズを指定できます。
例: 13.5pt @ 2px/pt = 27px

実行時にこの設定を変更すると、手動でフォントサイズを調整していない既存のターミナルにのみ影響します。

Linux GTKでは、デスクトップ環境のディスプレイスケールとテキストスケーリング設定に従ってスケーリングされます。`,

	"desc.theme": `使用するカラーテーマ。

利用可能なテーマは ghostty +list-themes で確認できます。

テーマファイルは通常のGhostty設定ファイルと同じ形式です。信頼できないソースのテーマは使用しないでください。

ライト/ダークモードで異なるテーマを指定するには:
light:テーマ名,dark:テーマ名
例: light:Rose Pine Dawn,dark:Rose Pine

background、foreground、palette等で追加の色を指定すると、テーマの色をオーバーライドします。`,

	"desc.background": `ウィンドウの背景色。

16進数（#RRGGBB または RRGGBB）またはX11カラー名で指定します。`,

	"desc.foreground": `ウィンドウの前景色（テキストの色）。

16進数（#RRGGBB または RRGGBB）またはX11カラー名で指定します。`,

	"desc.cursor-color": `カーソルの色。設定しない場合はデフォルトが選択されます。

16進数（#RRGGBB または RRGGBB）またはX11カラー名で指定します。

特殊な値も使用可能:
  * cell-foreground - セルの前景色に合わせる
  * cell-background - セルの背景色に合わせる`,

	"desc.cursor-style": `カーソルのスタイル。デフォルトスタイルを設定します。実行中のプログラムはエスケープシーケンスで明示的なスタイルを要求できます。

シェル統合により、プロンプトではカーソルが自動的にバーになります。これを無効にするには shell-integration-features = no-cursor を指定します。

有効な値:
  * block
  * bar
  * underline
  * block_hollow`,

	"desc.window-padding-x": `ウィンドウの水平パディング（左右の余白）。画面DPIに応じてスケーリングされます。

この値が大きすぎると、パディングでグリッドが完全に潰され何も表示されなくなります。

左右で異なるパディングを設定するには、カンマ区切りで2つの値を指定します。
例: window-padding-x = 2,4（左2、右4）

実行時の変更は新しいターミナル（新しいウィンドウ、タブ等）にのみ影響します。`,

	"desc.window-padding-y": `ウィンドウの垂直パディング（上下の余白）。画面DPIに応じてスケーリングされます。

この値が大きすぎると、パディングでグリッドが完全に潰され何も表示されなくなります。

上下で異なるパディングを設定するには、カンマ区切りで2つの値を指定します。
例: window-padding-y = 2,4（上2、下4）

実行時の変更は新しいターミナル（新しいウィンドウ、タブ等）にのみ影響します。`,

	"desc.window-decoration": `ウィンドウ装飾の表示方法。

有効な値:
  * auto - システムに応じて自動選択
  * none - 装飾なし
  * client - クライアントサイド装飾
  * server - サーバーサイド装飾

後方互換性のため、true（= auto）と false（= none）も受け付けます。

toggle_window_decorations キーバインドで実行時に切り替え可能です。

macOS: タイトルバーのみを隠す場合は macos-titlebar-style = hidden を使用してください。`,

	"desc.scrollback-limit": `スクロールバックバッファのサイズ（バイト単位）。アクティブ画面も含みます。

この制限に達すると、最も古い行がスクロールバックから削除されます。

スクロールバックは完全にメモリ上に存在します。値が大きいほどメモリ使用量が増える可能性があります。ただし、このlimitまで遅延割り当てされるため、大きな値を設定しても即座に大量のメモリを消費しません。

このサイズはターミナルサーフェスごとであり、アプリケーション全体ではありません。

実行時の変更は新しいターミナルサーフェスにのみ影響します。`,

	"desc.copy-on-select": `選択したテキストを自動的にクリップボードにコピーするかどうか。

true の場合、選択クリップボードにコピーを優先し、そうでなければシステムクリップボードにコピーします。

clipboard を指定すると、選択クリップボードとシステムクリップボードの両方にコピーします。

中クリック貼り付けは常に選択クリップボードを使用し、この設定が false でも常に有効です。

デフォルトは Linux と macOS で true です。`,

	"desc.confirm-close-surface": `サーフェスを閉じる前に確認を求めるかどうか。

デフォルトは true です。false に設定すると、確認なしで閉じます。

always に設定すると、シェル統合でプロセスが実行されていないと判定されても常に確認します。`,

	"desc.shell-integration": `シェル統合の自動注入を有効にするかどうか。シェル統合は以下の機能を提供します:

  * 作業ディレクトリの報告（新しいタブ/分割が前のターミナルのディレクトリを継承）
  * プロンプトマーキング（jump_to_prompt キーバインドを有効化）
  * プロンプトにいる場合、ターミナルを閉じても確認なし
  * 複雑なプロンプトでのウィンドウリサイズが改善

有効な値:
  * none - 自動注入しない（手動設定は可能）
  * detect - ファイル名からシェルを検出
  * bash, elvish, fish, zsh - 特定のシェル用

デフォルトは detect です。`,

	"desc.command": `実行するコマンド（通常はシェル）。絶対パスでない場合はPATHから検索されます。

設定されていない場合のデフォルト検索順:
  1. SHELL 環境変数
  2. passwd エントリ（ユーザー情報）

追加の引数を含めることができます。引数がある場合、/bin/sh -c でシェル引数展開が実行されます。

シェル展開を避けるには direct: プレフィックスを使用:
例: direct:nvim foo

このコマンドはすべての新しいターミナルサーフェスで使用されます。Ghostty起動時の最初のサーフェスのみに使用するには initial-command を使用してください。`,
}
