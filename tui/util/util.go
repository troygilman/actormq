package util

import (
	"log"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func MessageCmd(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func SetTableFocus(model table.Model, focus bool) table.Model {
	if focus && !model.Focused() {
		log.Println("Setting focus = true")
		model.Focus()
	} else if !focus && model.Focused() {
		log.Println("Setting focu = false")
		model.Blur()
	}
	return model
}
