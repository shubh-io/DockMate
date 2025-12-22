package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shubh-io/dockmate/internal/check"
	"github.com/shubh-io/dockmate/internal/config"
	"github.com/shubh-io/dockmate/internal/tui"
	"github.com/shubh-io/dockmate/internal/update"
	"github.com/shubh-io/dockmate/pkg/version"
)

// just a temporary file to indicate a restart is needed
const restartMarkerFile = ".dockmate_restart"

// ============================================================================
// Main
// ============================================================================

func main() {
	// Restart loop for settings changes
	for {
		if !runApp() {
			break
		}
	}
}

func getRestartMarkerPath() string {
	tmpDir := os.TempDir()
	return filepath.Join(tmpDir, restartMarkerFile)
}

func runApp() bool {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "--version", "-v":
			fmt.Printf("DockMate version: %s\n", version.Dockmate_Version)
			return false
		case "update":
			update.UpdateCommand()
			return false
		case "--runtime":
			runtimeSelector := tui.NewRuntimeSelectionModel()
			program := tea.NewProgram(runtimeSelector, tea.WithAltScreen())

			finalModel, err := program.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Runtime selection failed: %v\n", err)
				os.Exit(1)
			}

			rsModel, ok := finalModel.(tui.RuntimeSelectionModel)
			if !ok {
				fmt.Fprintf(os.Stderr, "Invalid model type returned\n")
				os.Exit(1)
			}

			selectedRuntime := strings.TrimSpace(rsModel.GetChoice())
			if selectedRuntime == "" {
				fmt.Fprintf(os.Stderr, "No runtime selected\n")
				os.Exit(1)
			}

			// load current config and update runtime
			cfg, _ := config.Load()
			cfg.Runtime.Type = selectedRuntime

			// Save updated config (if you dont know, config location is ~/.config/dockmate/config.yml or $XDG_CONFIG_HOME/dockmate/config.yml)
			if err := cfg.Save(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to save runtime selection: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Runtime set to %s.\n\n", selectedRuntime)
			fmt.Printf("To run the application: run 'dockmate'\n")
			fmt.Printf("To change runtime interactively later: 'dockmate --runtime'.\n")
			return false
		}
	}

	result := check.RunPreChecks()

	if !result.Passed {
		fmt.Fprintf(os.Stderr, "%s\n\n%s\n", result.ErrorMessage, result.SuggestedAction)
		os.Stderr.Sync()
		os.Exit(1)
	}

	// start the TUI with alternate screen mode
	// (alternate screen = your terminal history stays clean)

	p := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}

	// check if a restart was requested (if temp marker file exists)
	markerPath := getRestartMarkerPath()
	_, err := os.Stat(markerPath)
	if err == nil {
		// Marker file exists then we restarting ;-;
		os.Remove(markerPath) // Clean up the temp file
		return true           // Continue the loop to restart
	}

	// Normal quit
	return false
}
