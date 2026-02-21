package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	ports "github.com/Logstay/tech-challange-execution-service/internal/ports/message-broker"
	"github.com/Logstay/tech-challange-execution-service/pkg/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	*kafka.Consumer

	config *kafka.ConfigMap

	channelIDs []string
	logger     logger.Logger
}

func (c *Consumer) Consume(ctx context.Context, channelIDs ...string) (channel <-chan ports.Message, err error) {
	c.channelIDs = channelIDs

	if err = c.subscribe(); err != nil {

		return
	}

	ch := make(chan ports.Message, 100)

	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				kafkaMessage, readErr := c.ReadMessage(-1)

				if readErr == nil {
					var message ports.Message

					if err := json.Unmarshal(kafkaMessage.Value, &message); err != nil {
						continue
					}

					select {
					case ch <- message:
						if _, commitErr := c.CommitMessage(kafkaMessage); commitErr != nil {
							c.logger.Error("Failed to commit message", commitErr)
						}
					case <-ctx.Done():
						return
					}
					continue
				}

				c.handleReadError(readErr)
			}
		}
	}()

	return ch, nil
}

func (c *Consumer) Close() {
	if closeErr := c.Consumer.Close(); closeErr != nil {
		c.logger.Error("Failed to close consumer", closeErr)
	}
}

func (c *Consumer) subscribe() error {
	return c.SubscribeTopics(c.channelIDs, nil)
}

func (c *Consumer) handleReadError(readErr error) {
	var kafkaErr kafka.Error
	if !errors.As(readErr, &kafkaErr) {
		time.Sleep(5 * time.Second)
		return
	}

	if kafkaErr.IsFatal() {
		if recoverErr := c.recover(); recoverErr != nil {
			return
		}
	} else if !kafkaErr.IsTimeout() {
		time.Sleep(5 * time.Second)
	}
}

func (c *Consumer) recover() error {
	for _, backoff := range []time.Duration{5 * time.Second, 15 * time.Second, 30 * time.Second} {
		time.Sleep(backoff)

		if c.Consumer != nil {
			if closeErr := c.Consumer.Close(); closeErr != nil {
				c.logger.Error("Failed to close consumer during recovery", closeErr)
			}
		}

		if err := c.build(); err == nil {
			return nil
		}
	}

	return errors.New("failed to recover kafka consumer after 3 attempts")
}

func (c *Consumer) build() (err error) {
	if c.Consumer, err = kafka.NewConsumer(c.config); err != nil {
		return
	}

	if len(c.channelIDs) > 0 {
		if subErr := c.subscribe(); subErr != nil {
			c.logger.Error("Failed to subscribe to channels", subErr)
		}
	}

	return
}
