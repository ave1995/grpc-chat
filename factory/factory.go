package factory

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/ave1995/grpc-chat/config"
	"github.com/ave1995/grpc-chat/connector/kafka"
	"github.com/ave1995/grpc-chat/domain/connector"
	"github.com/ave1995/grpc-chat/domain/service"
	"github.com/ave1995/grpc-chat/domain/store"
	"github.com/ave1995/grpc-chat/service/message"
	"github.com/ave1995/grpc-chat/store/gormdb"
	"github.com/ave1995/grpc-chat/utils"
	"gorm.io/gorm"
)

type Factory struct {
	context context.Context
	config  config.Config

	logger     *slog.Logger
	loggerOnce sync.Once

	db     *gorm.DB
	dbOnce sync.Once

	messageStore     *gormdb.MessageStore
	messageStoreOnce sync.Once

	kafkaProducer     *kafka.Producer
	kafkaProducerOnce sync.Once

	hub     *message.Hub
	hubOnce sync.Once

	messageService     *message.Service
	messageServiceOnce sync.Once
}

func NewFactory(context context.Context, config config.Config) *Factory {
	return &Factory{
		context: context,
		config:  config,
	}
}

func (f *Factory) Logger() *slog.Logger {
	f.loggerOnce.Do(func() {
		f.logger = newLogger()
	})

	return f.logger
}

func (f *Factory) Database() *gorm.DB {
	f.dbOnce.Do(func() {
		ctxWithTimeout, cancel := context.WithTimeout(f.context, 5*time.Second)
		defer cancel()

		var err error
		f.db, err = gormdb.NewGormConnection(ctxWithTimeout, f.config.DBConfig())
		if err != nil {
			f.logger.Error("ini database connection: ", utils.SlogError(err))
			os.Exit(1)
		}
	})

	return f.db
}

func (f *Factory) MessageStore() store.MessageStore {
	f.messageStoreOnce.Do(func() {
		f.messageStore = gormdb.NewMessageStore(f.Database())
	})

	return f.messageStore
}

func (f *Factory) KafkaProducer() connector.Producer {
	f.kafkaProducerOnce.Do(func() {
		f.kafkaProducer = kafka.NewKafkaProducer(f.config.KafkaConfig())
	})

	return f.kafkaProducer
}

func (f *Factory) Hub() *message.Hub {
	f.hubOnce.Do(func() {
		f.hub = message.NewHub(f.context, f.Logger(), f.config.HubCapacity)
	})

	return f.hub
}

func (f *Factory) MessageService() service.MessageService {
	f.messageServiceOnce.Do(func() {
		f.messageService = message.NewService(
			f.config.MessageServiceConfig(),
			f.MessageStore(),
			f.Hub(),
		)
	})

	return f.messageService
}

func (f *Factory) Close() {
	logger := f.Logger()
	logger.Info("shutting down factory components...")

	if f.kafkaProducer != nil {
		err := f.kafkaProducer.Close()
		if err != nil {
			logger.Error("kafka producer close: ", utils.SlogError(err))
			return
		}
	}

	if f.db != nil {
		sqlDB, err := f.db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				logger.Error("close database", slog.String("error", err.Error()))
			} else {
				logger.Info("database connection closed")
			}
		}
	}

	logger.Info("factory shutdown complete.")
}
