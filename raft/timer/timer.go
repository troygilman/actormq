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
		for {
			select {
			case <-timer.C:
				engine.Send(pid, msg)
			case duration := <-resetChan:
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(duration)
			case <-stopChan:
				timer.Stop()
				break
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
