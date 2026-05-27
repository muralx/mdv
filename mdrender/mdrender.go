package mdrender

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// osc8Link wraps text in an OSC 8 terminal hyperlink that opens url on click.
// Uses BEL (\x07) as the OSC terminator — more widely supported than ST (\x1b\\),
// especially under tmux/screen and older terminals.
func osc8Link(text, url string) string {
	return "\x1b]8;;" + url + "\x07" + text + "\x1b]8;;\x07"
}

// MarkdownRenderer converts markdown text to ANSI-styled text for terminal display.
type MarkdownRenderer struct {
	H1Style        lipgloss.Style
	H2Style        lipgloss.Style
	H3Style        lipgloss.Style
	BoldStyle      lipgloss.Style
	CodeStyle      lipgloss.Style
	CodeBlockStyle lipgloss.Style
	LinkStyle      lipgloss.Style
}

// New creates a MarkdownRenderer with sensible defaults.
func New() *MarkdownRenderer {
	return &MarkdownRenderer{
		H1Style:        lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#4fc3f7")).Underline(true),
		H2Style:        lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#4fc3f7")),
		H3Style:        lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#81d4fa")),
		BoldStyle:      lipgloss.NewStyle().Bold(true),
		CodeStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#ce9178")),
		CodeBlockStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaaa")),
		LinkStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#4fc3f7")).Underline(true),
	}
}

// Render converts markdown text to ANSI-styled text. When maxWidth > 0,
// any line that would exceed maxWidth visible cells after OSC 8 substitution
// renders its links as plain styled text (no OSC 8) — this avoids triggering
// mtui's broken wrap path for OSC 8 sequences.
func (r *MarkdownRenderer) Render(md string, maxWidth int) string {
	if md == "" {
		return ""
	}

	lines := strings.Split(md, "\n")
	var out []string
	inCodeBlock := false

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}

		if inCodeBlock {
			line = strings.ReplaceAll(line, "\t", "    ")
			out = append(out, r.CodeBlockStyle.Render(line))
			continue
		}

		// Headings (longest prefix first).
		if strings.HasPrefix(line, "### ") {
			out = append(out, r.H3Style.Render(strings.TrimPrefix(line, "### ")))
			continue
		}
		if strings.HasPrefix(line, "## ") {
			out = append(out, r.H2Style.Render(strings.TrimPrefix(line, "## ")))
			continue
		}
		if strings.HasPrefix(line, "# ") {
			out = append(out, r.H1Style.Render(strings.TrimPrefix(line, "# ")))
			continue
		}

		// Horizontal rules.
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" || trimmed == "***" || trimmed == "___" {
			out = append(out, strings.Repeat("─", 60))
			continue
		}

		// Table separator — skip.
		if isTableSeparator(line) {
			continue
		}

		// Table rows — pass through with tab expansion.
		if strings.HasPrefix(strings.TrimSpace(line), "|") {
			out = append(out, strings.ReplaceAll(line, "\t", "    "))
			continue
		}

		// Inline formatting — links first so ANSI from bold/code can't corrupt link detection.
		line = strings.ReplaceAll(line, "\t", "    ")
		line = r.renderLinks(line, maxWidth)
		line = r.renderBold(line)
		line = r.renderInlineCode(line)
		out = append(out, line)
	}

	return strings.Join(out, "\n")
}

func (r *MarkdownRenderer) renderBold(line string) string {
	for {
		start := strings.Index(line, "**")
		if start < 0 {
			break
		}
		end := strings.Index(line[start+2:], "**")
		if end < 0 {
			break
		}
		end += start + 2
		line = line[:start] + r.BoldStyle.Render(line[start+2:end]) + line[end+2:]
	}
	return line
}

func (r *MarkdownRenderer) renderInlineCode(line string) string {
	for {
		start := strings.Index(line, "`")
		if start < 0 {
			break
		}
		end := strings.Index(line[start+1:], "`")
		if end < 0 {
			break
		}
		end += start + 1
		line = line[:start] + r.CodeStyle.Render(line[start+1:end]) + line[end+1:]
	}
	return line
}

// renderLinks replaces [text](url) with OSC 8 hyperlinks styled as links.
// When maxWidth > 0 and the resulting OSC 8 line would exceed maxWidth
// visible cells, falls back to plain styled link text (no OSC 8) for the
// whole line — protects against mtui's wrap path corrupting OSC 8 sequences.
func (r *MarkdownRenderer) renderLinks(line string, maxWidth int) string {
	osc8 := r.applyLinks(line, true)
	if maxWidth > 0 && lipgloss.Width(osc8) > maxWidth {
		return r.applyLinks(line, false)
	}
	return osc8
}

// applyLinks walks [text](url) markdown links left-to-right in a single pass.
// With osc8=true wraps in OSC 8; with osc8=false renders text only (styled).
// Single pass avoids re-scanning emitted ANSI codes — re-scanning was buggy:
// a `[` inside an emitted CSI sequence would pair with `](` from a later real
// link, consuming a huge span and exploding the line length each iteration.
func (r *MarkdownRenderer) applyLinks(line string, osc8 bool) string {
	var out strings.Builder
	out.Grow(len(line))
	i := 0
	for i < len(line) {
		if line[i] != '[' {
			out.WriteByte(line[i])
			i++
			continue
		}
		// Found '['; look for closing '](' then ')' on the remaining input.
		rest := line[i+1:]
		closeBracket := strings.Index(rest, "](")
		if closeBracket < 0 {
			out.WriteByte(line[i])
			i++
			continue
		}
		urlStart := i + 1 + closeBracket + 2
		urlLen := strings.Index(line[urlStart:], ")")
		if urlLen < 0 {
			out.WriteByte(line[i])
			i++
			continue
		}
		text := line[i+1 : i+1+closeBracket]
		url := line[urlStart : urlStart+urlLen]
		styled := r.LinkStyle.Render(text)
		if osc8 {
			styled = osc8Link(styled, url)
		}
		out.WriteString(styled)
		i = urlStart + urlLen + 1
	}
	return out.String()
}

func isTableSeparator(line string) bool {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "|") {
		return false
	}
	cleaned := strings.ReplaceAll(trimmed, "|", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, ":", "")
	return strings.TrimSpace(cleaned) == ""
}
