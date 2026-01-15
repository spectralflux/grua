package tui

import "github.com/charmbracelet/lipgloss"

var (
	ColorBorder      = lipgloss.Color("#44475A")
	ColorTitle       = lipgloss.Color("#FFD700")
	ColorStaged      = lipgloss.Color("#FF79C6")
	ColorUnstaged    = lipgloss.Color("#8BE9FD")
	ColorSelected    = lipgloss.Color("#BD93F9")
	ColorHunk        = lipgloss.Color("#00D7FF")
	ColorLineNum     = lipgloss.Color("#6272A4")
	ColorAddedBg     = lipgloss.Color("#1B4B1B")
	ColorRemovedBg   = lipgloss.Color("#4B1818")
	ColorAddedFg     = lipgloss.Color("#69FF94")
	ColorRemovedFg   = lipgloss.Color("#FF6B6B")
	ColorDim         = lipgloss.Color("#6272A4")
	ColorFg          = lipgloss.Color("#F8F8F2")
	ColorBg          = lipgloss.Color("#282A36")
	ColorHighlight   = lipgloss.Color("#44475A")
	ColorStatusBadge = lipgloss.Color("#50FA7B")
	ColorStatusBarBg = lipgloss.Color("#1E1F29")

	LogoGradient = []lipgloss.Color{
		lipgloss.Color("#E9B8FF"),
		lipgloss.Color("#D896FF"),
		lipgloss.Color("#C778FF"),
		lipgloss.Color("#B65EFF"),
		lipgloss.Color("#A855F7"),
		lipgloss.Color("#9333EA"),
		lipgloss.Color("#7E22CE"),
	}
)

// Styles holds all lipgloss styles for the TUI.
type Styles struct {
	App                  lipgloss.Style
	TitleBar             lipgloss.Style
	TitleText            lipgloss.Style
	FileListBorder       lipgloss.Style
	FileListBorderActive lipgloss.Style
	StagedHeader         lipgloss.Style
	UnstagedHeader       lipgloss.Style
	FileItem             lipgloss.Style
	FileItemSelected     lipgloss.Style
	StatusBadge          lipgloss.Style
	DiffBorder           lipgloss.Style
	DiffBorderActive     lipgloss.Style
	DiffTitle            lipgloss.Style
	HunkHeader           lipgloss.Style
	LineNumber           lipgloss.Style
	AddedLine            lipgloss.Style
	RemovedLine          lipgloss.Style
	ContextLine          lipgloss.Style
	StatusBar            lipgloss.Style
	HelpKey              lipgloss.Style
	HelpDesc             lipgloss.Style
}

func NewStyles() *Styles {
	s := &Styles{}

	s.App = lipgloss.NewStyle().Background(ColorBg)

	s.TitleBar = lipgloss.NewStyle().
		Foreground(ColorFg).
		Background(ColorBg).
		Padding(0, 1).
		Bold(true).
		Align(lipgloss.Center)

	s.TitleText = lipgloss.NewStyle().
		Foreground(ColorTitle).
		Bold(true)

	s.FileListBorder = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1)

	s.FileListBorderActive = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSelected).
		Padding(0, 1)

	s.StagedHeader = lipgloss.NewStyle().
		Foreground(ColorStaged).
		Bold(true).
		MarginBottom(0)

	s.UnstagedHeader = lipgloss.NewStyle().
		Foreground(ColorUnstaged).
		Bold(true).
		MarginTop(1).
		MarginBottom(0)

	s.FileItem = lipgloss.NewStyle().
		Foreground(ColorFg).
		PaddingLeft(2)

	s.FileItemSelected = lipgloss.NewStyle().
		Foreground(ColorBg).
		Background(ColorSelected).
		PaddingLeft(2).
		Bold(true)

	s.StatusBadge = lipgloss.NewStyle().
		Foreground(ColorStatusBadge).
		PaddingLeft(1)

	s.DiffBorder = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1)

	s.DiffBorderActive = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSelected).
		Padding(0, 1)

	s.DiffTitle = lipgloss.NewStyle().
		Foreground(ColorTitle).
		Bold(true)

	s.HunkHeader = lipgloss.NewStyle().
		Foreground(ColorHunk).
		Bold(true)

	s.LineNumber = lipgloss.NewStyle().
		Foreground(ColorLineNum).
		Width(6).
		Align(lipgloss.Right).
		PaddingRight(1)

	s.AddedLine = lipgloss.NewStyle().
		Background(ColorAddedBg).
		Foreground(ColorAddedFg)

	s.RemovedLine = lipgloss.NewStyle().
		Background(ColorRemovedBg).
		Foreground(ColorRemovedFg)

	s.ContextLine = lipgloss.NewStyle().
		Foreground(ColorFg)

	s.StatusBar = lipgloss.NewStyle().
		Background(ColorStatusBarBg).
		Foreground(ColorDim).
		Padding(0, 1)

	s.HelpKey = lipgloss.NewStyle().
		Background(ColorStatusBarBg).
		Foreground(ColorSelected).
		Bold(true)

	s.HelpDesc = lipgloss.NewStyle().
		Background(ColorStatusBarBg).
		Foreground(lipgloss.Color("#8B8B9E"))

	return s
}
