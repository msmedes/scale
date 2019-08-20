test:
	go test ./... -cover

deps:
	dep ensure

serve:
	go run cmd/scale/main.go
