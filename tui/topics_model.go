package tui

import (
	"log"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/troygilman/actormq/client"
	"github.com/troygilman/tui/common"
)

func NewTopicsModel(engine *actor.Engine, client *actor.PID) TopicsModel {
	return TopicsModel{
		engine:  engine,
		client:  client,
		adapter: common.NewAdapter(engine, common.BasicAdapterFunc),
	}
}

type TopicsModel struct {
	engine  *actor.Engine
	client  *actor.PID
	adapter common.Adapter
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
	msg, adapterCmd := model.adapter.Poll(msg)
	switch msg := msg.(type) {
	case client.CreateConsumerResult:
		log.Println("TOPIC - NEW CONSUMER", msg)
	}

	return model, tea.Batch(
		adapterCmd,
	)
}

func (model TopicsModel) View() string {
	return ""
}
