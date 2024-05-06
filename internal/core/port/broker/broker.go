package broker

import "context"

type Producer interface {
	ProduceMessage(ctx context.Context, topic string, key []byte, value []byte) error
	Close() error
}

type Consumer interface {
	ConsumeMessages(ctx context.Context, topic string, handler func(key []byte, value []byte) error) error
	Close() error
}
