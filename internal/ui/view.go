package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/and1truong/tui-timer/internal/timer"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Align(lipgloss.Center)

	timerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			Align(lipgloss.Center)

	modeWorkStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("82"))

	modeBreakStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("214"))

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Align(lipgloss.Center)

	progressFullStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("82"))

	progressEmptyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240"))

	stateStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Italic(true)
)

func renderView(e *timer.Engine, width int) string {
	var b strings.Builder

	// Top: mode + cycle
	modeStr := renderMode(e)
	cycleStr := fmt.Sprintf("Cycle: %d", e.Cycle)
	topLine := fmt.Sprintf("%s  |  %s", modeStr, cycleStr)
	b.WriteString(titleStyle.Width(width).Render(topLine))
	b.WriteString("\n\n")

	// State indicator
	stateStr := renderState(e.State)
	b.WriteString(lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(stateStr))
	b.WriteString("\n\n")

	// Center: big timer
	timeStr := formatDuration(e.Remaining)
	b.WriteString(timerStyle.Width(width).Render(timeStr))
	b.WriteString("\n\n")

	// Progress bar
	bar := renderProgressBar(e.Progress(), width-10)
	b.WriteString(lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(bar))
	b.WriteString("\n\n")

	// Bottom: key hints
	hints := "space: start/pause  |  r: reset  |  s: skip  |  c: config  |  q: quit"
	b.WriteString(hintStyle.Width(width).Render(hints))

	return b.String()
}

func renderMode(e *timer.Engine) string {
	label := e.Mode.String()
	switch e.Mode {
	case timer.ModeWork:
		return modeWorkStyle.Render(label)
	default:
		return modeBreakStyle.Render(label)
	}
}

func renderState(s timer.State) string {
	switch s {
	case timer.StateRunning:
		return stateStyle.Render("▶ Running")
	case timer.StatePaused:
		return stateStyle.Render("⏸ Paused")
	default:
		return stateStyle.Render("⏹ Ready")
	}
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", m, s)
}

func renderProgressBar(progress float64, width int) string {
	if width < 5 {
		width = 20
	}
	filled := int(float64(width) * progress)
	empty := width - filled

	bar := progressFullStyle.Render(strings.Repeat("█", filled)) +
		progressEmptyStyle.Render(strings.Repeat("░", empty))

	pct := fmt.Sprintf(" %3.0f%%", progress*100)
	return bar + pct
}
