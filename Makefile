.EXPORT_ALL_VARIABLES:
COMPOSE_FILE ?= ./build/docker-compose/docker-compose.yml
MIGRATIONS_DIRECTORY ?= ./migrations
SERVICE_NAME ?= gin-box

.PHONY: help
help: # Display available commands
	@echo "\n\033[0;33mAvailable make commands:\033[0m\n"
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

.PHONY: si
si: # Generate Swagger 2.0 docs
	swag init

.PHONY: sg
sg: # Generate Go queries from query.sql file
	sqlc generate

.PHONY: start
start: build up # Build and spin up the application in one command

.PHONY: restart
restart: down start # Restart the application

.PHONY: logs
logs: # Follow logs of all running containers
	docker compose logs --follow

.PHONY: up
up:  # Spin up the application
	docker compose -f $(COMPOSE_FILE) up -d
	docker compose ps

.PHONY: down
down: # Shut down the application
	docker compose down

.PHONY: connect
connect: # Connect to the running core container
	docker exec -it e2d985c5ad27 /bin/sh

.PHONY: create-migration
create-migration: # Create a new migration (SQL) file with specified name. Use 'name=' argument
	goose -dir $(MIGRATIONS_DIRECTORY) create $(name) sql

.PHONY: migrate-by-one
migrate-by-one: # Migrate the DB up by 1
	goose up-by-one

.PHONY: migrate-all
migrate-all: # Migrate the DB to the most recent version available
	goose up

.PHONY: migrate-down
migrate-down: # Roll back the DB version by 1
	goose down

.PHONY: db-status
db-status: # Dump the migration status of the database
	goose status

.PHONY: db-version
db-version: # Print the current version of the database
	goose version

.PHONY: build
build: # Build docker image of the application
	docker build --tag=$(SERVICE_NAME) --file=build/docker/Dockerfile .
