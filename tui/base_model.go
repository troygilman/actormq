package tui

import (
	"log/slog"
	"time"

	"github.com/anthdm/hollywood/actor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/troygilman/actormq/client"
	"github.com/troygilman/actormq/cluster"
	"github.com/troygilman/actormq/tui/util"
)

func Run() error {
	model, err := NewBaseModel()
	if err != nil {
		return err
	}

	file, err := tea.LogToFile("tmp/tui.log", "")
	if err != nil {
		return err
	}
	defer file.Close()

	program := tea.NewProgram(model, tea.WithAltScreen())

	_, err = program.Run()
	return err
}

func NewBaseModel() (*BaseModel, error) {
	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		return nil, err
	}

	slog.Default().Info("TEST")

	discovery := engine.Spawn(cluster.NewDiscovery(cluster.DiscoveryConfig{
		Topics: []string{},
		Logger: slog.Default(),
	}), "discovery")

	config := cluster.PodConfig{
		Discovery: discovery,
		Logger:    slog.Default(),
	}

	pods := []*actor.PID{
		engine.Spawn(cluster.NewPod(config), "pod"),
		engine.Spawn(cluster.NewPod(config), "pod"),
		engine.Spawn(cluster.NewPod(config), "pod"),
	}

	time.Sleep(time.Second)
	clientPID := engine.Spawn(client.NewClient(client.ClientConfig{Nodes: pods}), "client")

	adapter := util.NewAdapter(engine, util.BasicAdapterFunc)

	return &BaseModel{
		engine:      engine,
		client:      clientPID,
		adapter:     adapter,
		topicsModel: NewTopicsModel(engine, clientPID),
		tabsModel:   NewTabsModel(engine, clientPID),
	}, nil
}

type BaseModel struct {
	engine      *actor.Engine
	client      *actor.PID
	adapter     util.Adapter
	topicsModel tea.Model
	tabsModel   tea.Model
}

func (model BaseModel) Init() tea.Cmd {
	return tea.Batch(
		model.adapter.Init(),
		model.topicsModel.Init(),
		model.tabsModel.Init(),
	)
}

func (model *BaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := util.NewCommandBuilder()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return model, tea.Quit
		}
	}

	var topicsCmd tea.Cmd
	model.topicsModel, topicsCmd = model.topicsModel.Update(msg)
	cmds.AddCmd(topicsCmd)

	var tabsCmd tea.Cmd
	model.tabsModel, tabsCmd = model.tabsModel.Update(msg)
	cmds.AddCmd(tabsCmd)

	adapterMsg, adapterCmd := model.adapter.Message(msg)
	cmds.AddCmd(adapterCmd)
	switch adapterMsg.(type) {
	}

	return model, cmds.Build()
}

func (model BaseModel) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, model.topicsModel.View(), model.tabsModel.View())
}
