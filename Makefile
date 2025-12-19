LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-bidding-api

generate-bidding-api:
	mkdir -p internal/pb/bidding_api
	PATH=$(LOCAL_BIN):$$PATH protoc --proto_path api/bidding_api \
	--proto_path api \
	--go_out=internal/pb/bidding_api --go_opt=paths=source_relative \
	--go-grpc_out=internal/pb/bidding_api --go-grpc_opt=paths=source_relative \
	api/bidding_api/bidding.proto

run:
	go run cmd/app/main.go