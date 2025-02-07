package tui

import (
	"github.com/anthdm/hollywood/actor"
	tea "github.com/charmbracelet/bubbletea"
)

type ExternalMessage struct {
	Msg tea.Msg
}

func PollExternalMessages(msgs <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return ExternalMessage{
			Msg: <-msgs,
		}
	}
}

func SendWithSender(engine *actor.Engine, pid *actor.PID, msg any, sender *actor.PID) tea.Cmd {
	return func() tea.Msg {
		engine.SendWithSender(pid, msg, sender)
		return nil
	}
}
