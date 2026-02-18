package messagebroker

import "context"

type MessageBroker interface {
	NewConsumer(groupID string) (Consumer, error)
	NewProducer() (Producer, error)
}

type Consumer interface {
	Consume(ctx context.Context, channelID ...string) (channel <-chan Message, err error)
	Close()
}

type Producer interface {
	Produce(ctx context.Context, channelID string, message Message) error
	Close()
}

type Message struct {
	ChannelID string
	Headers   map[string]string
	Value     any
	Error     string
}
