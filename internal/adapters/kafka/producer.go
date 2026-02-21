package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	ports "github.com/Logstay/tech-challange-execution-service/internal/ports/message-broker"
	"github.com/Logstay/tech-challange-execution-service/pkg/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

type Producer struct {
	*kafka.Producer
	logger logger.Logger
}

func (p Producer) Produce(ctx context.Context, channelID string, message ports.Message) (err error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	key := uuid.NewString()
	message.ChannelID = channelID

	raw, err := json.Marshal(message)
	if err != nil {
		p.logger.Error("Failed to marshal message", err, map[string]any{
			"topic":   channelID,
			"message": message,
		})
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	deliveryChan := make(chan kafka.Event, 1)
	defer close(deliveryChan)

	err = p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &channelID,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(key),
		Value: raw,
	}, deliveryChan)

	if err != nil {
		p.logger.Error("Failed to produce message", err, map[string]any{
			"topic": channelID,
			"key":   key,
		})
		return fmt.Errorf("failed to produce message: %w", err)
	}

	select {
	case e := <-deliveryChan:
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			p.logger.Error("Message delivery failed", m.TopicPartition.Error, map[string]any{
				"topic":     channelID,
				"partition": m.TopicPartition.Partition,
				"offset":    m.TopicPartition.Offset,
			})
			return fmt.Errorf("message delivery failed: %w", m.TopicPartition.Error)
		}
		p.logger.Debug("Message delivered successfully", map[string]any{
			"topic":     channelID,
			"partition": m.TopicPartition.Partition,
			"offset":    m.TopicPartition.Offset,
			"key":       key,
		})
	case <-ctx.Done():
		p.logger.Warn("Context cancelled while waiting for delivery confirmation", map[string]any{
			"topic": channelID,
			"key":   key,
		})
		return ctx.Err()
	case <-time.After(30 * time.Second):
		p.logger.Error("Timeout waiting for delivery confirmation", errors.New("timeout after 30s"), map[string]any{
			"topic": channelID,
			"key":   key,
		})
		return errors.New("timeout waiting for message delivery confirmation")
	}

	return nil
}

func (p Producer) Close() {
	p.logger.Info("Closing Kafka producer", nil)

	remaining := p.Flush(30000)
	if remaining > 0 {
		p.logger.Warn("Producer closed with messages still in queue", map[string]any{
			"remaining_messages": remaining,
		})
	}

	p.Producer.Close()
	p.logger.Info("Kafka producer closed", nil)
}
