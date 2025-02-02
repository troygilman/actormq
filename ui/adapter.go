package ui

import "github.com/anthdm/hollywood/actor"

type adapterActor struct {
	c chan<- any
}

func NewAdapter(c chan<- any) actor.Producer {
	return func() actor.Receiver {
		return &adapterActor{
			c: c,
		}
	}
}

func (a *adapterActor) Receive(act *actor.Context) {
	switch msg := act.Message().(type) {
	case actor.Initialized, actor.Started:
	case actor.Stopped:
		close(a.c)
	default:
		a.c <- msg
	}
}
