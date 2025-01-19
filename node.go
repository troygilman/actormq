package actormq

import (
	"log"

	"github.com/anthdm/hollywood/actor"
)

type node struct {
	discoveryPID *actor.PID
}

func NewNode(discoveryPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &node{
			discoveryPID: discoveryPID,
		}
	}
}

func (n *node) Receive(act *actor.Context) {
	log.Printf("%s - %T: %v\n", act.PID().String(), act.Message(), act.Message())
	switch act.Message().(type) {
	case actor.Started:
		act.Send(n.discoveryPID, &RegisterNode{})
	case *ActiveNodes:
	case *actor.Ping:
		act.Send(act.Sender(), &actor.Pong{})
	}
}
