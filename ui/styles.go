package ui

import "github.com/charmbracelet/lipgloss"

// Tokyo Night palette
var (
	background = lipgloss.Color("#1a1b26")
	surface    = lipgloss.Color("#24283b")
	overlay    = lipgloss.Color("#414868")
	foreground = lipgloss.Color("#c0caf5")
	muted      = lipgloss.Color("#565f89")
	blue       = lipgloss.Color("#7aa2f7")
	cyan       = lipgloss.Color("#7dcfff")
	green      = lipgloss.Color("#9ece6a")
	red        = lipgloss.Color("#f7768e")
)

var (
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Foreground(blue).
			Bold(true)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(overlay).
			Padding(0, 1)

	ItemStyle = lipgloss.NewStyle().
			Foreground(foreground).
			Padding(0, 1)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(blue).
			Background(surface).
			Bold(true).
			Padding(0, 1)

	IndexStyle = lipgloss.NewStyle().
			Foreground(muted).
			Width(4).
			Align(lipgloss.Right)

	SelectedIndexStyle = lipgloss.NewStyle().
				Foreground(blue).
				Background(surface).
				Width(4).
				Align(lipgloss.Right)

	TimestampStyle = lipgloss.NewStyle().
			Foreground(cyan)

	SelectedTimestampStyle = lipgloss.NewStyle().
				Foreground(cyan).
				Background(surface)

	CursorStyle = lipgloss.NewStyle().
			Foreground(blue).
			Bold(true)

	HelpBarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(overlay).
			Foreground(muted).
			Padding(0, 1)

	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(foreground).
			Bold(true)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(muted)

	EmptyStyle = lipgloss.NewStyle().
			Foreground(muted).
			Italic(true).
			Padding(1, 2)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	DangerStyle = lipgloss.NewStyle().
			Foreground(red).
			Bold(true)
)
