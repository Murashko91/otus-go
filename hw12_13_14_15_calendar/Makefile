BIN := "./bin/calendar"
BIN_SHEDULER := "./bin/scheduler"
BIN_SENDER := "./bin/sender"
DOCKER_IMG="calendar:develop"

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
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

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

.PHONY: build run build-img run-img version test lint generate run-scheduler run-sender
