# Gruagach - Change Review Gremlin

A TUI tool for reviewing Go code changes before committing. Great for finding AI nonsense before it
gets committed!

![grua](https://img.shields.io/badge/go-1.25+-blue.svg)

## Installation
Build from source:

```bash
git clone https://github.com/spectralflux/grua
cd grua
go build -o grua .
```
and place executable on PATH somewhere.

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
| `?` | Toggle help |
| `q` / `Ctrl+c` | Quit |

## Name

Named after the changeling in gaelic mythology and Hellboy, Gruagach https://hellboy.fandom.com/wiki/Gruagach
