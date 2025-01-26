package cluster

import (
	"github.com/anthdm/hollywood/actor"
)

type nodeMetadata struct {
	pid        *actor.PID
	nextIndex  uint64
	matchIndex uint64
}

type commandMetadata struct {
	sender *actor.PID
}
