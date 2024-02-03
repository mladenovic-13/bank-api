.PHONY: build swag-init swag-fmt run

APP_MAIN_FILE := ./cmd/main.go

TEST_DB_URL := $(shell grep '^TEST_DB_URL' .env | cut -d '=' -f2 | tr -d '[:space:]')

build:
	@go build -o bin/bank-api $(APP_MAIN_FILE)

test: 
	@go test -cover -v ./...

swag-init:
	@swag init -g cmd/bank-api/main.go

swag-fmt:
	@swag init -g cmd/bank-api/main.go 

run: build  
	@./bin/bank-api
