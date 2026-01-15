package highlight

import (
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
)

type LineType int

const (
	LineContext LineType = iota
	LineAdded
	LineRemoved
)

var (
	AddedBg       = lipgloss.Color("#1B4B1B")
	RemovedBg     = lipgloss.Color("#4B1818")
	AddedFg       = lipgloss.Color("#69FF94")
	RemovedFg     = lipgloss.Color("#FF6B6B")
	ColorKeyword  = lipgloss.Color("#FF79C6")
	ColorString   = lipgloss.Color("#F1FA8C")
	ColorComment  = lipgloss.Color("#6272A4")
	ColorFunction = lipgloss.Color("#50FA7B")
	ColorType     = lipgloss.Color("#8BE9FD")
	ColorNumber   = lipgloss.Color("#BD93F9")
	ColorOperator = lipgloss.Color("#FF79C6")
	ColorDefault  = lipgloss.Color("#F8F8F2")
	ColorHunk     = lipgloss.Color("#00D7FF")
)

// Highlighter provides Go syntax highlighting with diff support.
type Highlighter struct {
	lexer chroma.Lexer
	style *chroma.Style
}

// New creates a new Go syntax highlighter.
func New() *Highlighter {
	lexer := lexers.Get("go")
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get("dracula")
	if style == nil {
		style = styles.Fallback
	}

	return &Highlighter{
		lexer: lexer,
		style: style,
	}
}

func (h *Highlighter) tokenColor(tt chroma.TokenType) lipgloss.Color {
	entry := h.style.Get(tt)
	if entry.Colour.IsSet() {
		return lipgloss.Color(entry.Colour.String())
	}

	switch {
	case tt == chroma.Keyword || tt == chroma.KeywordDeclaration ||
		tt == chroma.KeywordNamespace || tt == chroma.KeywordType:
		return ColorKeyword
	case tt == chroma.String || tt == chroma.StringChar || tt == chroma.StringBacktick:
		return ColorString
	case tt == chroma.Comment || tt == chroma.CommentSingle || tt == chroma.CommentMultiline:
		return ColorComment
	case tt == chroma.NameFunction || tt == chroma.NameBuiltin:
		return ColorFunction
	case tt == chroma.NameClass || tt == chroma.NameBuiltinPseudo:
		return ColorType
	case tt == chroma.Number || tt == chroma.NumberInteger || tt == chroma.NumberFloat:
		return ColorNumber
	case tt == chroma.Operator || tt == chroma.Punctuation:
		return ColorOperator
	default:
		return ColorDefault
	}
}

// HighlightLine syntax-highlights a line of Go code and applies diff background.
func (h *Highlighter) HighlightLine(line string, lineType LineType, width int) string {
	iterator, err := h.lexer.Tokenise(nil, line)
	if err != nil {
		return h.applyBackground(line, lineType, width)
	}

	var result strings.Builder
	tokens := iterator.Tokens()

	for _, token := range tokens {
		color := h.tokenColor(token.Type)
		style := lipgloss.NewStyle().Foreground(color)

		switch lineType {
		case LineAdded:
			style = style.Background(AddedBg)
		case LineRemoved:
			style = style.Background(RemovedBg)
		}

		result.WriteString(style.Render(token.Value))
	}

	rendered := result.String()

	if width > 0 {
		visibleLen := visibleLength(line)
		if visibleLen < width {
			padding := strings.Repeat(" ", width-visibleLen)
			var bgStyle lipgloss.Style
			switch lineType {
			case LineAdded:
				bgStyle = lipgloss.NewStyle().Background(AddedBg)
			case LineRemoved:
				bgStyle = lipgloss.NewStyle().Background(RemovedBg)
			default:
				bgStyle = lipgloss.NewStyle()
			}
			rendered += bgStyle.Render(padding)
		}
	}

	return rendered
}

func (h *Highlighter) applyBackground(text string, lineType LineType, width int) string {
	var style lipgloss.Style
	switch lineType {
	case LineAdded:
		style = lipgloss.NewStyle().Background(AddedBg).Foreground(AddedFg)
	case LineRemoved:
		style = lipgloss.NewStyle().Background(RemovedBg).Foreground(RemovedFg)
	default:
		style = lipgloss.NewStyle()
	}

	if width > 0 && len(text) < width {
		text += strings.Repeat(" ", width-len(text))
	}

	return style.Render(text)
}

func visibleLength(s string) int {
	count := 0
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		count++
	}
	return count
}

// HighlightHunkHeader styles a hunk header (@@ ... @@).
func (h *Highlighter) HighlightHunkHeader(header string) string {
	style := lipgloss.NewStyle().Foreground(ColorHunk).Bold(true)
	return style.Render(header)
}
