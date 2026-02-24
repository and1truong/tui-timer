package ui

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/and1truong/tui-timer/internal/config"
	"github.com/and1truong/tui-timer/internal/logger"
	"github.com/and1truong/tui-timer/internal/sound"
	"github.com/and1truong/tui-timer/internal/timer"
)

type tickMsg time.Time

type Model struct {
	engine *timer.Engine
	keys   keyMap
	cfg    *config.Config
	player sound.Player
	logger *logger.Logger
	width  int
	height int
}

func NewModel(cfg *config.Config, player sound.Player, log *logger.Logger) Model {
	e := timer.New(cfg.WorkDuration, cfg.ShortBreak, cfg.LongBreak, cfg.CyclesBeforeLong)
	return Model{
		engine: e,
		keys:   newKeyMap(),
		cfg:    cfg,
		player: player,
		logger: log,
		width:  60,
		height: 20,
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tickMsg:
		return m.handleTick()
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.keys.Toggle):
		evt := m.engine.Toggle()
		if evt == timer.EventStarted {
			m.playVoiceAsync(m.cfg.Voice.Messages.Start)
			m.log("Started %s session", m.engine.Mode)
		}
		return m, nil

	case key.Matches(msg, m.keys.Reset):
		m.engine.Reset()
		m.log("Reset %s session", m.engine.Mode)
		return m, nil

	case key.Matches(msg, m.keys.Skip):
		evt := m.engine.Skip()
		m.handleEvent(evt)
		m.log("Skipped to %s", m.engine.Mode)
		return m, nil

	case key.Matches(msg, m.keys.Config):
		return m, m.openConfig()
	}

	return m, nil
}

func (m Model) handleTick() (tea.Model, tea.Cmd) {
	evt := m.engine.Tick()

	switch evt {
	case timer.EventTick:
		if m.cfg.Sounds.Tick {
			go m.player.PlayBeep(context.Background())
		}
	default:
		m.handleEvent(evt)
	}

	return m, tickCmd()
}

func (m Model) handleEvent(evt timer.Event) {
	ctx := context.Background()

	switch evt {
	case timer.EventWorkDone:
		m.log("Work session completed (cycle %d)", m.engine.Cycle)
		if m.cfg.Sounds.Finish {
			go m.player.PlayBeep(ctx)
		}
		m.playVoiceAsync(m.cfg.Voice.Messages.WorkDone)

	case timer.EventBreakDone:
		m.log("Break completed, starting work")
		if m.cfg.Sounds.Break {
			go m.player.PlayBeep(ctx)
		}
		m.playVoiceAsync(m.cfg.Voice.Messages.BreakDone)
	}
}

func (m Model) playVoiceAsync(message string) {
	if m.cfg.Voice.Enabled && message != "" {
		go m.player.PlayVoice(context.Background(), m.cfg.Voice.Voice, message)
	}
}

func (m Model) log(format string, args ...any) {
	if m.logger != nil {
		m.logger.Log(format, args...)
	}
}

func (m Model) openConfig() tea.Cmd {
	return func() tea.Msg {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}
		path, err := config.ConfigPath()
		if err != nil {
			return nil
		}
		cmd := exec.Command(editor, path)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
		return nil
	}
}

func (m Model) View() string {
	return "\n" + renderView(m.engine, m.width) + "\n"
}
