package raft

import (
	"math/rand"
	"time"

	"github.com/anthdm/hollywood/actor"
)

func newElectionTimoutDuration() time.Duration {
	return time.Duration(minElectionTimeoutMs+rand.Intn(maxElectionTimeoutMs-minElectionTimeoutMs)) * time.Millisecond
}

func sendWithDelay(act *actor.Context, pid *actor.PID, msg any, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		act.Send(pid, msg)
	}()
}

func pidEquals(pid1 *actor.PID, pid2 *actor.PID) bool {
	if pid1 == nil && pid2 == nil {
		return true
	}
	if pid1 == nil || pid2 == nil {
		return false
	}
	return pid1.String() == pid2.String()
}
