package actormq

import "github.com/anthdm/hollywood/actor"

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
