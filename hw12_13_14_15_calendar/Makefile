BIN_DIR := "./bin"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

migrate:
	@if [ ! -f "./bin/goose" ]; then \
		echo "installing goose..." && \
		go build -o ./bin/goose github.com/pressly/goose/cmd/goose; \
	fi
	./bin/goose -dir=migrations postgres "$(DSN)" up

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.33.0

lint: install-lint-deps
	golangci-lint run ./...

generate:
	go generate ./...

up:
	docker-compose -f deployments/docker-compose.yaml up -d

down:
	docker-compose -f deployments/docker-compose.yaml down

integration-test: up
	ENV_FILE=.env go test -count=1 ./test

.PHONY: build run build-img run-img version test lint
