package main

import (
	"log"
	"net"

	"github.com/ave1995/grpc-chat/api/grpc/service"
	"google.golang.org/grpc"

	pb "github.com/ave1995/grpc-chat/api/grpc/proto"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, service.NewChatServer())

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
