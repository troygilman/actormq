package main

import (
	"io"
	"log/slog"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/troygilman0/actormq"
	"github.com/troygilman0/actormq/raft"
)

func main() {
	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}

	discoveryPID := engine.Spawn(raft.NewDiscovery(), "discovery")

	config := raft.NewNodeConfig().
		WithDiscoveryPID(discoveryPID).
		WithLogger(slog.New(slog.NewJSONHandler(io.Discard, nil))).
		// WithLogger(slog.Default()).
		WithMessageHandler(nil)

	nodes := []*actor.PID{
		engine.Spawn(raft.NewNode(config), "node"),
		engine.Spawn(raft.NewNode(config), "node"),
		engine.Spawn(raft.NewNode(config), "node"),
	}

	engine.Spawn(actormq.NewClient(actormq.ClientConfig{
		Nodes: nodes,
	}), "client")

	nodePID := nodes[0]
	for {
		// start := time.Now()
		result, err := engine.Request(nodePID, &raft.Message{}, time.Second).Result()
		if err != nil {
			panic(err)
		}
		messageResult, ok := result.(*raft.MessageResult)
		if !ok {
			panic("result is invalid type")
		}
		// log.Println("RESULT", messageResult, "duration:", time.Since(start))
		if messageResult.RedirectPID != nil {
			nodePID = raft.PIDToActorPID(messageResult.RedirectPID)
		}
		time.Sleep(time.Millisecond)
	}

}
