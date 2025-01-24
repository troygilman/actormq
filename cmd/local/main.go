package main

import (
	"io"
	"log/slog"

	"github.com/anthdm/hollywood/actor"
	"github.com/troygilman0/actormq"
	"github.com/troygilman0/actormq/discovery"
	"github.com/troygilman0/actormq/raft"
)

func main() {
	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}

	discoveryPID := engine.Spawn(discovery.NewDiscovery(), "discovery")

	config := raft.NewNodeConfig().
		WithDiscoveryPID(discoveryPID).
		WithLogger(slog.New(slog.NewJSONHandler(io.Discard, nil))).
		// WithLogger(slog.Default()).
		WithCommandHandler(func(command string) {})

	nodes := []*actor.PID{
		engine.Spawn(raft.NewNode(config), "node"),
		engine.Spawn(raft.NewNode(config), "node"),
		engine.Spawn(raft.NewNode(config), "node"),
	}

	engine.Spawn(actormq.NewClient(actormq.ClientConfig{
		Nodes: nodes,
	}), "client")

	select {}
	// nodePID := nodes[0]
	// for {
	// 	start := time.Now()
	// 	result, err := engine.Request(nodePID, &raft.Command{
	// 		Command: "Hello World",
	// 	}, time.Second).Result()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	commandResult, ok := result.(*raft.CommandResult)
	// 	if !ok {
	// 		panic("result is invalid type")
	// 	}
	// 	log.Println("RESULT", commandResult, "duration:", time.Since(start))
	// 	if commandResult.RedirectPID != nil {
	// 		nodePID = actormq.PIDToActorPID(commandResult.RedirectPID)
	// 	}
	// 	time.Sleep(time.Millisecond)
	// }

}
