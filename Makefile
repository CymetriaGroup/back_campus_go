APP=bin/api
MIGRATIONS=migrations
DATABASE_URL?=postgres://postgres:postgres@localhost:5432/app?sslmode=disable

.PHONY: run build test cover fmt vet lint migrate-up migrate-down docker-up docker-down
run:
	go run ./cmd/api
build:
	go build -o $(APP) ./cmd/api
test:
	go test ./... -race -v
cover:
	go test ./... -coverprofile=coverage.out
fmt:
	gofmt -w .
vet:
	go vet ./...
lint:
	golangci-lint run
migrate-up:
	migrate -path $(MIGRATIONS) -database "$(DATABASE_URL)" up
migrate-down:
	migrate -path $(MIGRATIONS) -database "$(DATABASE_URL)" down 1
docker-up:
	docker compose up -d --build
docker-down:
	docker compose down
