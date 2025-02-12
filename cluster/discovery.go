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
	Topics []string
	Logger *slog.Logger
}

type discoveryActor struct {
	config   DiscoveryConfig
	topics   map[string]map[uint64]struct{}
	nodes    map[uint64]*discoveryNodeMetadata
	repeater actor.SendRepeater
}

func NewDiscovery(config DiscoveryConfig) actor.Producer {
	return func() actor.Receiver {
		return &discoveryActor{
			config: config,
		}
	}
}

func (d *discoveryActor) Receive(act *actor.Context) {
	switch msg := act.Message().(type) {
	case actor.Initialized:
		d.nodes = make(map[uint64]*discoveryNodeMetadata)
		d.topics = make(map[string]map[uint64]struct{})

		for _, topic := range d.config.Topics {
			d.topics[topic] = make(map[uint64]struct{})
		}

	case actor.Started:
		d.repeater = act.SendRepeat(act.PID(), sendPing{}, time.Second)

	case actor.Stopped:
		d.repeater.Stop()

	case *RegisterNode:
		pid := act.Sender()
		topic, ok := d.topics[msg.Topic]
		if !ok {
			topic = make(map[uint64]struct{})
			d.topics[msg.Topic] = topic
		}
		topic[pid.LookupKey()] = struct{}{}
		d.nodes[pid.LookupKey()] = &discoveryNodeMetadata{
			pid:      pid,
			lastPong: time.Now(),
		}
		d.sendActiveNodes(act, msg.Topic)
		// log.Println("registered node", pid.String(), "on topic", msg.Topic)

	case sendPing:
		updatedTopics := make(map[string]struct{})
		for topic, keys := range d.topics {
			for key := range keys {
				node := d.nodes[key]
				if time.Since(node.lastPong) > 5*time.Second {
					delete(d.nodes, key)
					delete(keys, key)
					updatedTopics[topic] = struct{}{}
				} else {
					act.Send(node.pid, &actor.Ping{})
				}
			}
		}
		for topic := range updatedTopics {
			d.sendActiveNodes(act, topic)
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

func (d *discoveryActor) sendActiveNodes(act *actor.Context, topic string) {
	nodes := make([]*PID, 0)
	for key := range d.topics[topic] {
		nodes = append(nodes, ActorPIDToPID(d.nodes[key].pid))
	}
	for key := range d.topics[topic] {
		act.Send(d.nodes[key].pid, &ActiveNodes{
			Nodes: nodes,
		})
	}
}
