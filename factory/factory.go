package factory

import (
	"log/slog"
	"sync"

	"github.com/ave1995/grpc-chat/config"
)

type Factory struct {
	config config.Config

	logger     *slog.Logger
	loggerOnce sync.Once
}

func NewFactory(config config.Config) *Factory {
	return &Factory{
		config: config,
	}
}

func (f *Factory) Logger() *slog.Logger {
	f.loggerOnce.Do(func() {
		f.logger = newLogger()
	})

	return f.logger
}
