package util

import (
	"errors"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/troygilman/actormq/client"
)

func createProducer(engine *actor.Engine, clientPID *actor.PID, topic string) (*actor.PID, error) {
	result, err := engine.Request(clientPID, client.CreateProducer{
		ProducerConfig: client.ProducerConfig{
			Topic:      topic,
			Serializer: remote.ProtoSerializer{},
		},
	}, time.Second).Result()
	if err != nil {
		return nil, err
	}

	if result, ok := result.(client.CreateProducerResult); ok {
		return result.PID, nil
	} else {
		return nil, errors.New("wrong result type")
	}
}
