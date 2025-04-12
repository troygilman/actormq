package tui

import (
	"github.com/anthdm/hollywood/actor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/troygilman/actormq/tui/util"
)

var (
	tabStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Foreground(lipgloss.Color("255")).
		BorderForeground(lipgloss.Color("255")).
		Padding(0, 1)
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
	window      tea.WindowSizeMsg
}

func (model TabsModel) Init() tea.Cmd {
	return nil
}

func (model TabsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := util.NewCommandBuilder()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ">", "ctrl+right":
			model.selectedTab = min(model.selectedTab+1, len(model.tabs)-1)
		case "<", "ctrl+left":
			model.selectedTab = max(model.selectedTab-1, 0)
		case "D":
			model.tabs = append(model.tabs[:model.selectedTab], model.tabs[model.selectedTab+1:]...)
			model.selectedTab = min(model.selectedTab, len(model.tabs)-1)
			model.selectedTab = max(model.selectedTab, 0)
		}
	case tea.WindowSizeMsg:
		model.window = msg
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
	return lipgloss.NewStyle().
		Width(model.window.Width).
		Height(model.window.Height).
		MaxWidth(model.window.Width).
		MaxHeight(model.window.Height).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().MaxWidth(model.window.Width).Render(model.tabsView()),
				view,
			),
		)
}

func (model TabsModel) tabsView() string {
	tabs := []string{}
	length := 0
	for idx, tab := range model.tabs {
		style := tabStyle
		if idx == model.selectedTab {
			style = style.Foreground(lipgloss.Color("240")).BorderForeground(lipgloss.Color("240"))
		}
		renderedTab := style.Render(tab.title)
		renderedTabWidth := lipgloss.Width(renderedTab)
		if length+renderedTabWidth > model.window.Width {
			break
		}
		tabs = append(tabs, renderedTab)
		length += renderedTabWidth
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, tabs...)
}

type tabData struct {
	title string
	model tea.Model
}
