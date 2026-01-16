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
}

// FileList is the file list component.
type FileList struct {
	items  []FileListItem
	cursor int
	width  int
	height int
	styles *Styles
	keys   KeyMap
}

func NewFileList(styles *Styles, keys KeyMap) *FileList {
	return &FileList{
		styles: styles,
		keys:   keys,
	}
}

func (f *FileList) SetSize(width, height int) {
	f.width = width
	f.height = height
}

func (f *FileList) SetFiles(files []git.FileStatus) {
	prevSelected := f.SelectedFile()

	f.items = nil

	var staged, unstaged, unversioned []git.FileStatus
	for _, file := range files {
		if file.Unversioned {
			unversioned = append(unversioned, file)
		} else if file.Staged {
			staged = append(staged, file)
		} else {
			unstaged = append(unstaged, file)
		}
	}

	if len(staged) > 0 {
		f.items = append(f.items, FileListItem{
			IsHeader:   true,
			HeaderText: "STAGED",
		})
		for _, file := range staged {
			f.items = append(f.items, FileListItem{File: file})
		}
	}

	if len(unstaged) > 0 {
		f.items = append(f.items, FileListItem{
			IsHeader:   true,
			HeaderText: "UNSTAGED",
		})
		for _, file := range unstaged {
			f.items = append(f.items, FileListItem{File: file})
		}
	}

	if len(unversioned) > 0 {
		f.items = append(f.items, FileListItem{
			IsHeader:   true,
			HeaderText: "UNVERSIONED",
		})
		for _, file := range unversioned {
			f.items = append(f.items, FileListItem{File: file})
		}
	}

	if prevSelected != nil {
		for i, item := range f.items {
			if !item.IsHeader && item.File.Path == prevSelected.Path &&
				item.File.Staged == prevSelected.Staged &&
				item.File.Unversioned == prevSelected.Unversioned {
				f.cursor = i
				return
			}
		}
	}

	f.cursor = f.firstFileIndex()
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
			f.cursor = f.firstFileIndex()
		case key.Matches(msg, f.keys.Bottom):
			f.cursor = f.lastFileIndex()
		}
	}
	return f, nil
}

func (f *FileList) moveUp() {
	start := f.cursor
	for {
		f.cursor--
		if f.cursor < 0 {
			f.cursor = len(f.items) - 1
		}
		if f.cursor == start {
			break
		}
		if !f.items[f.cursor].IsHeader {
			break
		}
	}
}

func (f *FileList) moveDown() {
	start := f.cursor
	for {
		f.cursor++
		if f.cursor >= len(f.items) {
			f.cursor = 0
		}
		if f.cursor == start {
			break
		}
		if !f.items[f.cursor].IsHeader {
			break
		}
	}
}

func (f *FileList) firstFileIndex() int {
	for i, item := range f.items {
		if !item.IsHeader {
			return i
		}
	}
	return 0
}

func (f *FileList) lastFileIndex() int {
	for i := len(f.items) - 1; i >= 0; i-- {
		if !f.items[i].IsHeader {
			return i
		}
	}
	return 0
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
		isSelected := i == f.cursor && !item.IsHeader

		if item.IsHeader {
			headerStyle := f.styles.StagedHeader
			switch item.HeaderText {
			case "UNSTAGED":
				headerStyle = f.styles.UnstagedHeader
			case "UNVERSIONED":
				headerStyle = f.styles.UnversionedHeader
			}
			line = headerStyle.Render(fmt.Sprintf(" ▾ %s", item.HeaderText))
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
