# mdv — terminal markdown viewer

[![Go Reference](https://pkg.go.dev/badge/github.com/muralx/mdv.svg)](https://pkg.go.dev/github.com/muralx/mdv)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

`mdv` is a small, full-screen terminal viewer for Markdown files. It pages
large documents, supports mouse-wheel and vim-style scrolling, and renders
inline links as clickable OSC 8 hyperlinks in compatible terminals.

Built on [mate](https://github.com/muralx/mate), a Go TUI component framework.

## Install

```sh
go install github.com/muralx/mdv@latest
```

Or build from source:

```sh
git clone https://github.com/muralx/mdv
cd mdv
go build .
```

## Usage

```sh
mdv path/to/file.md
```

## Controls

| Key            | Action                  |
| -------------- | ----------------------- |
| `↑` / `k`      | Scroll up one line      |
| `↓` / `j`      | Scroll down one line    |
| `PgUp`         | Scroll up one page      |
| `PgDn`         | Scroll down one page    |
| `Home`         | Jump to top             |
| `End`          | Jump to bottom          |
| `Esc` / `q`    | Quit                    |
| `Shift`+drag   | Select and copy text    |

Mouse-wheel scrolling is also supported.

## Rendering

- **Headings** — `#`, `##`, `###`
- **Bold** — `**text**`
- **Inline code** — `` `code` ``
- **Fenced code blocks**
- **Links** — `[text](url)` rendered as clickable OSC 8 hyperlinks
- **Tables** — passed through as-is
- **Horizontal rules** — `---`, `***`, `___`

## Configuration

`MDV_COLOR` selects the color profile when terminal auto-detection isn't
appropriate — for example, inside `tmux`/`screen` with `TERM=screen-*` and
no `COLORTERM`, where auto-detection can block on terminal queries.

| Value                | Profile     |
| -------------------- | ----------- |
| _(unset)_            | True color  |
| `ansi256`            | 256 colors  |
| `ansi`               | 16 colors   |
| `off` / `none` / `ascii` | No color  |

## License

[Apache 2.0](LICENSE)
