.PHONY: migrate-up migrate-down generate run test

MAIN_FILE := ./cmd/bank-api/main.go

DB_URL := $(shell grep '^DB_URL' .env | cut -d '=' -f2 | tr -d '[:space:]')

build:
	@go build -o bin/bank-api $(MAIN_FILE)

run: build
	@./bin/bank-api

migrate-up:
	@goose -dir sql/schema postgres "$(DB_URL)" up 

migrate-down:
	@goose -dir sql/schema postgres "$(DB_URL)" down 

generate:
	@sqlc generate

test: 
	@go test -v ./...
