package raft

import (
	"github.com/anthdm/hollywood/actor"
)

type CommandHandler func(command string)

type nodeStatus uint

const (
	nodeStatusClosed nodeStatus = iota
	nodeStatusRunning
	nodeStatusClosing
)

type nodeMetadata struct {
	pid        *actor.PID
	nextIndex  uint64
	matchIndex uint64
}

type commandMetadata struct {
	sender *actor.PID
}
