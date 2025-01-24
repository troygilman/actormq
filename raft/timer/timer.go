package timer

import (
	"time"

	"github.com/anthdm/hollywood/actor"
)

type SendTimer struct {
	resetChan chan<- time.Duration
	stopChann chan<- struct{}
}

func NewSendTimer(engine *actor.Engine, pid *actor.PID, msg any, duration time.Duration) *SendTimer {
	resetChan := make(chan time.Duration)
	stopChan := make(chan struct{})
	timer := time.NewTimer(duration)
	go func() {
		stopped := false
		for {
			select {
			case <-timer.C:
				engine.Send(pid, msg)
				stopped = true
			case duration := <-resetChan:
				// log.Println("reset", pid.String(), reflect.TypeOf(msg))
				if !stopped && !timer.Stop() {
					// log.Println("reset-1", pid.String(), reflect.TypeOf(msg))
					<-timer.C
				}
				timer.Reset(duration)
				stopped = false
			case <-stopChan:
				timer.Stop()
				return
			}
		}
	}()
	return &SendTimer{
		resetChan: resetChan,
		stopChann: stopChan,
	}
}

func (st *SendTimer) Reset(duration time.Duration) {
	st.resetChan <- duration
}

func (st *SendTimer) Stop() {
	st.stopChann <- struct{}{}
}
