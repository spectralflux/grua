package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	"grua/internal/git"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FileListItem struct {
	File       git.FileStatus
	IsHeader   bool
	HeaderText string
	Expanded   bool
}

// FileList is the file list component.
type FileList struct {
	items            []FileListItem
	cursor           int
	width            int
	height           int
	styles           *Styles
	keys             KeyMap
	stagedExpanded   bool
	unstagedExpanded bool
}

func NewFileList(styles *Styles, keys KeyMap) *FileList {
	return &FileList{
		styles:           styles,
		keys:             keys,
		stagedExpanded:   true,
		unstagedExpanded: true,
	}
}

func (f *FileList) SetSize(width, height int) {
	f.width = width
	f.height = height
}

func (f *FileList) SetFiles(files []git.FileStatus) {
	f.items = nil

	var staged, unstaged []git.FileStatus
	for _, file := range files {
		if file.Staged {
			staged = append(staged, file)
		} else {
			unstaged = append(unstaged, file)
		}
	}

	if len(staged) > 0 {
		f.items = append(f.items, FileListItem{
			IsHeader:   true,
			HeaderText: "STAGED",
			Expanded:   f.stagedExpanded,
		})
		if f.stagedExpanded {
			for _, file := range staged {
				f.items = append(f.items, FileListItem{File: file})
			}
		}
	}

	if len(unstaged) > 0 {
		f.items = append(f.items, FileListItem{
			IsHeader:   true,
			HeaderText: "UNSTAGED",
			Expanded:   f.unstagedExpanded,
		})
		if f.unstagedExpanded {
			for _, file := range unstaged {
				f.items = append(f.items, FileListItem{File: file})
			}
		}
	}

	if f.cursor >= len(f.items) {
		f.cursor = max(0, len(f.items)-1)
	}
}

func (f *FileList) Update(msg tea.Msg) (*FileList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, f.keys.Up):
			f.moveUp()
		case key.Matches(msg, f.keys.Down):
			f.moveDown()
		case key.Matches(msg, f.keys.Top):
			f.cursor = 0
		case key.Matches(msg, f.keys.Bottom):
			f.cursor = max(0, len(f.items)-1)
		case key.Matches(msg, f.keys.Enter):
			f.toggleExpanded()
		}
	}
	return f, nil
}

func (f *FileList) moveUp() {
	if f.cursor > 0 {
		f.cursor--
	}
}

func (f *FileList) moveDown() {
	if f.cursor < len(f.items)-1 {
		f.cursor++
	}
}

func (f *FileList) toggleExpanded() {
	if f.cursor < 0 || f.cursor >= len(f.items) {
		return
	}
	item := f.items[f.cursor]
	if item.IsHeader {
		if item.HeaderText == "STAGED" {
			f.stagedExpanded = !f.stagedExpanded
		} else {
			f.unstagedExpanded = !f.unstagedExpanded
		}
	}
}

func (f *FileList) SelectedFile() *git.FileStatus {
	if f.cursor < 0 || f.cursor >= len(f.items) {
		return nil
	}
	item := f.items[f.cursor]
	if item.IsHeader {
		return nil
	}
	return &item.File
}

func (f *FileList) View(active bool) string {
	if len(f.items) == 0 {
		emptyMsg := lipgloss.NewStyle().
			Foreground(ColorDim).
			Italic(true).
			Render("No changes")
		return emptyMsg
	}

	var lines []string
	for i, item := range f.items {
		var line string
		isSelected := i == f.cursor

		if item.IsHeader {
			arrow := "▾"
			if (item.HeaderText == "STAGED" && !f.stagedExpanded) ||
				(item.HeaderText == "UNSTAGED" && !f.unstagedExpanded) {
				arrow = "▸"
			}

			headerStyle := f.styles.StagedHeader
			if item.HeaderText == "UNSTAGED" {
				headerStyle = f.styles.UnstagedHeader
			}

			if isSelected {
				line = lipgloss.NewStyle().
					Foreground(ColorBg).
					Background(ColorSelected).
					Bold(true).
					Render(fmt.Sprintf(" %s %s ", arrow, item.HeaderText))
			} else {
				line = headerStyle.Render(fmt.Sprintf(" %s %s", arrow, item.HeaderText))
			}
		} else {
			filename := filepath.Base(item.File.Path)
			status := item.File.Status

			maxNameLen := f.width - 8
			if maxNameLen < 10 {
				maxNameLen = 10
			}
			if len(filename) > maxNameLen {
				filename = filename[:maxNameLen-3] + "..."
			}

			paddedName := fmt.Sprintf("%-*s", maxNameLen, filename)

			if isSelected {
				line = f.styles.FileItemSelected.
					Width(f.width - 4).
					Render(fmt.Sprintf("%s %s", paddedName, status))
			} else {
				line = f.styles.FileItem.Render(paddedName) +
					f.styles.StatusBadge.Render(status)
			}
		}

		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")

	borderStyle := f.styles.FileListBorder
	if active {
		borderStyle = f.styles.FileListBorderActive
	}

	return borderStyle.
		Width(f.width).
		Height(f.height).
		Render(content)
}

func (f *FileList) Cursor() int {
	return f.cursor
}

func (f *FileList) ItemCount() int {
	return len(f.items)
}
