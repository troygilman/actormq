package tui

import (
	"strconv"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/troygilman/actormq/client"
	"github.com/troygilman/actormq/cluster"
	"github.com/troygilman/actormq/internal"
	"github.com/troygilman/actormq/tui/util"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func NewTopicsModel(engine *actor.Engine, clientPID *actor.PID) tea.Model {
	{
		result, err := internal.Request[client.CreateProducerResult](engine, clientPID, client.CreateProducer{
			ProducerConfig: client.ProducerConfig{
				Topic:      cluster.TopicClusterInternalTopic,
				Serializer: remote.ProtoSerializer{},
			},
		})
		if err != nil {
			panic(err)
		}

		internal.Request[client.ProduceMessageResult](engine, result.PID, client.ProduceMessage{
			Message: &cluster.TopicMetadata{
				TopicName: "test.1",
			},
		})

		internal.Request[client.ProduceMessageResult](engine, result.PID, client.ProduceMessage{
			Message: &cluster.TopicMetadata{
				TopicName: "test.2",
			},
		})
	}

	{
		time.Sleep(time.Second)
		result, err := internal.Request[client.CreateProducerResult](engine, clientPID, client.CreateProducer{
			ProducerConfig: client.ProducerConfig{
				Topic:      "test.1",
				Serializer: remote.ProtoSerializer{},
			},
		})
		if err != nil {
			panic(err)
		}

		go func() {
			for {
				engine.Send(result.PID, client.ProduceMessage{
					Message: &actor.Ping{},
				})
				time.Sleep(100 * time.Millisecond)
			}
		}()
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
		switch msg.String() {
		case "c":
			cmds.AddCmd(util.MessageCmd(CreateConsumerMsg{
				Topic: model.table.Rows()[model.table.Cursor()][0],
			}))
		}
	case tea.WindowSizeMsg:
		model.table.SetHeight(msg.Height - 2)
	case FocusMsg:
		model.table = util.SetTableFocus(model.table, msg.Focus)
	}

	var tableCmd tea.Cmd
	model.table, tableCmd = model.table.Update(msg)
	cmds.AddCmd(tableCmd)

	msg, adapterCmd := model.adapter.Message(msg)
	cmds.AddCmd(adapterCmd)
	switch msg := msg.(type) {
	case client.CreateConsumerResult:
	case client.ConsumeMessage:
		switch msg := msg.Message.(type) {
		case *cluster.TopicMetadata:
			if index, ok := model.topics[msg.TopicName]; ok {
				rows := model.table.Rows()
				rows[index] = []string{
					msg.TopicName,
					strconv.FormatUint(msg.NumMessages, 10),
				}
				model.table.SetRows(rows)
			} else {
				model.table.SetRows(append(model.table.Rows(), table.Row{
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
	return model.table.View()
}
