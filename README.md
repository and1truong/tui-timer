# tui-timer

A terminal Pomodoro timer TUI built with Go, Bubbletea, and Lipgloss.

## Install

```bash
go install github.com/and1truong/tui-timer/cmd/tui-timer@latest
```

Or build from source:

```bash
git clone https://github.com/and1truong/tui-timer.git
cd tui-timer
make build
```

## Usage

```bash
tui-timer
tui-timer --work 50m
tui-timer --work 50m --voice Alex
tui-timer --short-break 10m --long-break 20m
```

## Keybindings

| Key            | Action           |
|----------------|-----------------|
| `space`        | Start/Pause      |
| `r`            | Reset            |
| `s`            | Skip             |
| `c`            | Open config in $EDITOR |
| `shift+↑`      | +1 minute        |
| `shift+↓`      | -1 minute        |
| `shift+→`      | +10 minutes      |
| `shift+←`      | -10 minutes      |
| `q`            | Quit             |

## Config

Path: `~/.config/tui-timer/config.yaml`

Auto-created on first run with defaults:

```yaml
work_duration: 25m
short_break: 5m
long_break: 15m
cycles_before_long: 4

sounds:
  tick: true
  finish: true
  break: true

voice:
  enabled: true
  voice: "Samantha"
  messages:
    work_done: "Work session finished"
    break_done: "Break finished"
    start: "Focus time started"
```

## CLI Flags

| Flag | Description | Example |
|------|-------------|---------|
| `--work` | Work session duration | `--work 50m` |
| `--short-break` | Short break duration | `--short-break 10m` |
| `--long-break` | Long break duration | `--long-break 20m` |
| `--voice` | macOS voice name | `--voice Alex` |

CLI flags override config file values.

## macOS Voices

List available voices:

```bash
say -v '?'
```

Common voices: Samantha, Alex, Victoria, Daniel, Karen, Moira, Tessa.

## Logging

Session logs are written to `~/.local/share/tui-timer/log.txt`.

## Project Structure

```
cmd/tui-timer/main.go      — Entry point
internal/config/config.go  — YAML config + CLI flags
internal/timer/engine.go   — Timer state machine
internal/sound/sound.go    — Sound interface + macOS impl
internal/ui/model.go       — Bubbletea model
internal/ui/view.go        — Lipgloss rendering
internal/ui/keys.go        — Keybindings
internal/logger/logger.go  — File logger
```
