package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	pb "github.com/ave1995/grpc-chat/api/grpc/proto"
	"github.com/ave1995/grpc-chat/api/grpc/server"
	"github.com/ave1995/grpc-chat/config"
	"github.com/ave1995/grpc-chat/connector/kafka"
	"github.com/ave1995/grpc-chat/factory"
	"github.com/ave1995/grpc-chat/utils"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	factory := factory.NewFactory(ctx, cfg)
	logger := factory.Logger()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("net.listen", utils.SlogError(err))
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, server.NewChatServer(logger, factory.MessageService()))

	consumer := kafka.NewKafkaConsumer(logger, cfg.KafkaConfig().Brokers, "messages", "", factory.Hub())
	consumer.Start(ctx)

	go func() {
		logger.Info("gRPC server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("grpcServer.serve", utils.SlogError(err))
			os.Exit(1)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh
	logger.Info("Received signal. Shutting down gracefully...", "signal", sig)

	grpcServer.GracefulStop()

	if err := consumer.Stop(); err != nil {
		logger.Error("consumer.stop", utils.SlogError(err))
	}

	factory.Close()

	cancel()

	logger.Info("Server stopped cleanly.")
}
