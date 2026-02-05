package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/message-broker"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	*kafka.Consumer

	config *kafka.ConfigMap

	channelIDs []string
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
						}
					case <-ctx.Done():
						return
					}
					continue
				}

				if kafkaErr, ok := readErr.(kafka.Error); ok {
					if kafkaErr.IsFatal() {

						if recoverErr := c.recover(); recoverErr != nil {
							return
						}

					} else if !kafkaErr.IsTimeout() {

						time.Sleep(5 * time.Second)
					}

				} else {

					time.Sleep(5 * time.Second)

				}
			}
		}
	}()

	return ch, nil
}

func (c *Consumer) Close() {
	c.Consumer.Close()
}

func (c *Consumer) subscribe() error {
	return c.SubscribeTopics(c.channelIDs, nil)
}

func (c *Consumer) recover() error {

	for _, backoff := range []time.Duration{5, 15, 30} {
		time.Sleep(backoff * time.Second)

		if c.Consumer != nil {
			c.Consumer.Close()
		}

		if err := c.build(); err == nil {
			return nil
		} else {
			continue
		}
	}

	return fmt.Errorf("failed to recover kafka consumer after 3 attempts")
}

func (c *Consumer) build() (err error) {
	if c.Consumer, err = kafka.NewConsumer(c.config); err != nil {
		return
	}

	if len(c.channelIDs) > 0 {
		c.subscribe()
	}

	return
}
