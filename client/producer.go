package client

import (
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/troygilman/actormq/cluster"
)

type ProducerConfig struct {
	Topic      string
	Serializer remote.Serializer
}

type producerActor struct {
	config ProducerConfig
	pods   []*actor.PID
	leader *actor.PID
}

func NewProducer(config ProducerConfig, pods []*actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &producerActor{
			config: config,
			pods:   pods,
			leader: pods[0],
		}
	}
}

func (producer *producerActor) Receive(act *actor.Context) {
	switch msg := act.Message().(type) {
	case ProduceMessage:

		data, err := producer.config.Serializer.Serialize(msg.Message)
		if err != nil {
			panic(err)
		}

		envelope := &cluster.Envelope{
			Topic: producer.config.Topic,
			Message: &cluster.Message{
				TypeName: producer.config.Serializer.TypeName(msg.Message),
				Data:     data,
			},
		}

		for {
			result, err := handleResponse[*cluster.EnvelopeResult](act.Request(producer.leader, envelope, 10*time.Second))
			if err != nil {
				panic(err)
			}
			if !result.Success {
				if result.RedirectPID != nil {
					producer.leader = cluster.PIDToActorPID(result.RedirectPID)
					continue
				}
				time.Sleep(10 * time.Millisecond)
			}
			break
		}

		act.Respond(ProduceMessageResult{})
	}
}
