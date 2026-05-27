# mdv — terminal markdown viewer

[![CI](https://github.com/muralx/mdv/actions/workflows/ci.yml/badge.svg)](https://github.com/muralx/mdv/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/muralx/mdv)](https://github.com/muralx/mdv/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/muralx/mdv.svg)](https://pkg.go.dev/github.com/muralx/mdv)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`mdv` is a small, full-screen terminal viewer for Markdown files. It pages
large documents, supports mouse-wheel and vim-style scrolling, and renders
inline links as clickable OSC 8 hyperlinks in compatible terminals.

Built on [mate](https://github.com/muralx/mate), a Go TUI component framework.

## Install

### Prebuilt binaries

Grab the archive for your OS/arch from the
[latest release](https://github.com/muralx/mdv/releases/latest), extract,
and place `mdv` somewhere on your `PATH`.

### With Go

```sh
go install github.com/muralx/mdv@latest
```

### From source

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

[MIT](LICENSE)
