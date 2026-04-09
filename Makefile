APP_NAME=lunchvote
CMD_DIR=./cmd/api

.PHONY: tidy test run build up down

tidy:
	go mod tidy

test:
	go test ./...

run:
	go run $(CMD_DIR)

build:
	go build -o bin/$(APP_NAME) $(CMD_DIR)

up:
	docker compose up -d

down:
	docker compose down