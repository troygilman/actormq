package cluster

import (
	"log"
	"log/slog"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type (
	sendPing              struct{}
	discoveryNodeMetadata struct {
		pid      *actor.PID
		lastPong time.Time
	}
)

type DiscoveryConfig struct {
	Logger *slog.Logger
}

type discoveryActor struct {
	nodes    map[uint64]*discoveryNodeMetadata
	repeater actor.SendRepeater
}

func NewDiscovery() actor.Producer {
	return func() actor.Receiver {
		return &discoveryActor{}
	}
}

func (d *discoveryActor) Receive(act *actor.Context) {
	switch act.Message().(type) {
	case actor.Initialized:
		d.nodes = make(map[uint64]*discoveryNodeMetadata)

	case actor.Started:
		d.repeater = act.SendRepeat(act.PID(), sendPing{}, time.Second)

	case actor.Stopped:
		d.repeater.Stop()

	case *RegisterNode:
		pid := act.Sender()
		d.nodes[pid.LookupKey()] = &discoveryNodeMetadata{
			pid:      pid,
			lastPong: time.Now(),
		}
		d.sendActiveNodes(act)
		log.Println("registered node:", pid.String())

	case sendPing:
		shouldSendActiveNodes := false
		for key, node := range d.nodes {
			if time.Since(node.lastPong) > 5*time.Second {
				delete(d.nodes, key)
				shouldSendActiveNodes = true
			} else {
				act.Send(node.pid, &actor.Ping{})
			}
		}
		if shouldSendActiveNodes {
			d.sendActiveNodes(act)
		}

	case *actor.Pong:
		pid := act.Sender()
		node, ok := d.nodes[pid.LookupKey()]
		if !ok {
			log.Println("could not find node:", pid.String())
			return
		}
		node.lastPong = time.Now()
	}
}

func (d *discoveryActor) sendActiveNodes(act *actor.Context) {
	nodes := make([]*PID, 0)
	for _, node := range d.nodes {
		nodes = append(nodes, ActorPIDToPID(node.pid))
	}
	for _, node := range d.nodes {
		act.Send(node.pid, &ActiveNodes{
			Nodes: nodes,
		})
	}
}
