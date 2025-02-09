package util

import tea "github.com/charmbracelet/bubbletea"

type CommandBuilder struct {
	cmds []tea.Cmd
}

func NewCommandBuilder() CommandBuilder {
	return CommandBuilder{
		cmds: []tea.Cmd{},
	}
}

func (builder *CommandBuilder) AddCmd(cmd tea.Cmd) {
	builder.cmds = append(builder.cmds, cmd)
}

func (builder CommandBuilder) Build() tea.Cmd {
	return tea.Batch(builder.cmds...)
}
