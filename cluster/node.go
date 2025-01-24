package cluster

import (
	"github.com/anthdm/hollywood/actor"
)

type NodeConfig struct {
}

type nodeActor struct {
	config   NodeConfig
	raftNode *actor.PID
}

func NewNode(config NodeConfig) actor.Producer {
	return func() actor.Receiver {
		return &nodeActor{}
	}
}

func (node *nodeActor) Receive(act *actor.Context) {
	switch act.Message().(type) {
	case actor.Started:
		act.SpawnChild(NewRaftNode(NewRaftNodeConfig()), "raft-node")
	}
}
