# Podomoro Timer

A production-grade terminal Pomodoro timer TUI built with Go, Bubbletea, and Lipgloss.

## Install

```bash
go install github.com/htruong/podomoro/cmd/podomoro@latest
```

Or build from source:

```bash
git clone https://github.com/htruong/podomoro.git
cd podomoro
go build -o podomoro ./cmd/podomoro
```

## Usage

```bash
podomoro
podomoro --work 50m
podomoro --work 50m --voice Alex
podomoro --short-break 10m --long-break 20m
```

## Keybindings

| Key     | Action       |
|---------|-------------|
| `space` | Start/Pause |
| `r`     | Reset       |
| `s`     | Skip        |
| `c`     | Open config in $EDITOR |
| `q`     | Quit        |

## Config

Path: `~/.config/podomoro-timer/config.yaml`

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

Session logs are written to `~/.local/share/podomoro-timer/log.txt`.

## Project Structure

```
cmd/podomoro/main.go       — Entry point
internal/config/config.go  — YAML config + CLI flags
internal/timer/engine.go   — Timer state machine
internal/sound/sound.go    — Sound interface + macOS impl
internal/ui/model.go       — Bubbletea model
internal/ui/view.go        — Lipgloss rendering
internal/ui/keys.go        — Keybindings
internal/logger/logger.go  — File logger
```
