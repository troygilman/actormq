package actormq

import (
	"log"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/troygilman0/actormq/raft/timer"
)

type heartbeatTimeout struct {
}

type nodeMetadata struct {
	pid *actor.PID
}

type ClientConfig struct {
	Nodes []*actor.PID
}

type clientActor struct {
	config         ClientConfig
	nodes          map[string]*nodeMetadata
	leader         *actor.PID
	heartbeatTimer *timer.SendTimer
}

func NewClient(config ClientConfig) actor.Producer {
	return func() actor.Receiver {
		return &clientActor{
			config: config,
		}
	}
}

func (client *clientActor) Receive(act *actor.Context) {
	log.Printf("%s - %T: %+v\n", act.PID().String(), act.Message(), act.Message())
	switch act.Message().(type) {
	case actor.Initialized:
		client.nodes = make(map[string]*nodeMetadata)
		for _, pid := range client.config.Nodes {
			client.nodes[pid.String()] = &nodeMetadata{
				pid: pid,
			}
		}

	case actor.Started:
		client.heartbeatTimer = timer.NewSendTimer(act.Engine(), act.PID(), heartbeatTimeout{}, time.Second)
		client.sendHeartbeat(act)

	case actor.Stopped:
		client.heartbeatTimer.Stop()

	case heartbeatTimeout:
		client.sendHeartbeat(act)
		client.heartbeatTimer.Reset(time.Second)

	case *actor.Pong:

	}
}

func (client *clientActor) sendHeartbeat(act *actor.Context) {
	for _, node := range client.nodes {
		act.Send(node.pid, &actor.Ping{})
	}
}
