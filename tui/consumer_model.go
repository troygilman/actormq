package tui

import (
	"fmt"
	"log"

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
		table: table.New(
			table.WithWidth(100),
			table.WithColumns([]table.Column{
				{Title: "Offset", Width: 10},
				{Title: "Message", Width: 90},
			}),
		),
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
		tea.WindowSize(),
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
	var cmd tea.Cmd
	cmds := util.NewCommandBuilder()

	switch msg := msg.(type) {
	case tea.KeyMsg:
	case tea.WindowSizeMsg:
		model.table.SetHeight(msg.Height - 3)
	case FocusMsg:
		model.table = util.SetTableFocus(model.table, msg.Focus)
	}

	model.table, cmd = model.table.Update(msg)
	cmds.AddCmd(cmd)

	msg, adapterCmd := model.adapter.Message(msg)
	cmds.AddCmd(adapterCmd)
	switch msg := msg.(type) {
	case client.ConsumeMessages:
		log.Printf("model is consuming %d messages\n", len(msg.Messages))
		for _, msg := range msg.Messages {
			rows := model.table.Rows()
			rows = append(rows, table.Row{
				fmt.Sprintf("%d", len(model.table.Rows())),
				fmt.Sprintf("%T", msg),
			})
			model.table.SetRows(rows)
		}
	default:
		_ = msg
	}

	return model, cmds.Build()
}

func (model ConsumerModel) View() string {
	return model.table.View()
}
