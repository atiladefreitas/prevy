package ui

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/atiladefreitas/prevy/clipboard"
	"github.com/atiladefreitas/prevy/daemon"
	"github.com/atiladefreitas/prevy/store"
)

// ── Custom messages ──────────────────────────────────────────────────

// entriesLoadedMsg is sent when history finishes loading from disk.
type entriesLoadedMsg struct {
	entries []store.Entry
}

// tickMsg drives periodic updates (timestamp refresh, flash dismiss).
type tickMsg time.Time

// flashMsg sets a temporary status message.
type flashMsg struct {
	text  string
	style lipgloss.Style
}

// clearFlashMsg clears the flash after a timeout.
type clearFlashMsg struct{}

// ── Commands ─────────────────────────────────────────────────────────

func loadEntries() tea.Msg {
	entries, _ := store.Load()
	return entriesLoadedMsg{entries: entries}
}

func tickEvery() tea.Cmd {
	return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func dismissFlashAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(_ time.Time) tea.Msg {
		return clearFlashMsg{}
	})
}

// ── Model ────────────────────────────────────────────────────────────

// Model holds all TUI state.
type Model struct {
	entries      []store.Entry
	cursor       int
	width        int
	height       int
	ready        bool
	pasteContent string
	shouldPaste  bool
	shouldQuit   bool
	flashText    string
	flashStyle   lipgloss.Style
	daemonAlive  bool
}

// New creates a Model. Actual data loading happens in Init via a command.
func New() Model {
	return Model{
		daemonAlive: daemon.IsRunning(),
	}
}

func (m Model) ShouldPaste() bool    { return m.shouldPaste }
func (m Model) PasteContent() string { return m.pasteContent }

// Init kicks off async entry loading and the periodic tick.
func (m Model) Init() tea.Cmd {
	return tea.Batch(loadEntries, tickEvery())
}

// ── Update ───────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case entriesLoadedMsg:
		m.entries = msg.entries
		m.ready = true
		return m, nil

	case tickMsg:
		// Refresh timestamps on screen; just re-render.
		return m, tickEvery()

	case clearFlashMsg:
		m.flashText = ""
		return m, nil

	case flashMsg:
		m.flashText = msg.text
		m.flashStyle = msg.style
		return m, dismissFlashAfter(2 * time.Second)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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

	case keyTop:
		m.cursor = 0

	case keyBottom:
		if len(m.entries) > 0 {
			m.cursor = len(m.entries) - 1
		}

	case keySelect:
		if len(m.entries) > 0 {
			_ = clipboard.Write(m.entries[m.cursor].Content)
			m.shouldQuit = true
			return m, tea.Sequence(
				func() tea.Msg {
					return flashMsg{text: " Copied to clipboard!", style: successStyle}
				},
				tea.Quit,
			)
		}

	case keyPaste:
		if len(m.entries) > 0 {
			_ = clipboard.Write(m.entries[m.cursor].Content)
			m.pasteContent = m.entries[m.cursor].Content
			m.shouldPaste = true
			return m, tea.Quit
		}

	case keyDelete:
		if len(m.entries) > 0 {
			m.entries = store.Delete(m.entries, m.cursor)
			_ = store.Save(m.entries)
			if m.cursor >= len(m.entries) && m.cursor > 0 {
				m.cursor--
			}
			return m, func() tea.Msg {
				return flashMsg{text: " Entry deleted", style: dangerStyle}
			}
		}

	case keyClearAll:
		_ = store.Clear()
		m.entries = []store.Entry{}
		m.cursor = 0
		return m, func() tea.Msg {
			return flashMsg{text: " History cleared", style: dangerStyle}
		}
	}

	return m, nil
}

// ── View ─────────────────────────────────────────────────────────────

