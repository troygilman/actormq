package tui

import (
	"log"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/troygilman/actormq/client"
	"github.com/troygilman/actormq/cluster"
	"github.com/troygilman/actormq/tui/util"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func NewTopicsModel(engine *actor.Engine, clientPID *actor.PID) TopicsModel {
	result, err := engine.Request(clientPID, client.CreateProducer{
		ProducerConfig: client.ProducerConfig{
			Topic:      cluster.TopicClusterInternalTopic,
			Serializer: remote.ProtoSerializer{},
		},
	}, time.Second).Result()
	if err != nil {
		panic(err)
	}

	if result, ok := result.(client.CreateProducerResult); ok {
		engine.Send(result.PID, client.ProduceMessage{
			Message: &cluster.NewTopic{
				Name: "test1",
			},
		})
	} else {
		panic("could not create consumer")
	}

	return TopicsModel{
		engine:  engine,
		client:  clientPID,
		adapter: util.NewAdapter(engine, util.BasicAdapterFunc),
		table:   table.New(table.WithWidth(20)),
	}
}

type TopicsModel struct {
	engine  *actor.Engine
	client  *actor.PID
	adapter util.Adapter
	table   table.Model
}

func (model TopicsModel) Init() tea.Cmd {
	return tea.Batch(
		model.adapter.Init(),
		model.adapter.Send(model.client, client.CreateConsumer{
			ConsumerConfig: client.ConsumerConfig{
				Topic:        cluster.TopicClusterInternalTopic,
				PID:          model.adapter.PID(),
				Deserializer: remote.ProtoSerializer{},
			},
		}),
	)
}

func (model TopicsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := util.NewCommandBuilder()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.table.SetHeight(msg.Height - 2)
	}

	msg, adapterCmd := model.adapter.Message(msg)
	cmds.AddCmd(adapterCmd)
	switch msg := msg.(type) {
	case client.CreateConsumerResult:
		log.Println("TOPIC - NEW CONSUMER", msg)
	case client.ConsumeMessage:
		log.Println("ConsumeMessage", msg)
	}

	var tableCmd tea.Cmd
	model.table, tableCmd = model.table.Update(msg)
	cmds.AddCmd(tableCmd)

	return model, cmds.Build()
}

func (model TopicsModel) View() string {
	return baseStyle.Render(model.table.View())
}
