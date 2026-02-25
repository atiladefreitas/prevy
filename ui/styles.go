package ui

import "github.com/charmbracelet/lipgloss"

// Tokyo Night palette
var (
	surface    = lipgloss.Color("#24283b")
	overlay    = lipgloss.Color("#414868")
	foreground = lipgloss.Color("#c0caf5")
	muted      = lipgloss.Color("#565f89")
	subtle     = lipgloss.Color("#3b4261")
	blue       = lipgloss.Color("#7aa2f7")
	cyan       = lipgloss.Color("#7dcfff")
	green      = lipgloss.Color("#9ece6a")
	red        = lipgloss.Color("#f7768e")
	magenta    = lipgloss.Color("#bb9af7")
)

// ── Header ───────────────────────────────────────────────────────────

var logoStyle = lipgloss.NewStyle().
	Foreground(blue).
	Bold(true)

var headerAccentStyle = lipgloss.NewStyle().
	Foreground(magenta).
	Bold(true)

var countStyle = lipgloss.NewStyle().
	Foreground(muted)

var daemonOnStyle = lipgloss.NewStyle().
	Foreground(green).
	Bold(true)

var daemonOffStyle = lipgloss.NewStyle().
	Foreground(red).
	Bold(true)

// ── List rows ────────────────────────────────────────────────────────

var normalContentStyle = lipgloss.NewStyle().
	Foreground(foreground)

var selectedContentStyle = lipgloss.NewStyle().
	Foreground(foreground).
	Bold(true)

var indexStyle = lipgloss.NewStyle().
	Foreground(subtle)

var selectedIndexStyle = lipgloss.NewStyle().
	Foreground(blue).
	Bold(true)

var timestampStyle = lipgloss.NewStyle().
	Foreground(muted)

var selectedTimestampStyle = lipgloss.NewStyle().
	Foreground(cyan)

var cursorStyle = lipgloss.NewStyle().
	Foreground(blue).
	Bold(true)

// ── Scrollbar ────────────────────────────────────────────────────────

var scrollTrackStyle = lipgloss.NewStyle().
	Foreground(subtle)

var scrollThumbStyle = lipgloss.NewStyle().
	Foreground(blue)

// ── Preview pane ─────────────────────────────────────────────────────

var previewLabelStyle = lipgloss.NewStyle().
	Foreground(muted).
	Italic(true)

var previewContentStyle = lipgloss.NewStyle().
	Foreground(foreground)

// ── Main container ───────────────────────────────────────────────────

var outerBorderStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(overlay).
	Padding(1, 2)

// ── Help bar ─────────────────────────────────────────────────────────

var helpBarStyle = lipgloss.NewStyle().
	Foreground(muted).
	Padding(0, 1)

var helpKeyStyle = lipgloss.NewStyle().
	Foreground(blue).
	Bold(true)

var helpDescStyle = lipgloss.NewStyle().
	Foreground(muted)

var helpSepStyle = lipgloss.NewStyle().
	Foreground(subtle)

// ── Empty state ──────────────────────────────────────────────────────

var emptyStyle = lipgloss.NewStyle().
	Foreground(muted).
	Italic(true).
	Align(lipgloss.Center)

var emptyIconStyle = lipgloss.NewStyle().
	Foreground(subtle)

// ── Flash / status messages ──────────────────────────────────────────

var successStyle = lipgloss.NewStyle().
	Foreground(green).
	Bold(true)

var dangerStyle = lipgloss.NewStyle().
	Foreground(red).
	Bold(true)

var dividerStyle = lipgloss.NewStyle().
	Foreground(subtle)
