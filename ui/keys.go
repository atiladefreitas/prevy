package ui

import "github.com/charmbracelet/bubbletea"

type keyAction int

const (
	keyNone keyAction = iota
	keyUp
	keyDown
	keySelect
	keyClearAll
	keyQuit
)

func parseKey(msg tea.KeyMsg) keyAction {
	switch msg.String() {
	case "k", "up":
		return keyUp
	case "j", "down":
		return keyDown
	case "enter":
		return keySelect
	case "x":
		return keyClearAll
	case "q", "esc", "ctrl+c":
		return keyQuit
	}
	return keyNone
}
