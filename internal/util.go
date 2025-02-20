package internal

import (
	"fmt"
	"time"

	"github.com/anthdm/hollywood/actor"
)

func Request[Result any](engine *actor.Engine, pid *actor.PID, msg any) (result Result, err error) {
	maybeResult, err := engine.Request(pid, msg, time.Second).Result()
	if err != nil {
		return result, err
	}
	result, ok := maybeResult.(Result)
	if !ok {
		return result, fmt.Errorf("result was of incorrect type: %T", maybeResult)
	}
	return result, nil
}
