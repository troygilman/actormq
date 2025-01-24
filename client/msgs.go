package client

import "github.com/anthdm/hollywood/actor"

type CreateConsumer struct {
	Topic string
}

type CreateConsumerResult struct {
	PID   *actor.PID
	Error error
}
