package client

import (
	"log/slog"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/troygilman/actormq/cluster"
)

type heartbeatTimeout struct{}

type nodeMetadata struct {
	pid           *actor.PID
	lastHeartbeat time.Time
}

type ClientConfig struct {
	Nodes []*actor.PID
}

type clientActor struct {
	config            ClientConfig
	nodes             map[string]*nodeMetadata
	leader            *actor.PID
	heartbeatRepeater actor.SendRepeater
	clusterState      *cluster.ClusterState
}

func NewClient(config ClientConfig) actor.Producer {
	return func() actor.Receiver {
		return &clientActor{
			config: config,
		}
	}
}

func (client *clientActor) Receive(act *actor.Context) {
	// log.Printf("%s - %T: %+v\n", act.PID().String(), act.Message(), act.Message())
	switch msg := act.Message().(type) {
	case actor.Initialized:
		client.nodes = make(map[string]*nodeMetadata)
		for _, pid := range client.config.Nodes {
			client.nodes[pid.String()] = &nodeMetadata{
				pid:           pid,
				lastHeartbeat: time.Now(),
			}
		}

	case actor.Started:
		client.heartbeatRepeater = act.SendRepeat(act.PID(), heartbeatTimeout{}, time.Second)
		client.sendHeartbeat(act)

	case actor.Stopped:
		client.heartbeatRepeater.Stop()

	case heartbeatTimeout:
		client.sendHeartbeat(act)

	case *actor.Pong:
		if node, ok := client.nodes[act.Sender().String()]; ok {
			node.lastHeartbeat = time.Now()
		}

	case CreateConsumer:
		act.Respond(CreateConsumerResult{
			PID: act.SpawnChild(NewConsumer(msg.ConsumerConfig, client.config.Nodes), "consumer"),
		})

	case CreateProducer:
		slog.Default().Info("Create producer", "topic", msg.ProducerConfig.Topic)
		act.Respond(CreateProducerResult{
			PID: act.SpawnChild(NewProducer(msg.ProducerConfig, client.config.Nodes), "producer"),
		})

	case *cluster.ClusterState:

	}
}

func (client *clientActor) sendHeartbeat(act *actor.Context) {
	for _, node := range client.nodes {
		act.Send(node.pid, &actor.Ping{})
		act.Send(node.pid, &cluster.GetClusterState{})
	}
}
