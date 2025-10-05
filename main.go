package main

import (
	"context"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	pb "github.com/ave1995/grpc-chat/api/grpc/proto"
	"github.com/ave1995/grpc-chat/api/grpc/server"
	"github.com/ave1995/grpc-chat/config"
	"github.com/ave1995/grpc-chat/service/message"
	gormdb "github.com/ave1995/grpc-chat/store/gormdb"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		os.Exit(1)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	mainContext := context.Background()

	gorm, err := gormdb.NewGormConnection(mainContext, cfg.DBConfig())
	if err != nil {
		log.Fatalf("failed to ini database connection %v", err)
		os.Exit(1)
	}

	messageStore := gormdb.NewMessageStore(gorm)

	messageService := message.NewMessageService(messageStore)

	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, server.NewChatServer(messageService))

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
