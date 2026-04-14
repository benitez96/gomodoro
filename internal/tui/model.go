package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/benitez96/gomodoro/internal/domain"
)

type screen int

const (
	screenSetup screen = iota
	screenRunning
)

type tickMsg time.Time

type Model struct {
	screen    screen
	timer     *domain.Timer
	history   *domain.History
	focusMins int
	breakMins int
	cursor    int
	width     int
	height    int
}

func New() Model {
	return Model{
		screen:    screenSetup,
		history:   &domain.History{},
		focusMins: 25,
		breakMins: 5,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
