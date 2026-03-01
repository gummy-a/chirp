SHELL := /bin/bash

MEDIA_SERVICE_DATABASE_URL ?= postgres://postgres:password@localhost:5432/media_service?sslmode=disable

MIGRATE := migrate
SQLC := sqlc
PROTOC := protoc

.PHONY: help start-postgres stop-postgres migrate schema sqlc-generate openapi-generate tidy all

help:
	@printf "Usage:\n"
	@printf "  make start-postgres   Start local postgres container\n"
	@printf "  make stop-postgres    Stop and remove postgres container\n"
	@printf "  make start-redis      Start local redis container\n"
	@printf "  make stop-redis       Stop and remove redis container\n"
	@printf "  make migrate          Run DB migrations (uses MEDIA_SERVICE_DATABASE_URL)\n"
	@printf "  make schema           Dump DB schema to sql/schema.sql\n"
	@printf "  make sqlc-generate    Run sqlc generate in sql\n"
# 	@printf "  make protoc-generate  Generate gRPC code from backend/proto/media/media.proto\n"
	@printf "  make openapi-generate Run OpenAPI generation (go run api/generate.go from repo root)\n"
	@printf "  make tidy             Run go mod tidy in backend/servicies/media\n"
	@printf "  make all              Run migrate, schema, sqlc-generate, openapi-generate, tidy\n"

start-postgres:
	@docker run --name postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=media_service -p 5432:5432 -d postgres
	@sleep 5 # wait for postgres to start

stop-postgres:
	-@docker stop postgres
	-@docker rm postgres

start-redis:
	@docker run --name redis -d -p 6379:6379 redis:latest

stop-redis:
	-@docker stop redis
	-@docker rm redis

migrate:
	@echo "Running migrations..."
	@$(MIGRATE) -path sql/migrations/ -database "$(MEDIA_SERVICE_DATABASE_URL)" up

schema:
	@echo "Dumping schema to sql/schema.sql..."
	@docker run --rm --network host postgres:18.1 pg_dump --schema-only --no-owner --no-privileges --no-comments "$(MEDIA_SERVICE_DATABASE_URL)" > sql/schema.sql

sqlc-generate:
	@echo "Generating sqlc code..."
	@cd sql && $(SQLC) generate

# protoc-generate:
# 	@echo "Generating protobuf (gRPC) code..."
# 	@cd ../../proto/media && $(PROTOC) --go_out=../../servicies/media/internal/infrastructure/grpc/ --go-grpc_out=../../servicies/media/internal/infrastructure/grpc/ media.proto

openapi-generate:
	@echo "Generating OpenAPI code (run from repo root)..."
	@cd ../../../ && go run api/generate.go

tidy:
	@echo "Running go mod tidy in backend/servicies/media..."
	@go mod tidy

all: stop-postgres start-postgres stop-redis start-redis migrate schema sqlc-generate openapi-generate tidy
