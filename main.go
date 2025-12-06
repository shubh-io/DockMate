package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// ============================================================================
// Main
// ============================================================================

// start the TUI with alternate screen mode
// (alternate screen = your terminal history stays clean)
func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
