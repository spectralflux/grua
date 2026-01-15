# Gruagach - Change Previewer

A TUI tool for previewing Go code changes before committing.

![grua](https://img.shields.io/badge/go-1.25+-blue.svg)

## Installation

```bash
go install grua@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/grua
cd grua
go build -o grua .
```

## Usage

Run `grua` in any git repository with Go file changes:

```bash
cd your-go-project
grua
```

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `j` / `k` / `↑` / `↓` | Navigate up/down |
| `g` / `G` | Jump to top/bottom |
| `Ctrl+u` / `Ctrl+d` | Page up/down |
| `Tab` | Switch between file list and diff view |
| `Enter` | Toggle section expand/collapse |
| `?` | Toggle help |
| `q` / `Ctrl+c` | Quit |

## Name

Named after the changeling in Hellboy, Gruagach https://hellboy.fandom.com/wiki/Gruagach
