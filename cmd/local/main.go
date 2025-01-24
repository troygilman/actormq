package main

import (
	"io"
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

	config := cluster.NewRaftNodeConfig().
		WithDiscoveryPID(discoveryPID).
		WithLogger(slog.New(slog.NewJSONHandler(io.Discard, nil))).
		// WithLogger(slog.Default()).
		WithMessageHandler(nil)

	nodes := []*actor.PID{
		engine.Spawn(cluster.NewRaftNode(config), "node"),
		engine.Spawn(cluster.NewRaftNode(config), "node"),
		engine.Spawn(cluster.NewRaftNode(config), "node"),
	}

	engine.Spawn(client.NewClient(client.ClientConfig{
		Nodes: nodes,
	}), "client")

	nodePID := nodes[0]
	for {
		// start := time.Now()
		result, err := engine.Request(nodePID, &cluster.Message{}, time.Second).Result()
		if err != nil {
			panic(err)
		}
		messageResult, ok := result.(*cluster.MessageResult)
		if !ok {
			panic("result is invalid type")
		}
		// log.Println("RESULT", messageResult, "duration:", time.Since(start))
		if messageResult.RedirectPID != nil {
			nodePID = cluster.PIDToActorPID(messageResult.RedirectPID)
		}
		time.Sleep(time.Millisecond)
	}

}
