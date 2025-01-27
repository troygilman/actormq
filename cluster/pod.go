package cluster

import (
	"log/slog"

	"github.com/anthdm/hollywood/actor"
)

type PodConfig struct {
	Topics    []string
	Discovery *actor.PID
	Logger    *slog.Logger
}

type podActor struct {
	config PodConfig
	topics map[string]*actor.PID
}

func NewPod(config PodConfig) actor.Producer {
	return func() actor.Receiver {
		return &podActor{
			config: config,
		}
	}
}

func (pod *podActor) Receive(act *actor.Context) {
	switch msg := act.Message().(type) {
	case actor.Initialized:
		pod.topics = make(map[string]*actor.PID)

	case actor.Started:
		for _, topic := range pod.config.Topics {
			pod.topics[topic] = act.SpawnChild(NewTopic(TopicConfig{
				Topic:     topic,
				Discovery: pod.config.Discovery,
				Logger:    pod.config.Logger,
			}), "topic", actor.WithID(topic))
		}

	case actor.Stopped:

	case *Envelope:
		topic, ok := pod.topics[msg.Topic]
		if !ok {
			act.Respond(&EnvelopeResult{
				Success: false,
				Error:   "topic does not exist",
			})
			return
		}
		act.Engine().SendWithSender(topic, msg, act.Sender())

	case *RegisterConsumer:
		topic, ok := pod.topics[msg.Topic]
		if !ok {
			act.Respond(&RegisterConsumerResult{
				Success: false,
				Error:   "topic does not exist",
			})
			return
		}
		act.Engine().SendWithSender(topic, msg, act.Sender())
	}
}
