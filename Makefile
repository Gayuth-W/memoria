.PHONY: build run down test migrate-up migrate-down

build:
	docker-compose build

run:
	docker-compose up -d

down:
	docker-compose down

test:
	go test ./...

migrate-up:
	goose -dir migrations/schema postgres "postgres://postgres:postgres@localhost:5432/memoria?sslmode=disable" up

migrate-down:
	goose -dir migrations/schema postgres "postgres://postgres:postgres@localhost:5432/memoria?sslmode=disable" down
