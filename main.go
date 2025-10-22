package main

import (
	"context"
	"net"
	"os"

	"google.golang.org/grpc"

	pb "github.com/ave1995/grpc-chat/api/grpc/proto"
	"github.com/ave1995/grpc-chat/api/grpc/server"
	"github.com/ave1995/grpc-chat/config"
	"github.com/ave1995/grpc-chat/connector"
	"github.com/ave1995/grpc-chat/connector/kafka"
	"github.com/ave1995/grpc-chat/factory"
	"github.com/ave1995/grpc-chat/service/message"
	"github.com/ave1995/grpc-chat/store/gormdb"
	"github.com/ave1995/grpc-chat/utils"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	factory := factory.NewFactory(cfg)
	logger := factory.Logger()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("listen: ", utils.SlogError(err))
		os.Exit(1)
	}

	mainContext := context.Background()

	gorm, err := gormdb.NewGormConnection(mainContext, cfg.DBConfig())
	if err != nil {
		logger.Error("ini database connection: ", utils.SlogError(err))
		os.Exit(1)
	}

	messageStore := gormdb.NewMessageStore(gorm)

	producer := kafka.NewKafkaProducer(cfg.KafkaConfig())

	hub := connector.NewMessageHub(mainContext, logger, 10)

	// TODO: make topic configurable
	messageService := message.NewMessageService(messageStore, producer, "messages", hub)

	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, server.NewChatServer(logger, messageService))

	logger.Info("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("to serve: ", utils.SlogError(err))
		os.Exit(1)
	}
}
