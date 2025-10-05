generate-proto:
	protoc --go_out=. --go-grpc_out=. api/grpc/proto/chat.proto