package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	pb "github.com/ave1995/grpc-chat/api/grpc/proto"
	"github.com/ave1995/grpc-chat/api/grpc/server"
	"github.com/ave1995/grpc-chat/config"
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
		logger.Error("listen: ", utils.SlogError(err))
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, server.NewChatServer(logger, factory.MessageService()))

	go func() {
		logger.Info("gRPC server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("serve: ", utils.SlogError(err))
			os.Exit(1)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh
	logger.Info(fmt.Sprintf("Received signal: %v. Shutting down gracefully...", sig))

	grpcServer.GracefulStop()

	factory.Close()

	cancel()

	logger.Info("Server stopped cleanly.")
}
