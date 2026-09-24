APP := bin/api
CMD := ./cmd/api
MIGRATIONS := migrations
MODULES_DIR := internal/modules
DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/app?sslmode=disable

.PHONY: run build test cover fmt fmt-check vet lint check migrate-up migrate-down docker-up docker-down module clean

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

# Uso: make module MODULE=courses
module:
	@if [ -z "$(MODULE)" ]; then \
		echo "Error: indica el nombre del módulo. Uso: make module MODULE=courses"; \
		exit 1; \
	fi
	@if ! printf '%s' "$(MODULE)" | grep -Eq '^[a-z][a-z0-9_]*$$'; then \
		echo "Error: MODULE debe usar snake_case, comenzar con una letra y contener solo minúsculas, números o _."; \
		exit 1; \
	fi
	@if [ -e "$(MODULES_DIR)/$(MODULE)" ]; then \
		echo "Error: el módulo $(MODULE) ya existe en $(MODULES_DIR)/$(MODULE)."; \
		exit 1; \
	fi
	@set -eu; \
		base="$(MODULES_DIR)/$(MODULE)"; \
		mkdir -p \
			"$$base/domain" \
			"$$base/application" \
			"$$base/infrastructure/persistence/memory" \
			"$$base/infrastructure/persistence/postgres/models" \
			"$$base/infrastructure/persistence/postgres/mappers" \
			"$$base/infrastructure/security" \
			"$$base/infrastructure/mail" \
			"$$base/delivery/http/v1"; \
		printf 'package domain\n' > "$$base/domain/module.go"; \
		printf 'package domain\n' > "$$base/domain/errors.go"; \
		printf 'package application\n' > "$$base/application/ports.go"; \
		printf 'package application\n' > "$$base/application/service.go"; \
		printf 'package application\n' > "$$base/application/service_test.go"; \
		printf 'package v1\n' > "$$base/delivery/http/v1/controller.go"; \
		printf 'package v1\n' > "$$base/delivery/http/v1/dto.go"; \
		touch \
			"$$base/infrastructure/persistence/memory/.gitkeep" \
			"$$base/infrastructure/persistence/postgres/models/.gitkeep" \
			"$$base/infrastructure/persistence/postgres/mappers/.gitkeep" \
			"$$base/infrastructure/security/.gitkeep" \
			"$$base/infrastructure/mail/.gitkeep"; \
		echo "Módulo $(MODULE) creado en $$base"

clean:
	rm -rf bin coverage.out
