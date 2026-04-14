package domain

import "time"


type State int

const (
	StateIdle	State = iota
	StateFocus
	StateBreak
	StatePaused
)

func (s State) String() string {
	switch s {
	case StateFocus:
		return "Focus"
	case StateBreak:
		return "Break"
	case StatePaused:
		return "Paused"
	default:
		return "Idle"
	}
}

type Timer struct {
	FocusDuration time.Duration
	BreakDuration time.Duration
	State 				State
	Remaining			time.Duration

	beforePause		State
}

func New(focus, brk time.Duration) *Timer {
	return &Timer{
		FocusDuration: 	focus,
		BreakDuration: 	brk,
		State: 					StateIdle,
	}
}

func (t *Timer) StartFocus() {
	t.State = StateFocus
	t.Remaining = t.FocusDuration
}

func (t *Timer) StartBreak() {
	t.State = StateBreak
	t.Remaining = t.BreakDuration
}

func (t *Timer) Pause() {
	if t.State == StateFocus || t.State == StateBreak {
		t.beforePause = t.State
		t.State = StatePaused
	}
}

func (t *Timer) Resume() {
	if t.State == StatePaused {
		t.State = t.beforePause
	}
}

func (t *Timer) Tick() bool {
	if t.State != StateFocus && t.State != StateBreak {
		return false
	}
	if t.Remaining <= time.Second {
		t.Remaining = 0
		t.State = StateIdle
		return true
	}
	t.Remaining -= time.Second
	return false
}

func (t *Timer) IsRunning() bool {
	return t.State == StateFocus || t.State == StateBreak
}
