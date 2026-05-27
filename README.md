# mdv

A full-screen terminal markdown viewer built with [mate](https://github.com/muralx/mate).

## Usage

```
mdv <file.md>
```

## Controls

| Key | Action |
| --- | --- |
| `↑` / `k` | Scroll up one line |
| `↓` / `j` | Scroll down one line |
| `PgUp` | Scroll up one page |
| `PgDn` | Scroll down one page |
| `Home` | Jump to top |
| `End` | Jump to bottom |
| `Esc` / `q` | Quit |

Mouse wheel scrolling is also supported.

## Rendering

mdv renders the following markdown elements:

- **Headings** — `#`, `##`, `###` rendered in cyan bold
- **Bold** — `**text**`
- **Inline code** — `` `code` `` rendered in orange
- **Fenced code blocks** — rendered in grey
- **Tables** — passed through as-is
- **Horizontal rules** — `---`, `***`, `___`

## Building

```
go build .
```

Or install to `$GOPATH/bin`:

```
go install .
```

## Dependencies

- [mate](https://github.com/muralx/mate) — TUI component framework
- [bubbletea](https://github.com/charmbracelet/bubbletea) — terminal app framework
- [lipgloss](https://github.com/charmbracelet/lipgloss) — terminal styling
