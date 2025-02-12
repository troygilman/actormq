package client

import "github.com/anthdm/hollywood/actor"

type CreateConsumer struct {
	ConsumerConfig ConsumerConfig
}

type CreateConsumerResult struct {
	PID *actor.PID
}

type CreateProducer struct {
	ProducerConfig ProducerConfig
}

type CreateProducerResult struct {
	PID *actor.PID
}

type ProduceMessage struct {
	Message any
}

type ConsumeMessage struct {
	Message any
}
