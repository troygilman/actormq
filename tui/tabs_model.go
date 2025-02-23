package tui

import (
	"log"

	"github.com/anthdm/hollywood/actor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/troygilman/actormq/tui/util"
)

var (
	tabStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
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
	tabs        []*tabData
	selectedTab int
}

func (model TabsModel) Init() tea.Cmd {
	return nil
}

func (model TabsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := util.NewCommandBuilder()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Println(msg.String())
		switch msg.String() {
		case "ctrl+l":
			model.selectedTab = min(model.selectedTab+1, len(model.tabs)-1)
		case "ctrl+h":
			model.selectedTab = max(model.selectedTab-1, 0)
		}
	}

	switch msg := msg.(type) {
	case CreateConsumerMsg:
		consumerModel := NewConsumerModel(model.engine, model.client, msg.Topic)
		model.tabs = append(model.tabs, &tabData{
			title: msg.Topic + " consumer",
			model: consumerModel,
		})
		cmds.AddCmd(consumerModel.Init())
	case FocusMsg:
		if len(model.tabs) > 0 {
			tab := model.tabs[model.selectedTab]
			tab.model, cmd = tab.model.Update(msg)
			cmds.AddCmd(cmd)
		}
	}

	for _, tab := range model.tabs {
		tab.model, cmd = tab.model.Update(msg)
		cmds.AddCmd(cmd)
	}

	return model, cmds.Build()
}

func (model TabsModel) View() string {
	var view string
	if len(model.tabs) > 0 {
		view = model.tabs[model.selectedTab].model.View()
	}
	return baseStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			model.tabsView(),
			view,
		),
	)
}

func (model TabsModel) tabsView() string {
	tabs := []string{}
	for idx, tab := range model.tabs {
		style := tabStyle
		if idx == model.selectedTab {
			style = style.Foreground(lipgloss.Color("10"))
		}
		tabs = append(tabs, style.Render(tab.title))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, tabs...)
}

type tabData struct {
	title string
	model tea.Model
}
