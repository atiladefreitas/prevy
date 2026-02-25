package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/atiladefreitas/prevy/daemon"
	"github.com/atiladefreitas/prevy/ui"
)

// version is set at build time via -ldflags "-X main.version=..."
var version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--daemon", "-d":
			if err := daemon.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		case "--status":
			if daemon.IsRunning() {
				fmt.Println("prevy daemon is running")
			} else {
				fmt.Println("prevy daemon is not running")
			}
			return
		case "--version", "-v":
			fmt.Printf("prevy %s\n", version)
			return
		case "--help", "-h":
			fmt.Println("Usage: prevy [flags]")
			fmt.Println()
			fmt.Println("  (no flags)   Open clipboard history TUI")
			fmt.Println("  --daemon     Start background clipboard watcher")
			fmt.Println("  --status     Check if daemon is running")
			fmt.Println("  --version    Show version")
			fmt.Println("  --help       Show this help")
			return
		}
	}

	m := ui.New()
	p := tea.NewProgram(m, tea.WithAltScreen())

	result, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// if the user pressed 'p', print content to stdout after TUI closes
	// (content is also in clipboard from the Update handler)
	if final, ok := result.(ui.Model); ok && final.ShouldPaste() {
		fmt.Println(final.PasteContent())
	}
}
