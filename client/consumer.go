package client

import "github.com/anthdm/hollywood/actor"

type consumerNodeMetadata struct {
	pid *actor.PID
}

type ConsumerConfig struct {
	nodes []*actor.PID
}

type consumerActor struct {
	config ConsumerConfig
	nodes  map[string]*consumerNodeMetadata
	leader *actor.PID
}

func NewConsumerActor(config ConsumerConfig) actor.Producer {
	return func() actor.Receiver {
		return &consumerActor{
			config: config,
		}
	}
}

func (consumer *consumerActor) Receive(act *actor.Context) {
	switch act.Message().(type) {
	case actor.Initialized:
		consumer.nodes = make(map[string]*consumerNodeMetadata)
		for _, pid := range consumer.config.nodes {
			consumer.nodes[pid.String()] = &consumerNodeMetadata{
				pid: pid,
			}
		}

	}
}
