package ui

import tea "github.com/charmbracelet/bubbletea"

type keyAction int

const (
	keyNone keyAction = iota
	keyUp
	keyDown
	keyTop
	keyBottom
	keySelect
	keyPaste
	keyDelete
	keyClearAll
	keyQuit
)

func parseKey(msg tea.KeyMsg) keyAction {
	switch msg.String() {
	case "k", "up":
		return keyUp
	case "j", "down":
		return keyDown
	case "g", "home":
		return keyTop
	case "G", "end":
		return keyBottom
	case "enter":
		return keySelect
	case "p":
		return keyPaste
	case "d":
		return keyDelete
	case "x":
		return keyClearAll
	case "q", "esc", "ctrl+c":
		return keyQuit
	}
	return keyNone
}
