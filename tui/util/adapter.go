package util

import (
	"github.com/anthdm/hollywood/actor"
	tea "github.com/charmbracelet/bubbletea"
)

type AdapterFunc func(chan<- tea.Msg) actor.ReceiveFunc

func BasicAdapterFunc(msgs chan<- tea.Msg) actor.ReceiveFunc {
	return func(act *actor.Context) {
		switch msg := act.Message().(type) {
		case actor.Initialized, actor.Started, actor.Stopped:
		default:
			msgs <- msg
		}
	}
}

func NewAdapter(engine *actor.Engine, f AdapterFunc) Adapter {
	msgs := make(chan tea.Msg, 100)
	pid := engine.SpawnFunc(f(msgs), "adapter")
	return Adapter{
		engine: engine,
		pid:    pid,
		msgs:   msgs,
	}
}

type Adapter struct {
	engine *actor.Engine
	pid    *actor.PID
	msgs   chan tea.Msg
}

func (adapter Adapter) Init() tea.Cmd {
	return pollMessages(adapter.pid, adapter.msgs)
}

func (adapter Adapter) Message(msg tea.Msg) (tea.Msg, tea.Cmd) {
	switch msg := msg.(type) {
	case Message:
		if msg.pid == adapter.pid {
			return msg.msg, pollMessages(adapter.pid, adapter.msgs)
		}
	}
	return nil, nil
}

func pollMessages(pid *actor.PID, msgs <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		msg := <-msgs
		return Message{
			msg: msg,
			pid: pid,
		}
	}
}

func (adapter Adapter) Send(pid *actor.PID, msg any) tea.Cmd {
	return func() tea.Msg {
		adapter.engine.SendWithSender(pid, msg, adapter.pid)
		return nil
	}
}

func (adapter Adapter) PID() *actor.PID {
	return adapter.pid
}

type Message struct {
	pid *actor.PID
	msg tea.Msg
}
