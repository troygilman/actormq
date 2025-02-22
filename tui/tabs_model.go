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
	engine      *actor.Engine
	client      *actor.PID
	tabs        []tea.Model
	selectedTab int
}

func (model TabsModel) Init() tea.Cmd {
	return nil
}

func (model TabsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := util.NewCommandBuilder()

	switch msg := msg.(type) {
	case CreateConsumerMsg:
		consumerModel := NewConsumerModel(model.engine, model.client, msg.Topic)
		model.tabs = append(model.tabs, consumerModel)
		cmds.AddCmd(consumerModel.Init())
	case FocusMsg:
		if len(model.tabs) > 0 {
			model.tabs[model.selectedTab], cmd = model.tabs[model.selectedTab].Update(msg)
			cmds.AddCmd(cmd)
		}
	}

	for index, tab := range model.tabs {
		model.tabs[index], cmd = tab.Update(msg)
		cmds.AddCmd(cmd)
	}

	return model, cmds.Build()
}

func (model TabsModel) View() string {
	var view string
	if len(model.tabs) > 0 {
		view = model.tabs[0].View()
	}
	return baseStyle.Render(view)
}
