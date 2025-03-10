package cluster

import (
	"math/rand"
	"strings"
	"time"

	"github.com/anthdm/hollywood/actor"
)

func newElectionTimoutDuration(config NodeConfig) time.Duration {
	return config.ElectionMinInterval + time.Duration(rand.Intn(int(config.ElectionMaxInterval)-int(config.ElectionMinInterval)))
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

func PIDToActorPID(pid *PID) *actor.PID {
	if pid == nil {
		return nil
	}
	return actor.NewPID(pid.GetAddress(), pid.GetID())
}

func ActorPIDToPID(pid *actor.PID) *PID {
	if pid == nil {
		return nil
	}
	return &PID{
		Address: pid.GetAddress(),
		ID:      pid.GetID(),
	}
}

func ParentPID(pid *actor.PID) *actor.PID {
	id := pid.GetID()
	id = id[:strings.LastIndex(id, "/")]
	id = id[:strings.LastIndex(id, "/")]
	return actor.NewPID(pid.GetAddress(), id)
}
