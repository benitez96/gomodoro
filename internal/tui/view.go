package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/benitez96/gomodoro/internal/domain"
)

var bigDigits = [10][5]string{
	{"███", "█ █", "█ █", "█ █", "███"},
	{" █ ", "██ ", " █ ", " █ ", "███"},
	{"███", "  █", "███", "█  ", "███"},
	{"███", "  █", "███", "  █", "███"},
	{"█ █", "█ █", "███", "  █", "  █"},
	{"███", "█  ", "███", "  █", "███"},
	{"███", "█  ", "███", "█ █", "███"},
	{"███", "  █", "  █", "  █", "  █"},
	{"███", "█ █", "███", "█ █", "███"},
	{"███", "█ █", "███", "  █", "███"},
}

var bigColon = [5]string{" ", "█", " ", "█", " "}

func renderBigTime(d time.Duration, style lipgloss.Style) string {
	mins := int(d.Minutes())
	secs := int(d.Seconds()) % 60

	digits := [4]int{mins / 10, mins % 10, secs / 10, secs % 10}

	rows := make([]string, 5)
	for i := 0; i < 5; i++ {
		rows[i] = fmt.Sprintf("%s  %s   %s   %s  %s",
			bigDigits[digits[0]][i],
			bigDigits[digits[1]][i],
			bigColon[i],
			bigDigits[digits[2]][i],
			bigDigits[digits[3]][i],
		)
	}

	return style.Render(strings.Join(rows, "\n"))
}

func (m Model) View() string {
	switch m.screen {
	case screenSetup:
		return m.viewSetup()
	case screenRunning:
		return m.viewRunning()
	}
	return ""
}

func renderBigMins(mins int, style lipgloss.Style) string {
	d1, d2 := mins/10, mins%10
	rows := make([]string, 5)
	for i := 0; i < 5; i++ {
		rows[i] = bigDigits[d1][i] + "  " + bigDigits[d2][i]
	}
	return style.Render(strings.Join(rows, "\n"))
}

func makeCard(label string, mins int, active bool, color lipgloss.Color) string {
	var textStyle lipgloss.Style
	var borderColor lipgloss.Color

	if active {
		textStyle = lipgloss.NewStyle().Foreground(color).Bold(true)
		borderColor = color
	} else {
		textStyle = lipgloss.NewStyle().Foreground(colorMuted)
		borderColor = colorMuted
	}

	title := textStyle.Render(label)
	bigMins := renderBigMins(mins, textStyle)
	minLabel := textStyle.Render("min")

	inner := lipgloss.JoinVertical(lipgloss.Center, title, "", bigMins, minLabel)
	return styleCardBase.BorderForeground(borderColor).Render(inner)
}

func (m Model) viewSetup() string {
	title := styleTitle.Render("gomodoro")

	focusCard := makeCard("Focus", m.focusMins, m.cursor == 0, colorFocus)
	breakCard := makeCard("Break", m.breakMins, m.cursor == 1, colorBreak)

	cards := lipgloss.JoinHorizontal(lipgloss.Top, focusCard, "  ", breakCard)

	help := styleHelp.Render("← →  select   ↑ ↓  adjust   enter  start   q  quit")

	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		cards,
		help,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) viewRunning() string {
	var stateLabel string
	var timerStyle lipgloss.Style

	switch m.timer.State {
	case domain.StateFocus:
		stateLabel = styleTimerFocus.Render("● Focus")
		timerStyle = styleTimerFocus
	case domain.StateBreak:
		stateLabel = styleTimerBreak.Render("● Break")
		timerStyle = styleTimerBreak
	case domain.StatePaused:
		stateLabel = styleLabel.Render("⏸  Paused")
		timerStyle = styleLabel
	case domain.StateIdle:
		stateLabel = styleLabel.Render("✓ Done")
		timerStyle = styleLabel
	}

	bigTime := renderBigTime(m.timer.Remaining, timerStyle)

	history := m.renderHistory()

	var help string
	switch m.timer.State {
	case domain.StateFocus, domain.StateBreak:
		help = styleHelp.Render("space  pause   r  reset   q  quit")
	case domain.StatePaused:
		help = styleHelp.Render("space  resume   r  reset   q  quit")
	case domain.StateIdle:
		help = styleHelp.Render("f  focus   b  break   r  reset   q  quit")
	}

	content := lipgloss.JoinVertical(lipgloss.Center,
		stateLabel,
		"",
		bigTime,
		"",
		history,
		help,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) renderHistory() string {
	if len(m.history.Sessions) == 0 {
		return styleLabel.Render("no completed sessions")
	}

	var lines []string
	for _, s := range m.history.Sessions {
		mins := int(s.Duration.Minutes())
		clock := s.CompletedAt.Format("15:04")
		var line string
		if s.Type == domain.SessionFocus {
			line = styleHistoryFocus.Render(fmt.Sprintf("✓ Focus  %2d min  %s", mins, clock))
		} else {
			line = styleHistoryBreak.Render(fmt.Sprintf("✓ Break  %2d min  %s", mins, clock))
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
