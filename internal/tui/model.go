package tui

import (
	"fmt"
	"strings"

	"grua/internal/git"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Pane int

const (
	PaneFileList Pane = iota
	PaneDiffView
)

// Model is the main TUI model.
type Model struct {
	gitService *git.Service
	fileList   *FileList
	diffView   *DiffView
	styles     *Styles
	keys       KeyMap

	activePane  Pane
	showHelp    bool
	width       int
	height      int
	ready       bool
	files       []git.FileStatus
	currentFile *git.FileStatus
	err         error
}

type filesMsg struct {
	files []git.FileStatus
	err   error
}

type diffMsg struct {
	diff *git.FileDiff
	err  error
}

func NewModel(repoPath string) *Model {
	styles := NewStyles()
	keys := DefaultKeyMap()

	return &Model{
		gitService: git.NewService(repoPath),
		fileList:   NewFileList(styles, keys),
		diffView:   NewDiffView(styles, keys),
		styles:     styles,
		keys:       keys,
		activePane: PaneFileList,
	}
}

func (m *Model) Init() tea.Cmd {
	return m.loadFiles
}

func (m *Model) loadFiles() tea.Msg {
	files, err := m.gitService.GetChangedFiles()
	return filesMsg{files: files, err: err}
}

func (m *Model) loadDiff(file git.FileStatus) tea.Cmd {
	return func() tea.Msg {
		diff, err := m.gitService.GetDiff(file.Path, file.Staged)
		return diffMsg{diff: diff, err: err}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.showHelp = !m.showHelp
			return m, nil
		case key.Matches(msg, m.keys.Tab):
			if m.activePane == PaneFileList {
				m.activePane = PaneDiffView
			} else {
				m.activePane = PaneFileList
			}
			return m, nil
		}

		if m.activePane == PaneFileList {
			prevFile := m.fileList.SelectedFile()
			m.fileList, _ = m.fileList.Update(msg)
			newFile := m.fileList.SelectedFile()

			if newFile != nil && (prevFile == nil || *prevFile != *newFile) {
				m.currentFile = newFile
				cmds = append(cmds, m.loadDiff(*newFile))
			}
		} else {
			m.diffView, _ = m.diffView.Update(msg)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()
		m.ready = true

	case filesMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.files = msg.files
		m.fileList.SetFiles(msg.files)

		if len(msg.files) > 0 {
			for i := 0; i < m.fileList.ItemCount(); i++ {
				if file := m.fileList.SelectedFile(); file != nil {
					m.currentFile = file
					cmds = append(cmds, m.loadDiff(*file))
					break
				}
				m.fileList.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
		}

	case diffMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.diffView.SetDiff(msg.diff)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) updateLayout() {
	titleHeight := 3
	statusHeight := 2
	availableHeight := m.height - titleHeight - statusHeight

	fileListWidth := m.width / 4
	if fileListWidth < 20 {
		fileListWidth = 20
	}
	if fileListWidth > 35 {
		fileListWidth = 35
	}
	diffViewWidth := m.width - fileListWidth - 1

	m.fileList.SetSize(fileListWidth, availableHeight)
	m.diffView.SetSize(diffViewWidth, availableHeight)
}

func (m *Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress q to quit.", m.err)
	}

	if m.showHelp {
		return m.renderHelp()
	}

	var b strings.Builder

	title := m.styles.TitleText.Render("🔍 grua - preview changes")
	titleBar := m.styles.TitleBar.
		Width(m.width).
		Render(title)
	b.WriteString(titleBar)
	b.WriteString("\n\n")

	fileListView := m.fileList.View(m.activePane == PaneFileList)
	diffViewView := m.diffView.View(m.activePane == PaneDiffView)

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		fileListView,
		" ",
		diffViewView,
	)
	b.WriteString(content)
	b.WriteString("\n")

	statusBar := m.renderStatusBar()
	b.WriteString(statusBar)

	return b.String()
}

func (m *Model) renderStatusBar() string {
	var items []string

	items = append(items, m.styles.HelpKey.Render("j/k")+" "+m.styles.HelpDesc.Render("up/down"))
	items = append(items, m.styles.HelpKey.Render("Tab")+" "+m.styles.HelpDesc.Render("switch pane"))
	items = append(items, m.styles.HelpKey.Render("g/G")+" "+m.styles.HelpDesc.Render("top/bottom"))
	items = append(items, m.styles.HelpKey.Render("q")+" "+m.styles.HelpDesc.Render("quit"))
	items = append(items, m.styles.HelpKey.Render("?")+" "+m.styles.HelpDesc.Render("help"))

	return m.styles.StatusBar.
		Width(m.width).
		Render(strings.Join(items, "  │  "))
}

func (m *Model) renderHelp() string {
	var b strings.Builder

	title := lipgloss.NewStyle().
		Foreground(ColorTitle).
		Bold(true).
		MarginBottom(1).
		Render("Keyboard Shortcuts")

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(title))
	b.WriteString("\n\n")

	helpItems := []struct {
		key  string
		desc string
	}{
		{"j / k / ↑ / ↓", "Navigate up/down"},
		{"g / G", "Jump to top/bottom"},
		{"Ctrl+u / Ctrl+d", "Page up/down"},
		{"Tab", "Switch between file list and diff view"},
		{"Enter", "Toggle section expand/collapse"},
		{"?", "Toggle this help"},
		{"q / Ctrl+c", "Quit"},
	}

	for _, item := range helpItems {
		keyStyle := lipgloss.NewStyle().
			Foreground(ColorSelected).
			Bold(true).
			Width(20).
			Align(lipgloss.Right)
		descStyle := lipgloss.NewStyle().
			Foreground(ColorFg).
			PaddingLeft(2)

		line := keyStyle.Render(item.key) + descStyle.Render(item.desc)
		b.WriteString(lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	footer := lipgloss.NewStyle().
		Foreground(ColorDim).
		Italic(true).
		Render("Press ? to close")
	b.WriteString(lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(footer))

	return b.String()
}
