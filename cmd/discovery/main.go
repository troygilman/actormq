package main

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/troygilman0/actormq/discovery"
)

const (
	address = "127.0.0.1:8080"
)

func main() {
	remoter := remote.New(address, remote.NewConfig().WithMaxRetries(1))
	engine, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(remoter))
	if err != nil {
		panic(err)
	}

	engine.Spawn(discovery.NewDiscovery(), "discovery", actor.WithID("primary"))

	select {}
}
