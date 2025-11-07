.PHONY: help build run test clean docker-build docker-up docker-down

help:
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@echo '  help           - Show this help'
	@echo '  build          - Build application'
	@echo '  run            - Run application'
	@echo '  test           - Run tests'
	@echo '  clean          - Clean artifacts'
	@echo '  docker-build   - Build Docker image'
	@echo '  docker-up      - Start services'
	@echo '  docker-down    - Stop services'

build:
	go build -o bin/api ./cmd/api

run:
	go run ./cmd/api

test:
	go test -v ./...

clean:
	rm -rf bin/

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
