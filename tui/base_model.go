package tui

import (
	"log"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/troygilman/actormq/client"
	"github.com/troygilman/actormq/tui/util"
)

type Config struct {
	Pods []*actor.PID
}

func Run(config Config) error {
	remoter := remote.New("127.0.0.1:8081", remote.NewConfig())

	engine, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(remoter))
	if err != nil {
		return err
	}

	clientPID := engine.Spawn(client.NewClient(client.ClientConfig{Nodes: config.Pods}), "client")

	file, err := tea.LogToFile("tmp/tui.log", "")
	if err != nil {
		return err
	}
	defer file.Close()

	program := tea.NewProgram(NewBaseModel(engine, clientPID),
		tea.WithAltScreen(),
	)

	_, err = program.Run()
	return err
}

func NewBaseModel(engine *actor.Engine, client *actor.PID) BaseModel {
	return BaseModel{
		engine:          engine,
		client:          client,
		adapter:         util.NewAdapter(engine, util.BasicAdapterFunc),
		topicsModel:     NewTopicsModel(engine, client),
		tabsModel:       NewTabsModel(engine, client),
		focusedOnTopics: true,
	}
}

type BaseModel struct {
	engine          *actor.Engine
	client          *actor.PID
	adapter         util.Adapter
	topicsModel     tea.Model
	tabsModel       tea.Model
	focusedOnTopics bool
}

func (model BaseModel) Init() tea.Cmd {
	return tea.Batch(
		model.adapter.Init(),
		model.topicsModel.Init(),
		model.tabsModel.Init(),
	)
}

func (model BaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := util.NewCommandBuilder()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Println("KeyMsg", msg.String())
		switch msg.String() {
		case "ctrl+c", "q":
			return model, tea.Quit
		case "T":
			model.focusedOnTopics = !model.focusedOnTopics

			model.topicsModel, cmd = model.topicsModel.Update(FocusMsg{Focus: model.focusedOnTopics})
			cmds.AddCmd(cmd)

			model.tabsModel, cmd = model.tabsModel.Update(FocusMsg{Focus: !model.focusedOnTopics})
			cmds.AddCmd(cmd)
		}
	case tea.WindowSizeMsg:
		model.topicsModel, cmd = model.topicsModel.Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height,
		})
		cmds.AddCmd(cmd)

		model.tabsModel, cmd = model.tabsModel.Update(tea.WindowSizeMsg{
			Width:  msg.Width - 38,
			Height: msg.Height - 2,
		})
		cmds.AddCmd(cmd)

		return model, cmds.Build()
	}

	model.topicsModel, cmd = model.topicsModel.Update(msg)
	cmds.AddCmd(cmd)

	model.tabsModel, cmd = model.tabsModel.Update(msg)
	cmds.AddCmd(cmd)

	adapterMsg, cmd := model.adapter.Message(msg)
	cmds.AddCmd(cmd)
	switch adapterMsg.(type) {
	}

	return model, cmds.Build()
}

func (model BaseModel) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		baseStyle.Render(model.topicsModel.View()),
		baseStyle.Render(model.tabsModel.View()),
	)
}
