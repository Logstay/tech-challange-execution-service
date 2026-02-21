package generatedataset

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	channels "github.com/Logstay/tech-challange-execution-service/internal/events"
	cPorts "github.com/Logstay/tech-challange-execution-service/internal/ports/consumable"
	iPorts "github.com/Logstay/tech-challange-execution-service/internal/ports/inventory"
	mPorts "github.com/Logstay/tech-challange-execution-service/internal/ports/message-broker"
	"github.com/Logstay/tech-challange-execution-service/pkg/logger"
	"github.com/Logstay/tech-challange-execution-service/pkg/tracing"
)

const (
	groupID           = "inventory"
	defaultRetryDelay = time.Minute
	maxRetries        = 3
	processingTimeout = 30 * time.Second
	shutdownTimeout   = 30 * time.Second
)

type Handler struct {
	retryDelay        time.Duration
	maxRetries        int
	processingTimeout time.Duration
	log               logger.Logger
	consumer          mPorts.Consumer
	producer          mPorts.Producer
	tracer            *tracing.Tracer
	inventoryDAO      iPorts.DAO
	consumableDAO     cPorts.DAO

	processedCount atomic.Int64
	failedCount    atomic.Int64
	retriedCount   atomic.Int64

	wg           sync.WaitGroup
	shutdownOnce sync.Once
}

