// mdv renders a markdown file in a full-screen terminal viewer.
//
// Usage: mdv <file.md>
//
// Controls:
//
//	↑↓ / j k: scroll one line
//	PgUp/PgDn: scroll one page
//	Home/End: jump to top/bottom
//	Esc / q: quit
package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/muralx/mate/widget"
	"github.com/muralx/mate/window"

	"mdv/mdrender"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: mdv <file.md>")
		os.Exit(1)
	}

	// Pin lipgloss/termenv to a known color profile before any Render call.
	// Auto-detection would query the terminal (DCS sequences) and block forever
	// when running inside tmux/screen with TERM=screen-* and no COLORTERM.
	profile := termenv.TrueColor
	switch os.Getenv("MDV_COLOR") {
	case "ansi":
		profile = termenv.ANSI
	case "ansi256":
		profile = termenv.ANSI256
	case "off", "none", "ascii":
		profile = termenv.Ascii
	}
	lipgloss.SetColorProfile(profile)

	path := os.Args[1]
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mdv: %v\n", err)
		os.Exit(1)
	}

	renderer := mdrender.New()
	md := string(data)

	win, text := buildWindow(filepath.Base(path))
	text.SetContent(renderer.Render(md, 0))
	app := window.NewApp(win)
	v := &viewer{App: app, md: md, text: text, renderer: renderer}
	p := tea.NewProgram(v, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "mdv: %v\n", err)
		os.Exit(1)
	}
}

// viewer wraps window.App so we can re-render markdown when the terminal
// resizes — passing the new width to mdrender lets it suppress OSC 8 on
// lines that would otherwise be too long and trigger mtui's broken wrap path.
type viewer struct {
	*window.App
	md       string
	text     *widget.ScrollableText
	renderer *mdrender.MarkdownRenderer
	width    int
}

func (v *viewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if sm, ok := msg.(tea.WindowSizeMsg); ok && sm.Width != v.width {
		v.width = sm.Width
		v.text.SetContent(v.renderer.Render(v.md, sm.Width))
	}
	_, cmd := v.App.Update(msg)
	return v, cmd
}

func buildWindow(title string) (*window.MainWindow, *widget.ScrollableText) {
	win := window.NewWindow("mdv")

	text := widget.NewScrollableText("content", widget.DefaultScrollableTextStyles())

	statusStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1a1a2e")).
		Foreground(lipgloss.Color("#555577")).
		Padding(0, 1)
	status := widget.NewText("status",
		title+" | ↑↓/jk: scroll | PgUp/PgDn: page | Home/End: top/bottom | Esc: quit",
		statusStyle)
	status.SetPreferredHeight(1)

	win.Add(text, widget.TCBCenter)
	win.Add(status, widget.TCBBottom)

	win.OnKeyPress(func(msg tea.KeyMsg) tea.Cmd {
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return tea.Quit
		}
		return nil
	})

	return win, text
}
