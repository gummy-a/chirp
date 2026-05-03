SHELL := /bin/bash

DATABASE_URL ?= postgres://postgres:password@localhost:5432/chirp?sslmode=disable

MIGRATE := migrate
SQLC := sqlc

.PHONY: help start-postgres stop-postgres migrate schema sqlc-generate all

help:
	@printf "Usage:\n"
	@printf "  make start-postgres   Start local postgres container\n"
	@printf "  make stop-postgres    Stop and remove postgres container\n"
	@printf "  make start-redis      Start local redis container\n"
	@printf "  make stop-redis       Stop and remove redis container\n"
	@printf "  make migrate          Run DB migrations (uses DATABASE_URL)\n"
	@printf "  make schema           Dump DB schema to db/schema.sql\n"
	@printf "  make sqlc-generate    Run sqlc generate in sql\n"
	@printf "  make all              Run migrate, schema, sqlc-generate, openapi-generate, tidy\n"

start-postgres:
	@docker run --name postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=chirp -p 5432:5432 -d postgres
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
	@$(MIGRATE) -path db/migrations/ -database "$(DATABASE_URL)" up

schema:
	@echo "Dumping schema to db/schema.sql..."
	@docker run --rm --network host postgres:18.1 pg_dump --schema-only --no-owner --no-privileges --no-comments "$(DATABASE_URL)" > db/schema.sql

sqlc-generate:
	@echo "Generating sqlc code..."
	@cd db && $(SQLC) generate
	@python3 tools/migrate.py

all: stop-postgres start-postgres stop-redis start-redis migrate schema sqlc-generate
