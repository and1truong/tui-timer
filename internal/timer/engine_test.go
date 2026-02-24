package timer

import (
	"testing"
	"time"
)

func TestNewEngine(t *testing.T) {
	e := New(25*time.Minute, 5*time.Minute, 15*time.Minute, 4)
	if e.Mode != ModeWork {
		t.Errorf("expected ModeWork, got %v", e.Mode)
	}
	if e.State != StateIdle {
		t.Errorf("expected StateIdle, got %v", e.State)
	}
	if e.Remaining != 25*time.Minute {
		t.Errorf("expected 25m remaining, got %v", e.Remaining)
	}
}

func TestToggle(t *testing.T) {
	e := New(25*time.Minute, 5*time.Minute, 15*time.Minute, 4)

	evt := e.Toggle()
	if evt != EventStarted {
		t.Errorf("expected EventStarted, got %v", evt)
	}
	if e.State != StateRunning {
		t.Errorf("expected StateRunning, got %v", e.State)
	}

	e.Toggle()
	if e.State != StatePaused {
		t.Errorf("expected StatePaused, got %v", e.State)
	}

	e.Toggle()
	if e.State != StateRunning {
		t.Errorf("expected StateRunning, got %v", e.State)
	}
}

func TestTickCountdown(t *testing.T) {
	e := New(3*time.Second, 1*time.Second, 1*time.Second, 4)
	e.Toggle() // start

	evt := e.Tick()
	if evt != EventTick {
		t.Errorf("expected EventTick, got %v", evt)
	}
	if e.Remaining != 2*time.Second {
		t.Errorf("expected 2s, got %v", e.Remaining)
	}
}

func TestTickDoesNothingWhenIdle(t *testing.T) {
	e := New(25*time.Minute, 5*time.Minute, 15*time.Minute, 4)
	evt := e.Tick()
	if evt != EventNone {
		t.Errorf("expected EventNone when idle, got %v", evt)
	}
}

func TestWorkToBreakTransition(t *testing.T) {
	e := New(2*time.Second, 5*time.Minute, 15*time.Minute, 4)
	e.Toggle()

	e.Tick() // 1s left
	evt := e.Tick() // 0s -> transition

	if evt != EventWorkDone {
		t.Errorf("expected EventWorkDone, got %v", evt)
	}
	if e.Mode != ModeShortBreak {
		t.Errorf("expected ModeShortBreak, got %v", e.Mode)
	}
	if e.State != StateIdle {
		t.Errorf("expected StateIdle after transition, got %v", e.State)
	}
	if e.Cycle != 1 {
		t.Errorf("expected cycle 1, got %d", e.Cycle)
	}
}

func TestLongBreakAfterCycles(t *testing.T) {
	e := New(1*time.Second, 1*time.Second, 15*time.Minute, 2)

	// Complete 2 work cycles
	for i := 0; i < 2; i++ {
		e.Toggle()
		e.Tick() // work done

		if i < 1 {
			// First cycle -> short break
			e.Toggle()
			e.Tick() // break done
		}
	}

	if e.Mode != ModeLongBreak {
		t.Errorf("expected ModeLongBreak after %d cycles, got %v", e.CyclesBeforeLong, e.Mode)
	}
}

func TestReset(t *testing.T) {
	e := New(25*time.Minute, 5*time.Minute, 15*time.Minute, 4)
	e.Toggle()
	e.Tick()
	e.Reset()

	if e.State != StateIdle {
		t.Errorf("expected StateIdle after reset")
	}
	if e.Remaining != 25*time.Minute {
		t.Errorf("expected full duration after reset, got %v", e.Remaining)
	}
}

func TestSkip(t *testing.T) {
	e := New(25*time.Minute, 5*time.Minute, 15*time.Minute, 4)
	evt := e.Skip()

	if evt != EventWorkDone {
		t.Errorf("expected EventWorkDone on skip, got %v", evt)
	}
	if e.Mode != ModeShortBreak {
		t.Errorf("expected ModeShortBreak after skip, got %v", e.Mode)
	}
}

func TestProgress(t *testing.T) {
	e := New(4*time.Second, 1*time.Second, 1*time.Second, 4)
	if e.Progress() != 0 {
		t.Errorf("expected 0 progress at start")
	}

	e.Toggle()
	e.Tick()
	e.Tick()

	got := e.Progress()
	want := 0.5
	if got != want {
		t.Errorf("expected progress %v, got %v", want, got)
	}
}

func TestBreakToWorkTransition(t *testing.T) {
	e := New(25*time.Minute, 1*time.Second, 15*time.Minute, 4)
	e.Skip() // -> short break

	e.Toggle()
	evt := e.Tick() // break done

	if evt != EventBreakDone {
		t.Errorf("expected EventBreakDone, got %v", evt)
	}
	if e.Mode != ModeWork {
		t.Errorf("expected ModeWork after break, got %v", e.Mode)
	}
}
