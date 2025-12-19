LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

generate:
	make generate-bidding-api
	make generate-auction-api

generate-bidding-api:
	mkdir -p internal/pb/bidding_api
	mkdir -p api/swagger
	PATH=$(LOCAL_BIN):$$PATH protoc --proto_path api/bidding_api \
	--proto_path api \
	--proto_path third_party \
	--go_out=internal/pb/bidding_api --go_opt=paths=source_relative \
	--go-grpc_out=internal/pb/bidding_api --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=internal/pb/bidding_api --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=api/swagger --openapiv2_opt=allow_merge=true,merge_file_name=bidding_service \
	api/bidding_api/bidding.proto

generate-auction-api:
	mkdir -p internal/pb/auction_api
	mkdir -p api/swagger
	PATH=$(LOCAL_BIN):$$PATH protoc --proto_path api/auction_api \
	--proto_path api \
	--proto_path third_party \
	--go_out=internal/pb/auction_api --go_opt=paths=source_relative \
	--go-grpc_out=internal/pb/auction_api --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=internal/pb/auction_api --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=api/swagger --openapiv2_opt=allow_merge=true,merge_file_name=auction_service \
	api/auction_api/auction.proto

run-bidding:
	go run cmd/app/main.go

run-auction:
	go run cmd/auction_svc/main.go

run-gateway:
	go run cmd/gateway/main.go

run-notification:
	go run cmd/notification_svc/main.go