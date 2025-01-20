package main

import (
	"log/slog"
	"os"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/troygilman0/actormq/raft"
)

const (
	discoveryAddr = "127.0.0.1:8080"
	discoveryID   = "discovery/primary"
)

func main() {
	port := os.Args[1]

	remoter := remote.New("127.0.0.1:"+port, remote.NewConfig().WithMaxRetries(1))
	engine, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(remoter))
	if err != nil {
		panic(err)
	}

	engine.Spawn(raft.NewNode(raft.NodeConfig{
		DiscoveryPID: actor.NewPID(discoveryAddr, discoveryID),
		Handler:      handle,
		Logger:       slog.Default(),
	}), "node")

	select {}
}

func handle(command string) {

}
