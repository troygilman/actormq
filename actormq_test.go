package actormq

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/troygilman/actormq/client"
	"github.com/troygilman/actormq/cluster"
	"github.com/troygilman/actormq/internal"
)

func BenchmarkCluster(b *testing.B) {
	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		b.Error(err)
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	discovery := engine.Spawn(cluster.NewDiscovery(cluster.DiscoveryConfig{
		Topics: []string{},
		Logger: logger,
	}), "discovery")

	config := cluster.PodConfig{
		Discovery: discovery,
		Logger:    logger,
	}

	pods := []*actor.PID{
		engine.Spawn(cluster.NewPod(config), "pod", actor.WithID("A")),
		engine.Spawn(cluster.NewPod(config), "pod", actor.WithID("B")),
		engine.Spawn(cluster.NewPod(config), "pod", actor.WithID("C")),
	}

	clientPID := engine.Spawn(client.NewClient(client.ClientConfig{
		Nodes: pods,
	}), "client")

	time.Sleep(time.Second)
	internalTopicProducer, err := createProducer(engine, clientPID, cluster.TopicClusterInternalTopic)
	if err != nil {
		b.Error(err)
	}

	result, err := internal.Request[client.ProduceMessageResult](engine, internalTopicProducer, client.ProduceMessage{
		Message: &cluster.TopicMetadata{
			TopicName: "test",
		},
	})
	if err != nil {
		b.Error(err)
	}
	if result.Error != nil {
		b.Error(result.Error)
	}

	time.Sleep(time.Second)
	testProducer, err := createProducer(engine, clientPID, "test")
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	responses := make([]*actor.Response, b.N)
	for i := range b.N {
		responses[i] = sendMessage(engine, testProducer, &actor.Ping{})
	}

	for i := range b.N {
		result, err := responses[i].Result()
		if err != nil {
			b.Error(err)
		}

		if result, ok := result.(client.ProduceMessageResult); ok {
			if result.Error != nil {
				b.Error(err)
			}
		} else {
			b.Error("wrong result type")
		}
	}
}

func sendMessage(engine *actor.Engine, producerPID *actor.PID, msg any) *actor.Response {
	return engine.Request(producerPID, client.ProduceMessage{
		Message: msg,
	}, 10*time.Second)
}

func createProducer(engine *actor.Engine, clientPID *actor.PID, topic string) (*actor.PID, error) {
	result, err := internal.Request[client.CreateProducerResult](engine, clientPID, client.CreateProducer{
		ProducerConfig: client.ProducerConfig{
			Topic:      topic,
			Serializer: remote.ProtoSerializer{},
		},
	})
	return result.PID, err
}
