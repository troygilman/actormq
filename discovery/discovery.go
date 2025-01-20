package discovery

import (
	"log"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/troygilman0/actormq"
)

type (
	sendPing     struct{}
	nodeMetadata struct {
		pid      *actor.PID
		lastPong time.Time
	}
)

type discovery struct {
	nodes    map[string]*nodeMetadata
	repeater actor.SendRepeater
}

func NewDiscovery() actor.Producer {
	return func() actor.Receiver {
		return &discovery{}
	}
}

func (d *discovery) Receive(act *actor.Context) {
	log.Printf("%s - %T: %+v\n", act.PID().String(), act.Message(), act.Message())

	switch act.Message().(type) {
	case actor.Initialized:
		d.nodes = make(map[string]*nodeMetadata)

	case actor.Started:
		d.repeater = act.SendRepeat(act.PID(), sendPing{}, time.Second)

	case actor.Stopped:
		d.repeater.Stop()

	case *actormq.RegisterNode:
		pid := act.Sender()
		d.nodes[pid.String()] = &nodeMetadata{
			pid:      pid,
			lastPong: time.Now(),
		}
		d.sendActiveNodes(act)
		log.Println("registered node:", pid.String())

	case sendPing:
		shouldSendActiveNodes := false
		for key, node := range d.nodes {
			if time.Since(node.lastPong) > 3*time.Second {
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
		node, ok := d.nodes[pid.String()]
		if !ok {
			log.Println("could not find node:", pid.String())
			return
		}
		node.lastPong = time.Now()
	}
}

func (d *discovery) sendActiveNodes(act *actor.Context) {
	nodes := make([]*actormq.PID, 0)
	for _, node := range d.nodes {
		nodes = append(nodes, actormq.ActorPIDToPID(node.pid))
	}
	for _, node := range d.nodes {
		act.Send(node.pid, &actormq.ActiveNodes{
			Nodes: nodes,
		})
	}
}
