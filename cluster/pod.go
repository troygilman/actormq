package cluster

import (
	"fmt"
	"log/slog"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
)

const (
	TopicClusterInternalTopic = "cluster.internal.topic"
)

type PodConfig struct {
	Discovery *actor.PID
	Logger    *slog.Logger
}

type podActor struct {
	config       PodConfig
	topics       map[string]*actor.PID
	adapter      *actor.PID
	serializer   remote.Serializer
	deserializer remote.Deserializer
}

func NewPod(config PodConfig) actor.Producer {
	return func() actor.Receiver {
		return &podActor{
			config:       config,
			serializer:   remote.ProtoSerializer{},
			deserializer: remote.ProtoSerializer{},
		}
	}
}

func (pod *podActor) Receive(act *actor.Context) {
	switch msg := act.Message().(type) {
	case actor.Initialized:
		pod.topics = make(map[string]*actor.PID)

	case actor.Started:
		// INTERNAL TOPICS
		internalTopicsActor := act.SpawnChild(NewTopicActor(TopicConfig{
			Topic:     TopicClusterInternalTopic,
			Discovery: pod.config.Discovery,
			Logger:    pod.config.Logger,
		}), "topic", actor.WithID(TopicClusterInternalTopic))
		act.Send(internalTopicsActor, &RegisterConsumer{
			PID: ActorPIDToPID(act.PID()),
		})
		pod.topics[TopicClusterInternalTopic] = internalTopicsActor

	case actor.Stopped:

	case *Envelope:
		topic, ok := pod.topics[msg.Topic]
		if !ok {
			act.Respond(&EnvelopeResult{
				Success: false,
				Error:   fmt.Sprintf("topic does not exist: %s", msg.Topic),
			})
			return
		}
		act.Engine().SendWithSender(topic, msg, act.Sender())

	case *RegisterConsumer:
		topic, ok := pod.topics[msg.Topic]
		if !ok {
			act.Respond(&RegisterConsumerResult{
				Success: false,
				Error:   fmt.Sprintf("topic does not exist: %s", msg.Topic),
			})
			return
		}
		act.Engine().SendWithSender(topic, msg, act.Sender())

	case *ConsumerEnvelope:
		pod.config.Logger.Info("Consuming envelope", "msg", msg.Message)
		message, err := pod.deserializer.Deserialize(msg.Message.Data, msg.Message.TypeName)
		if err != nil {
			panic(err)
		}
		switch msg := message.(type) {
		case *TopicMetadata:
			if _, ok := pod.topics[msg.TopicName]; !ok {
				pod.topics[msg.TopicName] = act.SpawnChild(NewTopicActor(TopicConfig{
					Topic:        msg.TopicName,
					Discovery:    pod.config.Discovery,
					Logger:       pod.config.Logger,
					SendMetadata: true,
				}), "topic", actor.WithID(msg.TopicName))
			}
		}
		act.Respond(&ConsumerEnvelopeAck{})

	case *TopicMetadata:
		data, err := pod.serializer.Serialize(msg)
		if err != nil {
			panic(err)
		}
		typeName := pod.serializer.TypeName(msg)
		internalTopic := pod.topics[TopicClusterInternalTopic]
		act.Send(internalTopic, &Envelope{
			Message: &Message{
				TypeName: typeName,
				Data:     data,
			},
		})
	}
}
