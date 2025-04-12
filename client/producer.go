package client

import (
	"log"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/troygilman/actormq/cluster"
)

type ProducerConfig struct {
	Topic      string
	Serializer remote.Serializer
}

type producerActor struct {
	config      ProducerConfig
	pods        []*actor.PID
	leader      *actor.PID
	requests    []produceRequestMetadata
	nextOffset  uint64
	ackedOffset uint64
}

func NewProducer(config ProducerConfig, pods []*actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &producerActor{
			config: config,
			pods:   pods,
			leader: pods[0],
		}
	}
}

func (producer *producerActor) Receive(act *actor.Context) {
	switch msg := act.Message().(type) {
	case ProduceMessage:
		producer.handleProduceMessage(act, msg)
	case *cluster.EnvelopeResult:
		producer.handleEnvelopeResult(act, msg)
	}
}

func (producer *producerActor) handleProduceMessage(act *actor.Context, msg ProduceMessage) {
	start := time.Now()

	data, err := producer.config.Serializer.Serialize(msg.Message)
	if err != nil {
		panic(err)
	}

	envelope := &cluster.Envelope{
		Offset: producer.nextOffset,
		Topic:  producer.config.Topic,
		Message: &cluster.Message{
			TypeName: producer.config.Serializer.TypeName(msg.Message),
			Data:     data,
		},
	}

	act.Send(producer.leader, envelope)
	producer.requests = append(producer.requests, produceRequestMetadata{
		sender:   act.Sender(),
		offset:   producer.nextOffset,
		start:    start,
		envelope: envelope,
	})
}

func (producer *producerActor) handleEnvelopeResult(act *actor.Context, msg *cluster.EnvelopeResult) {
	index := msg.Offset - producer.ackedOffset
	requestMetadata := producer.requests[index]

	if !msg.Success {
		if msg.RedirectPID != nil {
			log.Println("redirecting msg - offset:", msg.Offset)
			producer.leader = cluster.PIDToActorPID(msg.RedirectPID)
			for _, request := range producer.requests[:index+1] {
				act.Send(producer.leader, request.envelope)
			}
			return
		}
		log.Println("cannot send msg: ", msg.Error)
		return
	}

	producer.requests = producer.requests[index+1:]
	producer.ackedOffset = msg.Offset
	duration := time.Since(requestMetadata.start)
	log.Printf("Produced message in %d ms\n", duration.Milliseconds())
	act.Send(requestMetadata.sender, ProduceMessageResult{})
}

type produceRequestMetadata struct {
	sender   *actor.PID
	offset   uint64
	start    time.Time
	envelope *cluster.Envelope
}
