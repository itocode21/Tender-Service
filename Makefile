# Variables
APP_NAME := tender-service
DOCKER_COMPOSE := docker compose

# Targets
.PHONY: build
build:
	$(DOCKER_COMPOSE) build

.PHONY: up
up:
	$(DOCKER_COMPOSE) up --build

.PHONY: down
down:
	$(DOCKER_COMPOSE) down

.PHONY: down-v
down-v:
	$(DOCKER_COMPOSE) down -v

.PHONY: network-prune
network-prune:
	docker network prune

.PHONY: run
run:
	go run cmd/main.go

.PHONY: test
test:
	go test ./internal/...

.PHONY: format
format:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run ./...