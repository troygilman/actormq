package main

import (
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/troygilman0/actormq"
)

const (
	local = "127.0.0.1"
)

func main() {
	discovery, err := newDiscovery(local + ":8080")
	if err != nil {
		panic(err)
	}

	_, err = newNode(local+":8081", discovery)
	if err != nil {
		panic(err)
	}

	_, err = newNode(local+":8082", discovery)
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(10 * time.Second)

	}()

	select {}
}

func newDiscovery(address string) (*actor.PID, error) {
	remoter := remote.New(address, remote.NewConfig())
	engine, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(remoter))
	if err != nil {
		return nil, err
	}
	pid := engine.Spawn(actormq.NewDiscovery(), "discovery")
	return pid, err
}

func newNode(address string, discoveryPID *actor.PID) (*actor.PID, error) {
	remoter := remote.New(address, remote.NewConfig())
	engine, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(remoter))
	if err != nil {
		return nil, err
	}
	pid := engine.Spawn(actormq.NewNode(discoveryPID), "node")
	return pid, nil
}
