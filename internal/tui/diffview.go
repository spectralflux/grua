package tui

import (
	"fmt"
	"strings"

	"grua/internal/git"
	"grua/internal/highlight"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DiffView displays the diff for a selected file.
type DiffView struct {
	diff        *git.FileDiff
	viewport    viewport.Model
	highlighter *highlight.Highlighter
	styles      *Styles
	keys        KeyMap
	width       int
	height      int
	ready       bool
	filePath    string
}

func NewDiffView(styles *Styles, keys KeyMap) *DiffView {
	return &DiffView{
		styles:      styles,
		keys:        keys,
		highlighter: highlight.New(),
	}
}

func (d *DiffView) SetSize(width, height int) {
	d.width = width
	d.height = height

	headerHeight := 1
	viewportHeight := height - headerHeight - 2
	viewportWidth := width - 4

	if !d.ready {
		d.viewport = viewport.New(viewportWidth, viewportHeight)
		d.viewport.YPosition = 0
		d.ready = true
	} else {
		d.viewport.Width = viewportWidth
		d.viewport.Height = viewportHeight
	}

	if d.diff != nil {
		d.renderDiff()
	}
}

func (d *DiffView) SetDiff(diff *git.FileDiff) {
	isNewFile := d.diff == nil || diff == nil ||
		d.diff.Path != diff.Path || d.diff.Staged != diff.Staged

	d.diff = diff
	if diff != nil {
		d.filePath = diff.Path
	} else {
		d.filePath = ""
	}

	prevYOffset := d.viewport.YOffset
	d.renderDiff()

	if isNewFile {
		d.viewport.GotoTop()
	} else {
		maxYOffset := d.viewport.TotalLineCount() - d.viewport.Height
		if maxYOffset < 0 {
			maxYOffset = 0
		}
		if prevYOffset > maxYOffset {
			d.viewport.SetYOffset(maxYOffset)
		} else {
			d.viewport.SetYOffset(prevYOffset)
		}
	}
}

func (d *DiffView) renderDiff() {
	if d.diff == nil || !d.ready {
		d.viewport.SetContent("")
		return
	}

	var lines []string
	contentWidth := d.width - 6

	for _, hunk := range d.diff.Hunks {
		header := d.highlighter.HighlightHunkHeader(hunk.Header)
		lines = append(lines, header)
		lines = append(lines, "")

		for _, line := range hunk.Lines {
			var lineNum string
			switch line.Type {
			case git.LineAdded:
				lineNum = fmt.Sprintf("%4d ", line.NewLineNum)
			case git.LineRemoved:
				lineNum = fmt.Sprintf("%4d ", line.OldLineNum)
			default:
				lineNum = fmt.Sprintf("%4d ", line.NewLineNum)
			}

			lineNumStyled := d.styles.LineNumber.Render(lineNum)

			var hlType highlight.LineType
			switch line.Type {
			case git.LineAdded:
				hlType = highlight.LineAdded
			case git.LineRemoved:
				hlType = highlight.LineRemoved
			default:
				hlType = highlight.LineContext
			}

			content := d.highlighter.HighlightLine(line.Content, hlType, contentWidth)

			var indicator string
			switch line.Type {
			case git.LineAdded:
				indicator = lipgloss.NewStyle().
					Foreground(ColorAddedFg).
					Background(ColorAddedBg).
					Render("+")
			case git.LineRemoved:
				indicator = lipgloss.NewStyle().
					Foreground(ColorRemovedFg).
					Background(ColorRemovedBg).
					Render("-")
			default:
				indicator = " "
			}

			fullLine := lineNumStyled + indicator + " " + content
			lines = append(lines, fullLine)
		}

		lines = append(lines, "")
	}

	d.viewport.SetContent(strings.Join(lines, "\n"))
}

func (d *DiffView) Update(msg tea.Msg) (*DiffView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.keys.Up):
			d.viewport.ScrollUp(1)
		case key.Matches(msg, d.keys.Down):
			d.viewport.ScrollDown(1)
		case key.Matches(msg, d.keys.Top):
			d.viewport.GotoTop()
		case key.Matches(msg, d.keys.Bottom):
			d.viewport.GotoBottom()
		case key.Matches(msg, d.keys.PageUp):
			d.viewport.HalfViewUp()
		case key.Matches(msg, d.keys.PageDown):
			d.viewport.HalfViewDown()
		default:
			d.viewport, cmd = d.viewport.Update(msg)
		}
	default:
		d.viewport, cmd = d.viewport.Update(msg)
	}

	return d, cmd
}

func (d *DiffView) View(active bool) string {
	title := "No file selected"
	if d.filePath != "" {
		title = d.filePath
		if d.diff != nil && d.diff.Staged {
			title += " (staged)"
		}
	}
	titleStyled := d.styles.DiffTitle.Render(title)

	var content string
	if d.diff == nil {
		content = lipgloss.NewStyle().
			Foreground(ColorDim).
			Italic(true).
			Render("Select a file to view diff")
	} else if len(d.diff.Hunks) == 0 {
		content = lipgloss.NewStyle().
			Foreground(ColorDim).
			Italic(true).
			Render("No changes in this file")
	} else {
		content = d.viewport.View()
	}

	fullContent := titleStyled + "\n" + content

	borderStyle := d.styles.DiffBorder
	if active {
		borderStyle = d.styles.DiffBorderActive
	}

	return borderStyle.
		Width(d.width).
		Height(d.height).
		Render(fullContent)
}

func (d *DiffView) ScrollPercent() float64 {
	return d.viewport.ScrollPercent()
}
