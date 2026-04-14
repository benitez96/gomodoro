package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/benitez96/gomodoro/internal/domain"
	"github.com/benitez96/gomodoro/internal/sound"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		if m.timer == nil {
			return m, nil
		}
		return m.handleTick()

	case tea.KeyMsg:
		switch m.screen {
		case screenSetup:
			return m.updateSetup(msg)
		case screenRunning:
			return m.updateRunning(msg)
		}
	}

	return m, nil
}

func (m Model) handleTick() (Model, tea.Cmd) {
	prevState := m.timer.State
	finished := m.timer.Tick()

	if !finished {
		if m.timer.IsRunning() {
			return m, tick()
		}
		return m, nil
	}

	switch prevState {
	case domain.StateFocus:
		m.history.Add(domain.SessionFocus, m.timer.FocusDuration)
	case domain.StateBreak:
		m.history.Add(domain.SessionBreak, m.timer.BreakDuration)
	}

	return m, bellCmd()
}

func bellCmd() tea.Cmd {
	return func() tea.Msg {
		sound.Bell()
		return nil
	}
}

func (m Model) updateSetup(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "left", "h", "tab", "shift+tab":
		m.cursor = 1 - m.cursor

	case "right", "l":
		m.cursor = 1 - m.cursor

	case "up", "k":
		if m.cursor == 0 {
			m.focusMins++
		} else {
			m.breakMins++
		}

	case "down", "j":
		if m.cursor == 0 && m.focusMins > 1 {
			m.focusMins--
		} else if m.cursor == 1 && m.breakMins > 1 {
			m.breakMins--
		}

	case "enter", " ":
		m.timer = domain.New(
			time.Duration(m.focusMins)*time.Minute,
			time.Duration(m.breakMins)*time.Minute,
		)
		m.timer.StartFocus()
		m.screen = screenRunning
		return m, tick()
	}

	return m, nil
}

func (m Model) updateRunning(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "r":
		m.screen = screenSetup
		m.timer = nil
		return m, nil

	case " ":
		if m.timer.State == domain.StatePaused {
			m.timer.Resume()
			return m, tick()
		}
		m.timer.Pause()

	case "b":
		if m.timer.State == domain.StateIdle {
			m.timer.StartBreak()
			return m, tick()
		}

	case "f":
		if m.timer.State == domain.StateIdle {
			m.timer.StartFocus()
			return m, tick()
		}
	}

	return m, nil
}