func (m Model) View() string {
	termW := m.width
	termH := m.height
	if termW == 0 {
		termW = 80
	}
	if termH == 0 {
		termH = 24
	}

	if !m.ready {
		loading := lipgloss.NewStyle().Foreground(muted).Render("Loading...")
		return lipgloss.Place(termW, termH, lipgloss.Center, lipgloss.Center, loading)
	}

	// Cap inner width for readability.
	const maxInner = 78
	innerW := termW - 6
	if innerW > maxInner {
		innerW = maxInner
	}
	if innerW < 40 {
		innerW = 40
	}

	// ── Sections ─────────────────────────────────────────────────

	header := m.viewHeader(innerW)
	div := dividerStyle.Render(strings.Repeat("─", innerW-4))

	// How many list rows fit.
	listH := termH - 18 // borders + header + dividers + preview + help + padding
	if listH < 3 {
		listH = 3
	}

	var list string
	if len(m.entries) == 0 {
		list = m.viewEmpty(innerW)
	} else {
		list = m.viewList(innerW, listH)
	}

	// Preview (only when entries exist).
	var preview string
	if len(m.entries) > 0 {
		preview = "\n" + div + "\n" + m.viewPreview(innerW)
	}

	inner := header + "\n" + div + "\n" + list + preview

	box := outerBorderStyle.Width(innerW).Render(inner)
	help := m.viewHelp(innerW)
	full := box + "\n" + help

	return lipgloss.Place(termW, termH, lipgloss.Center, lipgloss.Center, full)
}

// ── Header ───────────────────────────────────────────────────────────

func (m Model) viewHeader(w int) string {
	left := logoStyle.Render("  Prevy") + headerAccentStyle.Render(" clipboard")

	// Right side: flash message OR daemon + count.
	var right string
	if m.flashText != "" {
		right = m.flashStyle.Render(m.flashText)
	} else {
		var badge string
		if m.daemonAlive {
			badge = daemonOnStyle.Render("  ") + countStyle.Render("on")
		} else {
			badge = daemonOffStyle.Render("  ") + countStyle.Render("off")
		}
		items := countStyle.Render(fmt.Sprintf("  %d items", len(m.entries)))
		right = badge + "  " + items
	}

	gap := w - 4 - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}
	return left + strings.Repeat(" ", gap) + right
}

// ── List ─────────────────────────────────────────────────────────────

func (m Model) viewList(w, maxRows int) string {
	total := len(m.entries)
	if maxRows > total {
		maxRows = total
	}

	// Scroll window.
	start := 0
	if m.cursor >= maxRows {
		start = m.cursor - maxRows + 1
	}
	end := start + maxRows
	if end > total {
		end = total
	}

	scrollbar := buildScrollbar(maxRows, start, total)

	// Available width for row content (minus scrollbar column + gap).
	rowW := w - 6

	var rows []string
	for i := start; i < end; i++ {
		row := m.viewRow(i, rowW)
		rows = append(rows, row+"  "+scrollbar[i-start])
	}

	return strings.Join(rows, "\n")
}

func (m Model) viewRow(idx, w int) string {
	e := m.entries[idx]
	sel := idx == m.cursor

	// Fixed-width columns: index(3) + gap(1) + content(flex) + gap(1) + time(7)
	// The cursor prefix is 2 chars: "> " or "  ".
	idxW := 3
	tsW := 7
	fixedW := 2 + idxW + 1 + 1 + tsW
	contentW := w - fixedW
	if contentW < 8 {
		contentW = 8
	}

	// Format index right-aligned.
	idxStr := fmt.Sprintf("%*d", idxW, idx+1)

	// Time string, padded to fixed width.
	ts := relativeTime(e.Timestamp)
	if runeWidth(ts) < tsW {
		ts = strings.Repeat(" ", tsW-runeWidth(ts)) + ts
	}

	// Content: force single line, then truncate to exact rune width.
	content := singleLine(e.Content)
	content = runesTruncate(content, contentW)
	// Pad content to fixed width so timestamp aligns.
	if runeWidth(content) < contentW {
		content = content + strings.Repeat(" ", contentW-runeWidth(content))
	}

	if sel {
		return cursorStyle.Render("> ") +
			selectedIndexStyle.Render(idxStr) + " " +
			selectedContentStyle.Render(content) + " " +
			selectedTimestampStyle.Render(ts)
	}

	return "  " +
		indexStyle.Render(idxStr) + " " +
		normalContentStyle.Render(content) + " " +
		timestampStyle.Render(ts)
}

