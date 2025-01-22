package main

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/troygilman0/actormq/discovery"
	"github.com/troygilman0/actormq/raft"
)

func main() {
	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}

	discoveryPID := engine.Spawn(discovery.NewDiscovery(), "discovery")

	engine.Spawn(raft.NewNode(), "node")
	engine.Spawn(raft.NewNode(), "node")
	node := engine.Spawn(raft.NewNode(), "node")

}
