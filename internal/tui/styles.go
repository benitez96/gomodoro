package tui

import "github.com/charmbracelet/lipgloss"

var (
	colorFocus   = lipgloss.Color("#E06C75")
	colorBreak   = lipgloss.Color("#98C379")
	colorMuted   = lipgloss.Color("#5C6370")
	colorAccent  = lipgloss.Color("#61AFEF")
	colorSurface = lipgloss.Color("#282C34")

	styleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent).
			MarginBottom(1)

	styleTimerFocus = lipgloss.NewStyle().
			Foreground(colorFocus).
			Bold(true)

	styleTimerBreak = lipgloss.NewStyle().
			Foreground(colorBreak).
			Bold(true)

	styleLabel = lipgloss.NewStyle().
			Foreground(colorMuted)

	styleSelected = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	styleHistoryFocus = lipgloss.NewStyle().
				Foreground(colorFocus)

	styleHistoryBreak = lipgloss.NewStyle().
				Foreground(colorBreak)

	styleHelp = lipgloss.NewStyle().
			Foreground(colorMuted).
			MarginTop(1)

	styleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSurface).
			Padding(1, 3)

	styleCardBase = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 4).
			Width(20).
			Align(lipgloss.Center)
)
