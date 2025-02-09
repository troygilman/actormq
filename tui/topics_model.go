package tui

import (
	"log"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/troygilman/actormq/client"
	"github.com/troygilman/actormq/tui/util"
)

func NewTopicsModel(engine *actor.Engine, client *actor.PID) TopicsModel {
	return TopicsModel{
		engine:  engine,
		client:  client,
		adapter: util.NewAdapter(engine, util.BasicAdapterFunc),
	}
}

type TopicsModel struct {
	engine  *actor.Engine
	client  *actor.PID
	adapter util.Adapter
}

func (model TopicsModel) Init() tea.Cmd {
	return tea.Batch(
		model.adapter.Init(),
		model.adapter.Send(model.client, client.CreateConsumer{
			ConsumerConfig: client.ConsumerConfig{
				Topic:        "test1",
				Deserializer: remote.ProtoSerializer{},
			},
		}),
	)
}

func (model TopicsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := util.NewCommandBuilder()

	msg, adapterCmd := model.adapter.Message(msg)
	cmds.AddCmd(adapterCmd)
	switch msg := msg.(type) {
	case client.CreateConsumerResult:
		log.Println("TOPIC - NEW CONSUMER", msg)
	}

	return model, cmds.Build()
}

func (model TopicsModel) View() string {
	return ""
}
