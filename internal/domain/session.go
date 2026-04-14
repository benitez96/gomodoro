package domain

import "time"

type SessionType int

const (
	SessionFocus SessionType = iota
	SessionBreak
)

func (st SessionType) String() string {
	if st == SessionFocus {
		return "Focus"
	}
	return "Break"
}

type Session struct {
	Type        SessionType
	Duration    time.Duration
	CompletedAt time.Time
}

type History struct {
	Sessions []Session
}

func (h *History) Add(t SessionType, d time.Duration) {
	h.Sessions = append(h.Sessions, Session{
		Type:        t,
		Duration:    d,
		CompletedAt: time.Now().Local(),
	})
}
