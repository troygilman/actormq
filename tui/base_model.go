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
)

func Run() error {
	model, err := NewBaseModel()
	if err != nil {
		return err
	}
	program := tea.NewProgram(model) // tea.WithAltScreen(),

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

	msgs := make(chan tea.Msg, 100)
	adapter := engine.SpawnFunc(func(act *actor.Context) {
		switch msg := act.Message().(type) {
		case actor.Initialized, actor.Started, actor.Stopped:
		default:
			msgs <- msg
		}
	}, "adapter")

	return &BaseModel{
		engine:  engine,
		client:  client,
		adapter: adapter,
		msgs:    msgs,
	}, nil
}

type BaseModel struct {
	engine  *actor.Engine
	client  *actor.PID
	adapter *actor.PID
	msgs    chan tea.Msg
}

func (model BaseModel) Init() tea.Cmd {
	return tea.Batch(
		PollExternalMessages(model.msgs),
		model.send(model.client, client.CreateConsumer{
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
	case ExternalMessage:
		switch msg.Msg.(type) {
		case client.CreateConsumerResult:
			log.Println("NEW CONSUMER")
		}
		return model, tea.Batch(
			PollExternalMessages(model.msgs),
		)
	}
	return model, nil
}

func (model BaseModel) View() string {
	return ""
}

func (model BaseModel) send(pid *actor.PID, msg any) tea.Cmd {
	return SendWithSender(model.engine, pid, msg, model.adapter)
}
