package main

import (
	"log/slog"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/troygilman/actormq/cluster"
)

func main() {
	remoter := remote.New("127.0.0.1:8080", remote.NewConfig().WithMaxRetries(1))

	engine, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(remoter))
	if err != nil {
		panic(err)
	}

	discovery := engine.Spawn(cluster.NewDiscovery(cluster.DiscoveryConfig{
		Topics: []string{},
		Logger: slog.Default(),
	}), "discovery")

	config := cluster.PodConfig{
		Discovery: discovery,
		Logger:    slog.Default(),
	}

	engine.Spawn(cluster.NewPod(config), "pod", actor.WithID("A"))
	engine.Spawn(cluster.NewPod(config), "pod", actor.WithID("B"))
	engine.Spawn(cluster.NewPod(config), "pod", actor.WithID("C"))

	select {}
}