func New(messageBroker mPorts.MessageBroker, tracer *tracing.Provider, log logger.Logger, inventoryDAO iPorts.DAO, consumableDAO cPorts.DAO) (h *Handler, err error) {
	log.Info("Initializing inventory handler", map[string]any{
		"group_id":     groupID,
		"retry_delay":  defaultRetryDelay,
		"max_retries":  maxRetries,
		"proc_timeout": processingTimeout,
	})

	h = &Handler{
		retryDelay:        defaultRetryDelay,
		maxRetries:        maxRetries,
		processingTimeout: processingTimeout,
		log:               log,
		tracer:            tracer.NewTracer("async-decrease-inventory"),
		inventoryDAO:      inventoryDAO,
		consumableDAO:     consumableDAO,
	}

	if h.consumer, err = messageBroker.NewConsumer(groupID); err != nil {
		log.Error("Failed to create consumer", err, map[string]any{
			"group_id": groupID,
		})
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	if h.producer, err = messageBroker.NewProducer(); err != nil {
		log.Error("Failed to create producer", err)
		h.consumer.Close()
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	log.Info("Inventory handler initialized successfully", nil)
	return h, nil
}

func (h *Handler) Handle(ctx context.Context) (err error) {
	h.log.Info("Starting event handler", nil)

	defer func() {
		h.log.Info("Shutting down event handler", nil)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		done := make(chan struct{})
		go func() {
			h.wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			h.log.Info("All messages processed, closing gracefully", nil)
		case <-shutdownCtx.Done():
			h.log.Warn("Shutdown timeout reached, forcing close", nil)
		}

		h.shutdownOnce.Do(func() {
			if h.consumer != nil {
				h.consumer.Close()
			}
			if h.producer != nil {
				h.producer.Close()
			}
		})

		h.log.Info("Handler metrics", map[string]any{
			"processed": h.processedCount.Load(),
			"failed":    h.failedCount.Load(),
			"retried":   h.retriedCount.Load(),
		})
	}()

	err = h.listenQueue(ctx)

	return
}

func (h *Handler) listenQueue(ctx context.Context) (err error) {
	ch, err := h.consumer.Consume(ctx, channels.InventoryDecreased, channels.InventoryDecreasedRetry)

	if err != nil {
		h.log.Error("Failed to start consuming", err)
		return fmt.Errorf("failed to consume: %w", err)
	}

	h.log.Info("Started listening to queue", map[string]any{
		"channels": []string{channels.InventoryDecreased, channels.InventoryDecreasedRetry},
	})

	for {
		select {
		case <-ctx.Done():
			h.log.Info("Context cancelled, stopping queue listener", nil)
			return ctx.Err()

		case message, ok := <-ch:
			if !ok {
				h.log.Warn("Consumer channel closed", nil)
				return nil
			}

			h.wg.Add(1)

			h.handleMessage(ctx, message)
		}
	}
}

func (h *Handler) handleMessage(ctx context.Context, message mPorts.Message) {
	defer h.wg.Done()

	start := time.Now()
	h.log.Debug("Processing message", map[string]any{
		"channel": message.ChannelID,
	})

	procCtx, cancel := context.WithTimeout(ctx, h.processingTimeout)
	defer cancel()

	err := h.processEvent(procCtx, message)

	duration := time.Since(start)

	if err != nil {
		h.failedCount.Add(1)
		h.log.Error("Failed to process message", err, map[string]any{
			"channel":     message.ChannelID,
			"duration_ms": duration.Milliseconds(),
		})

		if message.ChannelID != channels.InventoryDecreasedRetry {
			h.scheduleRetry(ctx, message, err)
		} else {
			h.log.Error("Message failed after retry, giving up", err)
		}
		return
	}

	h.processedCount.Add(1)
	h.log.Info("Message processed successfully", map[string]any{
		"channel":     message.ChannelID,
		"duration_ms": duration.Milliseconds(),
	})
}

func (h *Handler) processEvent(ctx context.Context, message mPorts.Message) (err error) {
	span := h.tracer.GetSpanFromContext(ctx)
	var consumableID string

	span.AddEvent("processing-event-start", tracing.Attributes{
		"channel": message.ChannelID,
	})

	consumableID, ok := message.Value.(string)
	if !ok {
		err = errors.New("invalid message value, expected string but got")
		span.AddEvent("processing-event-invalid-message", tracing.Attributes{
			"error": err,
		})
		return err
	}

	consumable, err := h.consumableDAO.Get(ctx, consumableID)
	if err != nil {
		h.log.Error("Failed to get consumable", err)
		span.AddEvent("processing-event-consumable-get-failed", tracing.Attributes{
			"error": err,
		})
		return err
	}

	err = h.inventoryDAO.Decrease(ctx, iPorts.InventoryDAO{
		ConsumableID: consumable.ID,
	})
	if err != nil {
		h.log.Error("Failed to decrease inventory", err)
		span.AddEvent("processing-event-inventory-decrease-failed", tracing.Attributes{
			"error": err,
		})
		return err
	}

	span.AddEvent("processing-event-success", nil)

	return nil
}

func (h *Handler) scheduleRetry(ctx context.Context, message mPorts.Message, originalErr error) {
	h.log.Info("Scheduling message retry", map[string]any{
		"retry_delay": h.retryDelay,
	})

	h.retriedCount.Add(1)

	select {
	case <-ctx.Done():
		h.log.Warn("Context cancelled, skipping retry")
		return
	case <-time.After(h.retryDelay):
	}

	h.sendToRetryQueue(ctx, message, originalErr)
}

func (h *Handler) sendToRetryQueue(ctx context.Context, message mPorts.Message, originalErr error) {
	span := h.tracer.GetSpanFromContext(ctx)
	span.AddEvent("sending-to-retry-queue", tracing.Attributes{
		"original_error": originalErr.Error(),
	})

	if message.Headers == nil {
		message.Headers = make(map[string]string)
	}
	message.Headers["original_error"] = originalErr.Error()
	message.Headers["retry_timestamp"] = time.Now().Format(time.RFC3339)

	retryCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := h.producer.Produce(retryCtx, channels.InventoryDecreasedRetry, message)
	if err != nil {
		h.log.Error("Failed to send message to retry queue", err)
		span.AddEvent("retry-queue-send-failed", tracing.Attributes{
			"error": err.Error(),
		})
		return
	}

	h.log.Info("Message sent to retry queue")
	span.AddEvent("retry-queue-send-success", nil)
}
