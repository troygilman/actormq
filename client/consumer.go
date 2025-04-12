package client

import (
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/troygilman/actormq/cluster"
)

type ConsumerConfig struct {
	Topic        string
	PID          *actor.PID
	Deserializer remote.Deserializer
}

type consumerActor struct {
	config ConsumerConfig
	pods   []*actor.PID
	leader *actor.PID
}

func NewConsumer(config ConsumerConfig, pods []*actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &consumerActor{
			config: config,
			pods:   pods,
			leader: pods[0],
		}
	}
}

func (consumer *consumerActor) Receive(act *actor.Context) {
	switch msg := act.Message().(type) {
	case actor.Started:
		result, err := handleResponse[*cluster.RegisterConsumerResult](act.Request(consumer.leader, &cluster.RegisterConsumer{
			Topic: consumer.config.Topic,
			PID:   cluster.ActorPIDToPID(act.PID()),
		}, 10*time.Second))
		if err != nil {
			panic(err)
		}
		if !result.Success {
			panic(result.Error)
		}
		// log.Println("registered consumer")

	case *cluster.ConsumerEnvelope:
		// log.Printf("consumer of topic %s received envelope with %d messages\n", consumer.config.Topic, len(msg.Messages))
		messages := make([]any, len(msg.Messages))
		for idx, msg := range msg.Messages {
			message, err := consumer.config.Deserializer.Deserialize(msg.Data, msg.TypeName)
			if err != nil {
				panic(err)
			}
			messages[idx] = message
		}

		act.Send(consumer.config.PID, ConsumeMessages{
			Messages: messages,
		})
	}
}
