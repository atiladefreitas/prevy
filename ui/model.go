package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/atiladefreitas/prevy/clipboard"
	"github.com/atiladefreitas/prevy/store"
)

type status int

const (
	statusBrowsing status = iota
	statusCopied
	statusCleared
)

type Model struct {
	entries []store.Entry
	cursor  int
	width   int
	height  int
	status  status
}

func New() Model {
	entries, _ := store.Load()

	// read current clipboard and add if new
	current, err := clipboard.Read()
	if err == nil && current != "" {
		entries = store.Add(entries, current)
		_ = store.Save(entries)
	}

	return Model{
		entries: entries,
		cursor:  0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		action := parseKey(msg)

		switch action {
		case keyQuit:
			return m, tea.Quit

		case keyUp:
			if m.cursor > 0 {
				m.cursor--
			}

		case keyDown:
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}

		case keySelect:
			if len(m.entries) > 0 {
				_ = clipboard.Write(m.entries[m.cursor].Content)
				m.status = statusCopied
				return m, tea.Quit
			}

		case keyClearAll:
			_ = store.Clear()
			m.entries = []store.Entry{}
			m.cursor = 0
			m.status = statusCleared
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.status == statusCopied {
		msg := SuccessStyle.Render("  Copied to clipboard!")
		return AppStyle.Render(msg) + "\n"
	}

	width := m.width
	if width == 0 {
		width = 60
	}

	// inner content width (accounting for border + padding)
	innerWidth := width - 8
	if innerWidth < 30 {
		innerWidth = 30
	}

	// build list
	var listContent string

	if len(m.entries) == 0 {
		empty := EmptyStyle.Render("No clipboard history yet. Copy something!")
		listContent = empty
	} else {
		var rows []string
		// calculate how many rows fit (reserve space for borders, padding, help bar)
		maxVisible := m.height - 10
		if maxVisible < 3 {
			maxVisible = 10
		}

		// scrolling window
		start := 0
		if m.cursor >= maxVisible {
			start = m.cursor - maxVisible + 1
		}
		end := start + maxVisible
		if end > len(m.entries) {
			end = len(m.entries)
		}

		for i := start; i < end; i++ {
			rows = append(rows, m.renderRow(i, innerWidth))
		}
		listContent = strings.Join(rows, "\n")
	}

	// main box
	mainBox := BorderStyle.
		Width(innerWidth).
		Render(
			TitleStyle.Render(" Prevy") + "\n\n" +
				listContent + "\n",
		)

	// help bar
	help := m.renderHelp(innerWidth)

	return AppStyle.Render(mainBox + "\n" + help)
}

func (m Model) renderRow(index int, maxWidth int) string {
	e := m.entries[index]
	isSelected := index == m.cursor

	// index number
	idxStr := fmt.Sprintf("%d", index+1)

	// relative time
	timeStr := relativeTime(e.Timestamp)

	// content preview -- truncate to fit
	// available width = maxWidth - index(4) - cursor(2) - time(~8) - spacing(6)
	contentWidth := maxWidth - 4 - 2 - len(timeStr) - 6
	if contentWidth < 10 {
		contentWidth = 10
	}
	preview := truncate(singleLine(e.Content), contentWidth)

	if isSelected {
		cursor := CursorStyle.Render(">")
		idx := SelectedIndexStyle.Render(idxStr)
		content := SelectedStyle.Width(contentWidth).Render(preview)
		ts := SelectedTimestampStyle.Render(timeStr)
		return lipgloss.JoinHorizontal(lipgloss.Top, cursor, " ", idx, " ", content, " ", ts)
	}

	cursor := "  "
	idx := IndexStyle.Render(idxStr)
	content := ItemStyle.Width(contentWidth).Render(preview)
	ts := TimestampStyle.Render(timeStr)
	return lipgloss.JoinHorizontal(lipgloss.Top, cursor, idx, " ", content, " ", ts)
}

func (m Model) renderHelp(width int) string {
	pairs := []struct{ key, desc string }{
		{"enter", "copy"},
		{"x", "clear all"},
		{"q", "quit"},
	}

	var parts []string
	for _, p := range pairs {
		k := HelpKeyStyle.Render(p.key)
		d := HelpDescStyle.Render(" " + p.desc)
		parts = append(parts, k+d)
	}

	helpText := strings.Join(parts, HelpDescStyle.Render("  "))

	return HelpBarStyle.
		Width(width).
		Align(lipgloss.Center).
		Render(helpText)
}

func singleLine(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\t", " ")
	return s
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

func relativeTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		s := int(d.Seconds())
		if s < 1 {
			s = 1
		}
		return fmt.Sprintf("%ds ago", s)
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	}
}
