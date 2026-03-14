.PHONY: run build migrate-up migrate-down db-create db-drop tidy

# Application
run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

# Database
db-create:
	createdb discipline

db-drop:
	dropdb discipline

migrate-up:
	psql -d discipline -f migrations/001_initial.up.sql

migrate-down:
	psql -d discipline -f migrations/001_initial.down.sql

migrate-reset: migrate-down migrate-up

# Dependencies
tidy:
	go mod tidy
