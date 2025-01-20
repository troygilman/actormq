package actormq

import "github.com/anthdm/hollywood/actor"

func PIDToActorPID(pid *PID) *actor.PID {
	return actor.NewPID(pid.GetAddress(), pid.GetID())
}

func ActorPIDToPID(pid *actor.PID) *PID {
	return &PID{
		Address: pid.GetAddress(),
		ID:      pid.GetID(),
	}
}
