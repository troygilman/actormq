package main

import (
	"io"
	"log"
	"log/slog"
	"time"

	"github.com/anthdm/hollywood/actor"
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
		Topic: "test",
	})

	pid := pods[0]
	for {
		start := time.Now()
		result, err := engine.Request(pid, &cluster.Envelope{
			Topic: "test",
			Message: &cluster.Message{
				Data: []byte("dwwdwdw"),
			},
		}, time.Second).Result()
		if err != nil {
			panic(err)
		}
		envelopeResult, ok := result.(*cluster.EnvelopeResult)
		if !ok {
			panic("result is invalid type")
		}
		log.Println("RESULT", envelopeResult, "duration:", time.Since(start), pid)
		if envelopeResult.RedirectPID != nil {
			pid = cluster.PIDToActorPID(envelopeResult.RedirectPID)
		}

		time.Sleep(time.Second)
	}

}
