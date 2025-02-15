package cluster

import (
	"log/slog"

	"github.com/anthdm/hollywood/actor"
)

type TopicConfig struct {
	Topic        string
	Discovery    *actor.PID
	Logger       *slog.Logger
	SendMetadata bool
}

type topicActor struct {
	config      TopicConfig
	messagesPID *actor.PID
	consumerPID *actor.PID
	consumers   map[uint64]*actor.PID
	envelopes   []*ConsumerEnvelope
}

func NewTopicActor(config TopicConfig) actor.Producer {
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
		config := NewNodeConfig().
			WithDiscoveryPID(topic.config.Discovery).
			WithLogger(topic.config.Logger)
		config.Topic = topic.config.Topic
		config.StateMachine = act.PID()
		topic.messagesPID = act.SpawnChild(NewNode(config), "node")

	case *actor.Ping:
		act.Send(act.Sender(), &actor.Pong{})

	case *Envelope:
		topic.config.Logger.Info("Received msg")
		act.Engine().SendWithSender(topic.messagesPID, msg, act.Sender())

	case *ConsumerEnvelope:
		topic.envelopes = append(topic.envelopes, msg)
		for _, pid := range topic.consumers {
			act.Send(pid, msg)
		}
		act.Respond(&ConsumerEnvelopeAck{})
		if msg.IsLeader && topic.config.SendMetadata {
			act.Send(act.PID(), &TopicMetadata{
				TopicName:   topic.config.Topic,
				NumMessages: uint64(len(topic.envelopes)),
			})
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
		topic.config.Logger.Info("Registering consumer", "topic", topic.config.Topic, "pid", pid.String())
		topic.consumers[key] = pid
		act.Respond(&RegisterConsumerResult{
			Success: true,
		})

		// Replay all envelopes to new consumer
		for _, envelope := range topic.envelopes {
			act.Send(pid, envelope)
		}

	}
}
