APP_NAME=chess-api

.PHONY: dev up down logs test tidy fmt lint sqlc migrate-up migrate-down build run

up:
	docker compose up -d db

down:
	docker compose down

logs:
	docker compose logs -f

tidy:
	go mod tidy

fmt:
	go fmt ./...

test:
	go test ./...

# --- sqlc ---
sqlc:
	sqlc generate

# --- migrations (golang-migrate CLI required) ---
# Install (one-time): go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate-up:
	migrate -database "$(DATABASE_URL)" -path internal/db/migrations up

migrate-down:
	migrate -database "$(DATABASE_URL)" -path internal/db/migrations down 1

build:
	go build -o bin/$(APP_NAME) ./cmd/api

run:
	DATABASE_URL="$(DATABASE_URL)" JWT_SECRET="$(JWT_SECRET)" JWT_TTL_MIN="$(JWT_TTL_MIN)" HTTP_ADDR=":8080" go run ./cmd/api
