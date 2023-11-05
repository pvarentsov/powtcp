.DEFAULT_GOAL := help

help:
	@echo
	@echo "Usage: make [command]"
	@echo
	@echo "Commands:"
	@echo
	@echo " build-server          Build server app"	
	@echo " build-client          Build client app"
	@echo	
	@echo " run-server            Run server app"
	@echo " run-client            Run client app"
	@echo
	@echo " test                  Run tests"
	@echo " fmt                   Format code"
	@echo

build-server:
	@go build -o ./bin/server ./cmd/server/*.go

build-client:
	@go build -o ./bin/client ./cmd/client/*.go

run-server:
	@./bin/server --config ./config/template.yaml

run-client:
	@./bin/client --config ./config/template.yaml

test:
	@go test ./... -v

fmt:
	@go fmt ./...