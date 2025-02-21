package tui

import (
	"github.com/anthdm/hollywood/actor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/troygilman/actormq/tui/util"
)

func NewTabsModel(engine *actor.Engine, client *actor.PID) tea.Model {
	return TabsModel{
		engine: engine,
		client: client,
	}
}

type TabsModel struct {
	engine *actor.Engine
	client *actor.PID
	tabs   []tea.Model
}

func (model TabsModel) Init() tea.Cmd {
	return nil
}

func (model TabsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := util.NewCommandBuilder()

	switch msg := msg.(type) {
	case CreateConsumerMsg:
		consumerModel := NewConsumerModel(model.engine, model.client, msg.Topic)
		model.tabs = append(model.tabs, consumerModel)
		cmds.AddCmd(consumerModel.Init())
	}

	for index, tab := range model.tabs {
		var tabCmd tea.Cmd
		model.tabs[index], tabCmd = tab.Update(msg)
		cmds.AddCmd(tabCmd)
	}

	return model, cmds.Build()
}

func (model TabsModel) View() string {
	if len(model.tabs) > 0 {
		return model.tabs[0].View()
	}
	return ""
}
