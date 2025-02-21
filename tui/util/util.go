package util

import tea "github.com/charmbracelet/bubbletea"

func MessageCmd(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
