package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/and1truong/tui-timer/internal/config"
	"github.com/and1truong/tui-timer/internal/logger"
	"github.com/and1truong/tui-timer/internal/sound"
	"github.com/and1truong/tui-timer/internal/ui"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.ApplyCLIFlags(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "flags: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	player := sound.NewMacPlayer()

	model := ui.NewModel(cfg, player, log)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
