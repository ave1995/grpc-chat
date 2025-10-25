package message

import (
	"context"
	"log/slog"
	"time"

	"github.com/ave1995/grpc-chat/config"
	"github.com/ave1995/grpc-chat/domain/connector"
	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/domain/store"
)

type Processor struct {
	logger   *slog.Logger
	config   config.MessageProcessorConfig
	store    store.OutboxStore
	producer connector.Producer
}

func NewProcessor(logger *slog.Logger, config config.MessageProcessorConfig, store store.OutboxStore, producer connector.Producer) *Processor {
	return &Processor{
		logger:   logger,
		config:   config,
		store:    store,
		producer: producer,
	}
}

func (p *Processor) processPending(ctx context.Context) error {
	events, err := p.store.GetPendingEvents(ctx, model.SendMessage, p.config.OutboxBatchSize)
	if err != nil {
		p.logger.Error("failed to get pending events", "error", err)
		return err
	}

	for _, event := range events {
		if err := p.producer.Send(ctx, p.config.Topic, event.AggregateID.String(), event.Payload); err != nil {
			p.logger.Error("failed to send message", "error", err, "id", event.ID)
			err := p.store.MarkFailed(ctx, event.ID)
			if err != nil {
				// TODO: how to approach this
				p.logger.Error("failed to mark failed", "error", err, "id", event.ID)
				return err
			}
			continue
		}

		if err := p.store.MarkProcessed(ctx, event.ID); err != nil {
			// TODO: how to approach this
			p.logger.Error("failed to mark processed", "error", err, "id", event.ID)
			return err
		}
	}

	return nil
}

func (p *Processor) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(p.config.OutboxInterval)
		defer ticker.Stop()

		p.logger.Info("message processor started")

		for {
			select {
			case <-ctx.Done():
				p.logger.Info("message processor shutting down gracefully")
				return

			case <-ticker.C:
				if err := p.processPending(ctx); err != nil {
					p.logger.Error("failed to process pending events", "error", err)
				}
			}
		}
	}()
}
