generate-proto:
	protoc -I. \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/grpc/proto/chat.proto
		
up: 
	docker-compose up -d

down:
	docker-compose down

clean-go:
	go mod tidy
