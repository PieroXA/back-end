.PHONY: run build test migrate-up migrate-down sqlc docker-up docker-down

# ── Development ────────────────────────────────────────────────────────────────

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test ./... -v -count=1

# ── Database ───────────────────────────────────────────────────────────────────

migrate-up:
	migrate -path sql/migrations -database "$$DATABASE_URL" up

migrate-down:
	migrate -path sql/migrations -database "$$DATABASE_URL" down 1

migrate-create:
	migrate create -ext sql -dir sql/migrations -seq $(name)

# ── Code Generation ────────────────────────────────────────────────────────────

sqlc:
	sqlc generate

# ── Docker ─────────────────────────────────────────────────────────────────────

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-reset:
	docker compose down -v && docker compose up -d

# ── Lint ───────────────────────────────────────────────────────────────────────

lint:
	golangci-lint run ./...

# ── Help ───────────────────────────────────────────────────────────────────────

help:
	@echo "Available targets:"
	@echo "  make run            - Run the API server"
	@echo "  make build          - Build the binary"
	@echo "  make test           - Run all tests"
	@echo "  make migrate-up     - Apply all pending migrations"
	@echo "  make migrate-down   - Rollback last migration"
	@echo "  make sqlc           - Generate sqlc code"
	@echo "  make docker-up      - Start Docker containers"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make docker-reset   - Reset Docker (delete volumes)"
