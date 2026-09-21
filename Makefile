APP := bin/api
CMD := ./cmd/api
MIGRATIONS := migrations
DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/app?sslmode=disable

.PHONY: run build test cover fmt fmt-check vet lint check migrate-up migrate-down docker-up docker-down clean

run:
	go run $(CMD)

build:
	go build -o $(APP) $(CMD)

test:
	go test ./... -race -v

cover:
	go test ./... -race -coverprofile=coverage.out

fmt:
	gofmt -w .

fmt-check:
	test -z "$$(gofmt -l .)"

vet:
	go vet ./...

lint:
	golangci-lint run

check: fmt-check vet test build

migrate-up:
	migrate -path $(MIGRATIONS) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS) -database "$(DATABASE_URL)" down 1

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

clean:
	rm -rf bin coverage.out
