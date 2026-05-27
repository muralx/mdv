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

	win, text := buildWindow(filepath.Base(path))
	text.SetMarkdown(string(data))
	app := window.NewApp(win)
	p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "mdv: %v\n", err)
		os.Exit(1)
	}
}

func buildWindow(title string) (*window.MainWindow, *widget.MarkdownTextArea) {
	win := window.NewWindow("mdv")

	text := widget.NewMarkdownTextArea("content", widget.DefaultMarkdownTextAreaStyles())

	statusStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1a1a2e")).
		Foreground(lipgloss.Color("#555577")).
		Padding(0, 1)
	status := widget.NewText("status",
		title+" | ↑↓/jk: scroll | PgUp/PgDn: page | Home/End: top/bottom | Esc: quit | Shift to select and copy text",
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
