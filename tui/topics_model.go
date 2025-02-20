package tui

import (
	"log"
	"strconv"
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
			Message: &cluster.TopicMetadata{
				TopicName: "test.1",
			},
		})
		engine.Send(result.PID, client.ProduceMessage{
			Message: &cluster.TopicMetadata{
				TopicName: "test.2",
			},
		})
	} else {
		panic("could not create consumer")
	}

	return TopicsModel{
		engine:  engine,
		client:  clientPID,
		adapter: util.NewAdapter(engine, util.BasicAdapterFunc),
		table: table.New(
			table.WithKeyMap(table.DefaultKeyMap()),
			table.WithFocused(true),
			table.WithWidth(30),
			table.WithColumns([]table.Column{
				{Title: "Topic", Width: 20},
				{Title: "Messages", Width: 10},
			}),
		),
		topics: make(map[string]int),
	}
}

type TopicsModel struct {
	engine  *actor.Engine
	client  *actor.PID
	adapter util.Adapter
	table   table.Model
	topics  map[string]int
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
	case tea.KeyMsg:
	case tea.WindowSizeMsg:
		model.table.SetHeight(msg.Height - 2)
	}

	var tableCmd tea.Cmd
	model.table, tableCmd = model.table.Update(msg)
	cmds.AddCmd(tableCmd)

	msg, adapterCmd := model.adapter.Message(msg)
	cmds.AddCmd(adapterCmd)
	switch msg := msg.(type) {
	case client.CreateConsumerResult:
		log.Println("TOPIC - NEW CONSUMER", msg)
	case client.ConsumeMessage:
		log.Printf("ConsumeMessage %T\n", msg)
		switch msg := msg.Message.(type) {
		case *cluster.TopicMetadata:
			if index, ok := model.topics[msg.TopicName]; ok {
				model.table.Rows()[index] = []string{
					msg.TopicName,
					strconv.FormatUint(msg.NumMessages, 10),
				}
			} else {
				model.table.SetRows(append(model.table.Rows(), []string{
					msg.TopicName,
					strconv.FormatUint(msg.NumMessages, 10),
				}))
				model.topics[msg.TopicName] = len(model.table.Rows()) - 1
			}
		}
	}

	return model, cmds.Build()
}

func (model TopicsModel) View() string {
	return baseStyle.Render(model.table.View())
}
