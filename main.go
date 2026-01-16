package main

import (
	"fmt"
	"os"

	"grua/internal/git"
	"grua/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Find git repository root
	repoPath, err := git.GetRepoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: not a git repository (or any of the parent directories)")
		os.Exit(1)
	}

	// Create and run the TUI
	model := tui.NewModel(repoPath)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running grua: %v\n", err)
		os.Exit(1)
	}
}
