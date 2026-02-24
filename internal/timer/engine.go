package timer

import "time"

// Mode represents the current timer mode.
type Mode int

const (
	ModeWork Mode = iota
	ModeShortBreak
	ModeLongBreak
)

func (m Mode) String() string {
	switch m {
	case ModeWork:
		return "Work"
	case ModeShortBreak:
		return "Short Break"
	case ModeLongBreak:
		return "Long Break"
	default:
		return "Unknown"
	}
}

// State represents the timer's running state.
type State int

const (
	StateIdle State = iota
	StateRunning
	StatePaused
)

// Event is emitted when notable things happen.
type Event int

const (
	EventNone Event = iota
	EventTick
	EventWorkDone
	EventBreakDone
	EventStarted
)

// Engine is the core timer logic, decoupled from any TUI or clock.
type Engine struct {
	WorkDuration     time.Duration
	ShortBreak       time.Duration
	LongBreak        time.Duration
	CyclesBeforeLong int

	Mode      Mode
	State     State
	Remaining time.Duration
	Cycle     int // completed work cycles
}

// New creates a new timer engine.
func New(work, shortBreak, longBreak time.Duration, cyclesBeforeLong int) *Engine {
	return &Engine{
		WorkDuration:     work,
		ShortBreak:       shortBreak,
		LongBreak:        longBreak,
		CyclesBeforeLong: cyclesBeforeLong,
		Mode:             ModeWork,
		State:            StateIdle,
		Remaining:        work,
	}
}

// Toggle starts or pauses the timer. Returns EventStarted on first start.
func (e *Engine) Toggle() Event {
	switch e.State {
	case StateIdle:
		e.State = StateRunning
		return EventStarted
	case StateRunning:
		e.State = StatePaused
	case StatePaused:
		e.State = StateRunning
	}
	return EventNone
}

// Reset resets the current session to its full duration.
func (e *Engine) Reset() {
	e.State = StateIdle
	e.Remaining = e.currentDuration()
}

// Skip moves to the next session.
func (e *Engine) Skip() Event {
	return e.advance()
}

// Tick decrements the timer by 1 second. Returns the event that occurred.
func (e *Engine) Tick() Event {
	if e.State != StateRunning {
		return EventNone
	}

	e.Remaining -= time.Second
	if e.Remaining <= 0 {
		return e.advance()
	}
	return EventTick
}

// AdjustTime adds delta to both Remaining and the current mode's duration.
// Remaining is clamped to [1s, currentDuration].
func (e *Engine) AdjustTime(delta time.Duration) {
	switch e.Mode {
	case ModeWork:
		e.WorkDuration += delta
		if e.WorkDuration < time.Minute {
			e.WorkDuration = time.Minute
		}
	case ModeShortBreak:
		e.ShortBreak += delta
		if e.ShortBreak < time.Minute {
			e.ShortBreak = time.Minute
		}
	case ModeLongBreak:
		e.LongBreak += delta
		if e.LongBreak < time.Minute {
			e.LongBreak = time.Minute
		}
	}
	e.Remaining += delta
	if e.Remaining < time.Second {
		e.Remaining = time.Second
	}
	if max := e.currentDuration(); e.Remaining > max {
		e.Remaining = max
	}
}

// Progress returns a value from 0.0 to 1.0.
func (e *Engine) Progress() float64 {
	total := e.currentDuration()
	if total == 0 {
		return 0
	}
	elapsed := total - e.Remaining
	return float64(elapsed) / float64(total)
}

func (e *Engine) currentDuration() time.Duration {
	switch e.Mode {
	case ModeWork:
		return e.WorkDuration
	case ModeShortBreak:
		return e.ShortBreak
	case ModeLongBreak:
		return e.LongBreak
	default:
		return e.WorkDuration
	}
}

func (e *Engine) advance() Event {
	var evt Event

	switch e.Mode {
	case ModeWork:
		e.Cycle++
		evt = EventWorkDone
		if e.CyclesBeforeLong > 0 && e.Cycle%e.CyclesBeforeLong == 0 {
			e.Mode = ModeLongBreak
			e.Remaining = e.LongBreak
		} else {
			e.Mode = ModeShortBreak
			e.Remaining = e.ShortBreak
		}
	case ModeShortBreak, ModeLongBreak:
		evt = EventBreakDone
		e.Mode = ModeWork
		e.Remaining = e.WorkDuration
	}

	e.State = StateIdle
	return evt
}
