package tui

import (
	"fmt"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/troygilman/actormq/client"
	"github.com/troygilman/actormq/tui/util"
)

func NewConsumerModel(engine *actor.Engine, client *actor.PID, topic string) tea.Model {
	return ConsumerModel{
		engine:  engine,
		client:  client,
		topic:   topic,
		adapter: util.NewAdapter(engine, util.BasicAdapterFunc),
		table: table.New(table.WithWidth(100), table.WithColumns([]table.Column{
			{Title: "Message", Width: 100},
		})),
	}
}

type ConsumerModel struct {
	engine  *actor.Engine
	client  *actor.PID
	topic   string
	adapter util.Adapter
	table   table.Model
}

func (model ConsumerModel) Init() tea.Cmd {
	return tea.Batch(
		model.adapter.Init(),
		model.adapter.Send(model.client, client.CreateConsumer{
			ConsumerConfig: client.ConsumerConfig{
				Topic:        model.topic,
				PID:          model.adapter.PID(),
				Deserializer: remote.ProtoSerializer{},
			},
		}),
	)
}

func (model ConsumerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := util.NewCommandBuilder()

	msg, adapterCmd := model.adapter.Message(msg)
	cmds.AddCmd(adapterCmd)
	switch msg := msg.(type) {
	case client.ConsumeMessage:
		model.table.SetRows(append(model.table.Rows(), table.Row{fmt.Sprintf("%T", msg.Message)}))
	}

	return model, cmds.Build()
}

func (model ConsumerModel) View() string {
	return model.table.View()
}
