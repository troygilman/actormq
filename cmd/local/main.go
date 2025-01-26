package main

import (
	"io"
	"log/slog"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/troygilman0/actormq/client"
	"github.com/troygilman0/actormq/cluster"
)

func main() {
	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}

	discoveryPID := engine.Spawn(cluster.NewDiscovery(), "discovery")

	config := cluster.PodConfig{
		Topics:    []string{"test"},
		Discovery: discoveryPID,
		Logger:    slog.New(slog.NewJSONHandler(io.Discard, nil)),
		// Logger: slog.Default(),
	}

	pods := []*actor.PID{
		engine.Spawn(cluster.NewPod(config), "pod"),
		engine.Spawn(cluster.NewPod(config), "pod"),
		engine.Spawn(cluster.NewPod(config), "pod"),
	}

	clientPID := engine.Spawn(client.NewClient(client.ClientConfig{
		Nodes: pods,
	}), "client")

	engine.Send(clientPID, client.CreateConsumer{
		ConsumerConfig: client.ConsumerConfig{
			Topic:        "test",
			Deserializer: remote.ProtoSerializer{},
		},
	})

	response := engine.Request(clientPID, client.CreateProducer{
		ProducerConfig: client.ProducerConfig{
			Topic:      "test",
			Serializer: remote.ProtoSerializer{},
		},
	}, time.Second)
	result, err := response.Result()
	if err != nil {
		panic(err)
	}
	producerPID := result.(client.CreateProducerResult).PID

	for {
		engine.Send(producerPID, client.ProduceMessage{
			Message: &actor.Ping{},
		})
		time.Sleep(time.Second)
	}
}
