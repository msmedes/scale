.DEFAULT_GOAL := default

default: test lint

test:
	@go test ./... -cover -v

lint:
	@$(shell go list -f {{.Target}} golang.org/x/lint/golint) ./...

serve:
	@go run cmd/scale/main.go

codegen:
	@protoc -I internal/pkg/rpc internal/pkg/rpc/proto/scale.proto --go_out=plugins=grpc:internal/pkg/rpc

docker.build:
	@docker build -t msmedes/scale:dev .

docker.run:
	@docker run -p 3000:3000 msmedes/scale:dev

docker: docker.build docker.run