// ── Scrollbar ────────────────────────────────────────────────────────

func buildScrollbar(viewH, start, total int) []string {
	chars := make([]string, viewH)

	if total <= viewH {
		for i := range chars {
			chars[i] = " "
		}
		return chars
	}

	thumbSz := viewH * viewH / total
	if thumbSz < 1 {
		thumbSz = 1
	}

	thumbPos := 0
	if total-viewH > 0 {
		thumbPos = start * (viewH - thumbSz) / (total - viewH)
	}

	for i := range chars {
		if i >= thumbPos && i < thumbPos+thumbSz {
			chars[i] = scrollThumbStyle.Render("┃")
		} else {
			chars[i] = scrollTrackStyle.Render("│")
		}
	}
	return chars
}

// ── Preview ──────────────────────────────────────────────────────────

func (m Model) viewPreview(w int) string {
	if m.cursor < 0 || m.cursor >= len(m.entries) {
		return ""
	}

	raw := m.entries[m.cursor].Content
	pw := w - 6
	if pw < 10 {
		pw = 10
	}

	lines := strings.Split(raw, "\n")
	maxLines := 3
	overflow := len(lines) > maxLines
	if overflow {
		lines = lines[:maxLines]
	}

	var out []string
	for _, line := range lines {
		line = strings.ReplaceAll(line, "\t", "  ")
		line = runesTruncate(line, pw)
		out = append(out, previewContentStyle.Render(line))
	}
	if overflow {
		extra := len(strings.Split(raw, "\n")) - maxLines
		out = append(out, previewLabelStyle.Render(fmt.Sprintf("  ... +%d more lines", extra)))
	}

	label := previewLabelStyle.Render("  Preview")
	return label + "\n" + strings.Join(out, "\n")
}

// ── Empty state ──────────────────────────────────────────────────────

func (m Model) viewEmpty(w int) string {
	art := emptyIconStyle.Render(
		"    ┌─────────┐\n" +
			"    │         │\n" +
			"    │  empty  │\n" +
			"    │         │\n" +
			"    └─────────┘")

	msg := emptyStyle.Width(w - 4).Render(
		art + "\n\n" +
			"No clipboard history yet.\n" +
			"Copy something to get started!")

	return "\n" + msg + "\n"
}

// ── Help bar ─────────────────────────────────────────────────────────

func (m Model) viewHelp(w int) string {
	type bind struct{ key, desc string }
	binds := []bind{
		{"enter", "copy"},
		{"p", "paste"},
		{"d", "delete"},
		{"x", "clear"},
		{"g/G", "top/end"},
		{"q", "quit"},
	}

	sep := helpSepStyle.Render(" | ")

	var parts []string
	for _, b := range binds {
		parts = append(parts,
			helpKeyStyle.Render(b.key)+helpDescStyle.Render(" "+b.desc),
		)
	}

	text := strings.Join(parts, sep)
	return helpBarStyle.Width(w).Align(lipgloss.Center).Render(text)
}

// ── String helpers ───────────────────────────────────────────────────

// singleLine collapses all whitespace into spaces.
func singleLine(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\t", " ")
	return s
}

// runesTruncate truncates s so its display width is at most max runes,
// appending "..." when truncated. This guarantees a single-line result.
func runesTruncate(s string, max int) string {
	if max <= 0 {
		return ""
	}
	count := utf8.RuneCountInString(s)
	if count <= max {
		return s
	}
	if max <= 3 {
		runes := []rune(s)
		return string(runes[:max])
	}
	runes := []rune(s)
	return string(runes[:max-3]) + "..."
}

// runeWidth returns the rune count of s.
func runeWidth(s string) int {
	return utf8.RuneCountInString(s)
}

// relativeTime formats a timestamp as a short human-readable string.
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
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}
