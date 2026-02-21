package increase

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
	groupID           = "inventory-increase"
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
	h = &Handler{
		retryDelay:        defaultRetryDelay,
		maxRetries:        maxRetries,
		processingTimeout: processingTimeout,
		log:               log,
		tracer:            tracer.NewTracer("async-increase-inventory"),
		inventoryDAO:      inventoryDAO,
		consumableDAO:     consumableDAO,
	}

	if h.consumer, err = messageBroker.NewConsumer(groupID); err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	if h.producer, err = messageBroker.NewProducer(); err != nil {
		h.consumer.Close()
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return h, nil
}

func (h *Handler) Handle(ctx context.Context) (err error) {
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		done := make(chan struct{})
		go func() {
			h.wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-shutdownCtx.Done():
		}

		h.shutdownOnce.Do(func() {
			if h.consumer != nil {
				h.consumer.Close()
			}
			if h.producer != nil {
				h.producer.Close()
			}
		})
	}()

	err = h.listenQueue(ctx)
	return
}

func (h *Handler) listenQueue(ctx context.Context) (err error) {
	ch, err := h.consumer.Consume(ctx, channels.InventoryIncreased, channels.InventoryIncreasedRetry)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case message, ok := <-ch:
			if !ok {
				return nil
			}

			h.wg.Add(1)
			h.handleMessage(ctx, message)
		}
	}
}

func (h *Handler) handleMessage(ctx context.Context, message mPorts.Message) {
	defer h.wg.Done()

	procCtx, cancel := context.WithTimeout(ctx, h.processingTimeout)
	defer cancel()

	err := h.processEvent(procCtx, message)
	if err != nil {
		h.failedCount.Add(1)
		h.log.Error("failed to process increase message", err, map[string]any{
			"channel": message.ChannelID,
		})

		if message.ChannelID != channels.InventoryIncreasedRetry {
			h.scheduleRetry(ctx, message, err)
		}
		return
	}

	h.processedCount.Add(1)
}

func (h *Handler) processEvent(ctx context.Context, message mPorts.Message) (err error) {
	span := h.tracer.GetSpanFromContext(ctx)
	var consumableID string

	consumableID, ok := message.Value.(string)
	if !ok {
		err = errors.New("invalid message value, expected string")
		span.AddEvent("processing-increase-event-invalid-message", tracing.Attributes{
			"error": err,
		})
		return err
	}

	consumable, err := h.consumableDAO.Get(ctx, consumableID)
	if err != nil {
		span.AddEvent("processing-increase-event-consumable-get-failed", tracing.Attributes{
			"error": err,
		})
		return err
	}

	err = h.inventoryDAO.Increase(ctx, iPorts.InventoryDAO{
		ConsumableID: consumable.ID,
	})
	if err != nil {
		span.AddEvent("processing-increase-event-inventory-increase-failed", tracing.Attributes{
			"error": err,
		})
		return err
	}

	span.AddEvent("processing-increase-event-success", nil)
	return nil
}

func (h *Handler) scheduleRetry(ctx context.Context, message mPorts.Message, originalErr error) {
	h.retriedCount.Add(1)

	select {
	case <-ctx.Done():
		return
	case <-time.After(h.retryDelay):
	}

	h.sendToRetryQueue(ctx, message, originalErr)
}

func (h *Handler) sendToRetryQueue(ctx context.Context, message mPorts.Message, originalErr error) {
	if message.Headers == nil {
		message.Headers = make(map[string]string)
	}
	message.Headers["original_error"] = originalErr.Error()
	message.Headers["retry_timestamp"] = time.Now().Format(time.RFC3339)

	retryCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_ = h.producer.Produce(retryCtx, channels.InventoryIncreasedRetry, message)
}
