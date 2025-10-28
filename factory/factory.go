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
	"github.com/ave1995/grpc-chat/domain/store"
	"github.com/ave1995/grpc-chat/service/message"
	"github.com/ave1995/grpc-chat/store/cached"
	"github.com/ave1995/grpc-chat/store/gormdb"
	"github.com/ave1995/grpc-chat/store/memory"
	"github.com/ave1995/grpc-chat/store/redis"
	"github.com/ave1995/grpc-chat/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Factory struct {
	context context.Context
	config  config.Config

	logger     *slog.Logger
	loggerOnce sync.Once

	db     *gorm.DB
	dbOnce sync.Once

	memoryCache     *memory.Cache
	memoryCacheOnce sync.Once

	redisCache     *redis.Cache
	redisCacheOnce sync.Once

	messageStore     *gormdb.MessageStore
	messageStoreOnce sync.Once

	cachedMessageStore     *cached.MessageStore
	cachedMessageStoreOnce sync.Once

	kafkaProducer     *kafka.Producer
	kafkaProducerOnce sync.Once

	kafkaConsumer     *kafka.Consumer
	kafkaConsumerOnce sync.Once

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

func (f *Factory) MemoryCache() *memory.Cache {
	f.memoryCacheOnce.Do(func() {
		// TODO: config
		f.memoryCache = memory.NewCache(10)
	})
	return f.memoryCache
}

func (f *Factory) RedisCache() *redis.Cache {
	f.redisCacheOnce.Do(func() {
		f.redisCache = redis.NewCache("localhost:6379", "", 0)
	})
	return f.redisCache
}

func (f *Factory) MessageStore() store.MessageStore {
	f.messageStoreOnce.Do(func() {
		f.messageStore = gormdb.NewMessageStore(f.Database())
	})

	return f.messageStore
}

func (f *Factory) CachedMessageStore() store.MessageStore {
	f.cachedMessageStoreOnce.Do(func() {
		f.cachedMessageStore = cached.NewMessageStore(
			f.Logger(),
			f.MessageStore(),
			//f.MemoryCache(),
			f.RedisCache(),
		)
	})

	return f.cachedMessageStore
}

func (f *Factory) KafkaProducer() connector.Producer {
	f.kafkaProducerOnce.Do(func() {
		f.kafkaProducer = kafka.NewKafkaProducer(f.Logger(), f.config.KafkaConfig())
	})

	return f.kafkaProducer
}

func (f *Factory) KafkaConsumer() connector.Consumer {
	f.kafkaConsumerOnce.Do(func() {
		f.kafkaConsumer = kafka.NewKafkaConsumer(f.Logger(), f.config.KafkaConfig().Brokers, f.config.MessageTopic, uuid.New().String())
	})

	return f.kafkaConsumer
}

func (f *Factory) Hub() *message.Hub {
	f.hubOnce.Do(func() {
		f.hub = message.NewHub(f.context, f.Logger(), f.config.HubCapacity)
	})

	return f.hub
}

func (f *Factory) MessageService() *message.Service {
	f.messageServiceOnce.Do(func() {
		f.messageService = message.NewService(
			f.Logger(),
			f.config.MessageServiceConfig(),
			//f.MessageStore(),
			f.CachedMessageStore(),
			f.Hub(),
			f.KafkaConsumer(),
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

	if f.redisCache != nil {
		err := f.redisCache.Close()
		if err != nil {
			logger.Error("close redis cache: ", utils.SlogError(err))
		}
	}

	if f.memoryCache != nil {
		err := f.memoryCache.Close()
		if err != nil {
			logger.Error("close memory cache: ", utils.SlogError(err))
		}
	}

	logger.Info("factory shutdown complete.")
}
