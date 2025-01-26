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

	// config := cluster.NewRaftNodeConfig().
	// 	WithDiscoveryPID(discoveryPID).
	// 	WithLogger(slog.New(slog.NewJSONHandler(io.Discard, nil))).
	// 	// WithLogger(slog.Default()).
	// 	WithMessageHandler(nil)

	config := cluster.TopicConfig{
		Discovery: discoveryPID,
		Logger:    slog.New(slog.NewJSONHandler(io.Discard, nil)),
	}

	nodes := []*actor.PID{
		engine.Spawn(cluster.NewTopic(config), "topic"),
		engine.Spawn(cluster.NewTopic(config), "topic"),
		engine.Spawn(cluster.NewTopic(config), "topic"),
	}

	clientPID := engine.Spawn(client.NewClient(client.ClientConfig{
		Nodes: nodes,
	}), "client")
	engine.Send(clientPID, client.CreateConsumer{})

	nodePID := nodes[0]
	for {
		start := time.Now()
		result, err := engine.Request(nodePID, &cluster.Message{
			// Data: []byte("dwwdwdw"),
		}, time.Second).Result()
		if err != nil {
			panic(err)
		}
		messageResult, ok := result.(*cluster.MessageResult)
		if !ok {
			panic("result is invalid type")
		}
		log.Println("RESULT", messageResult, "duration:", time.Since(start))
		if messageResult.RedirectPID != nil {
			nodePID = cluster.PIDToActorPID(messageResult.RedirectPID)
		}

		time.Sleep(time.Millisecond)
	}

}
