package tui

import (
	"io"
	"log"
	"log/slog"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/troygilman/actormq/client"
	"github.com/troygilman/actormq/cluster"
	"github.com/troygilman/actormq/tui/util"
)

func Run() error {
	model, err := NewBaseModel()
	if err != nil {
		return err
	}

	file, err := tea.LogToFile("tmp/application.log", "")
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

	discovery := engine.Spawn(cluster.NewDiscovery(), "discovery")

	config := cluster.PodConfig{
		Topics:    []string{"test", "test1", "test2"},
		Discovery: discovery,
		Logger:    slog.New(slog.NewJSONHandler(io.Discard, nil)),
		// Logger: slog.Default(),
	}

	pods := []*actor.PID{
		engine.Spawn(cluster.NewPod(config), "pod"),
		engine.Spawn(cluster.NewPod(config), "pod"),
		engine.Spawn(cluster.NewPod(config), "pod"),
	}

	client := engine.Spawn(client.NewClient(client.ClientConfig{Nodes: pods}), "client")

	adapter := util.NewAdapter(engine, util.BasicAdapterFunc)

	return &BaseModel{
		engine:      engine,
		client:      client,
		adapter:     adapter,
		topicsModel: NewTopicsModel(engine, client),
	}, nil
}

type BaseModel struct {
	engine      *actor.Engine
	client      *actor.PID
	adapter     util.Adapter
	topicsModel tea.Model
}

func (model BaseModel) Init() tea.Cmd {
	return tea.Batch(
		model.adapter.Init(),
		model.topicsModel.Init(),
		model.adapter.Send(model.client, client.CreateConsumer{
			ConsumerConfig: client.ConsumerConfig{
				Topic:        "test",
				Deserializer: remote.ProtoSerializer{},
			},
		}),
	)
}

func (model *BaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return model, tea.Quit
		}
	}

	adapterMsg, adapterCmd := model.adapter.Poll(msg)
	if adapterMsg != nil {
		switch adapterMsg := adapterMsg.(type) {
		case client.CreateConsumerResult:
			log.Println("BASE - CONSUMER", adapterMsg)
		}
	}

	var topicsCmd tea.Cmd
	model.topicsModel, topicsCmd = model.topicsModel.Update(msg)

	return model, tea.Batch(
		adapterCmd,
		topicsCmd,
	)
}

func (model BaseModel) View() string {
	return model.topicsModel.View()
}
