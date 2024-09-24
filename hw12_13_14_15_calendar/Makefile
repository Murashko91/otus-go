BIN := "./bin/calendar"
BIN_SHEDULER := "./bin/scheduler"
BIN_SENDER := "./bin/sender"
DOCKER_IMG="calendar-app:develop"
SCH_DOCKER_IMG="scheduler-app:develop"
SEN_DOCKER_IMG="sender-app:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender
	go build -v -o $(BIN_SHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler

run: build
	$(BIN) -calendar-conf ./configs/calendar_config.yaml

run-scheduler: build
	$(BIN_SHEDULER) -sheduler-conf ./configs/scheduler_config.yaml


run-sender: build
	$(BIN_SENDER) -sender-conf ./configs/sender_config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/calendar/Dockerfile .
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(SCH_DOCKER_IMG) \
		-f build/scheduler/Dockerfile .
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(SEN_DOCKER_IMG) \
		-f build/sender/Dockerfile .

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.60.3

lint: install-lint-deps
	golangci-lint run ./...

generate:
	protoc ./proto/api/EventService.proto --go_out=./proto/gen/go --go-grpc_out=./proto/gen/go
migrate-up:
	goose -dir ./migrations postgres "$(psqlInfo)" up
migrate-down:
	goose -dir ./migrations postgres "$(psqlInfo)" down
up:
	docker compose -f docker-compose.yaml up -d 

up-rebuild: build-img up

down:
	docker compose -f docker-compose.yaml down 

integration-tests:
	go test -v ./test/integration/integration_test.go

helm-up:
	helm install calendar-app ks8
helm-down:
	helm uninstall calendar-app

.PHONY: build run build-img version test lint generate run-scheduler run-sender up up-rebuild down integration-test helm-up helm-down
