APP_MAIN_FILE := ./cmd/bank-api/main.go

DB_URL := $(shell grep '^DB_URL' .env | cut -d '=' -f2 | tr -d '[:space:]')

build:
	@go build -o bin/bank-api $(APP_MAIN_FILE)

migrate-up:
	@goose -dir sql/schema postgres "$(DB_URL)" up 

migrate-down:
	@goose -dir sql/schema postgres "$(DB_URL)" down 

migrate-status:
	@goose -dir sql/schema postgres "$(DB_URL)" status 

migrate-reset:
	@goose -dir sql/schema postgres "$(DB_URL)" reset 

generate:
	@sqlc generate

test: 
	@go test -v ./...

swag-init:
	@swag init -g cmd/bank-api/main.go

swag-fmt:
	@swag init -g cmd/bank-api/main.go 

run: swag-init build  
	@./bin/bank-api
