package cluster

import (
	"log/slog"

	"github.com/anthdm/hollywood/actor"
)

type NodeConfig struct {
	Discovery *actor.PID
	Logger    *slog.Logger
}

type nodeActor struct {
	config   NodeConfig
	raftNode *actor.PID
}

func NewNode(config NodeConfig) actor.Producer {
	return func() actor.Receiver {
		return &nodeActor{
			config: config,
		}
	}
}

func (node *nodeActor) Receive(act *actor.Context) {
	switch msg := act.Message().(type) {
	case actor.Started:
		node.raftNode = act.SpawnChild(NewRaftNode(NewRaftNodeConfig().
			WithDiscoveryPID(node.config.Discovery).
			WithLogger(node.config.Logger),
		), "raft-node")

	case *actor.Ping:
		act.Send(act.Sender(), &actor.Pong{})

	case *Message:
		act.Engine().SendWithSender(node.raftNode, msg, act.Sender())
	}
}
