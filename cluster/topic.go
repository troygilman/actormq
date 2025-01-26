package cluster

import (
	"log/slog"

	"github.com/anthdm/hollywood/actor"
)

type TopicConfig struct {
	Discovery *actor.PID
	Logger    *slog.Logger
}

type topicActor struct {
	config      TopicConfig
	messagesPID *actor.PID
	consumerPID *actor.PID
	consumers   map[uint64]*actor.PID
}

func NewTopic(config TopicConfig) actor.Producer {
	return func() actor.Receiver {
		return &topicActor{
			config: config,
		}
	}
}

func (topic *topicActor) Receive(act *actor.Context) {
	switch msg := act.Message().(type) {
	case actor.Initialized:
		topic.consumers = make(map[uint64]*actor.PID)

	case actor.Started:
		topic.messagesPID = act.SpawnChild(NewNode(NewNodeConfig().
			WithDiscoveryPID(topic.config.Discovery).
			WithLogger(topic.config.Logger),
		), "node")
		// topic.consumerPID = act.SpawnChild(NewRaftNode(NewRaftNodeConfig().
		// 	WithDiscoveryPID(topic.config.Discovery).
		// 	WithLogger(topic.config.Logger),
		// ), "node", actor.WithID("consumer"))

	case *actor.Ping:
		act.Send(act.Sender(), &actor.Pong{})

	case *Envelope:
		act.Engine().SendWithSender(topic.messagesPID, msg.Message, act.Sender())

	case *MessageProcessedEvent:
		for _, pid := range topic.consumers {
			act.Send(pid, msg)
		}

	case *RegisterConsumer:
		pid := PIDToActorPID(msg.PID)
		key := pid.LookupKey()
		if _, ok := topic.consumers[key]; ok {
			act.Respond(&RegisterConsumerResult{
				Success: false,
				Error:   "consumer is already registered",
			})
			return
		}
		topic.consumers[key] = pid
		act.Respond(&RegisterConsumerResult{
			Success: true,
		})

	}
}
