.DEFAULT_GOAL := default

default: test

test:
	go test ./... -cover

serve:
	go run cmd/scale/main.go

codegen:
	protoc -I internal/app/scale internal/app/scale/proto/scale.proto --go_out=plugins=grpc:internal/app/scale

docker.build:
	docker build -t msmedes/scale:dev .

docker.run:
	docker run -p 3000:3000 msmedes/scale:dev

docker: docker.build docker.run
