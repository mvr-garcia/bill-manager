# Makefile for common development tasks

SHELL := /bin/bash

# If you want to override the JWT token (for curl targets), set it like:
#   make JWT_TOKEN=eyJ... create-bill

.DEFAULT_GOAL := help

.PHONY: help db-up db-down env run health create-bill list-bills approve-bill audits

help:
	@echo "Usage: make <target> [VARIABLE=value]"
	@echo ""
	@echo "Available targets:"
	@echo "  db-up            Start PostgreSQL (docker-compose up -d)"
	@echo "  db-down          Stop PostgreSQL (docker-compose down)"
	@echo "  env              Copy .env.example to .env (if missing)"
	@echo "  run              Run the API server (go run cmd/api/main.go)"
	@echo "  test             Run unit/integration tests (starts docker-compose services)"
	@echo "  health           Call the health check endpoint"
	@echo "  create-bill      Create a bill via curl (requires JWT_TOKEN)"
	@echo "  list-bills       List bills via curl (requires JWT_TOKEN)"
	@echo "  approve-bill     Approve a bill via curl (requires JWT_TOKEN and BILL_ID)"
	@echo "  audits           View audit logs for a bill via curl (requires JWT_TOKEN and BILL_ID)"
	@echo ""
	@echo "Examples:" 
	@echo "  make db-up"
	@echo "  make env"
	@echo "  make run"
	@echo "  make JWT_TOKEN=... create-bill"

# Start/stop database

db-up:
	docker-compose up -d

db-down:
	docker-compose down

# Environment

env:
	@if [ ! -f .env ]; then \
	  cp .env.example .env; \
	  echo "Created .env from .env.example"; \
	else \
	  echo ".env already exists"; \
	fi

# Run server

run:
	@echo "Running API server on http://localhost:8080 (or $${SERVER_PORT:-8080})"
	@go run cmd/api/main.go

# Run tests (starts docker-compose services first)

test: db-up
	@echo "Starting docker-compose services for tests..."
	@set -e; \
	ret=0; \
	go test ./... || ret=$$?; \
	@echo "Stopping docker-compose services..."; \
	$(MAKE) db-down; \
	exit $$ret

# Health check

health:
	@curl --fail --show-error http://localhost:8080/health

# Helpers for curl requests

JWT_TOKEN ?= 
BILL_ID ?= 

create-bill:
	@if [ -z "$(JWT_TOKEN)" ]; then \
	  echo "Missing JWT_TOKEN (e.g., make JWT_TOKEN=... create-bill)"; exit 1; \
	fi
	@curl -sS -X POST http://localhost:8080/api/v1/bills \
	  -H "Authorization: Bearer $(JWT_TOKEN)" \
	  -H "Content-Type: application/json" \
	  -d '{"description":"Office Supplies","amount":1500.50,"due_date":"2026-04-15T00:00:00Z"}'

list-bills:
	@if [ -z "$(JWT_TOKEN)" ]; then \
	  echo "Missing JWT_TOKEN (e.g., make JWT_TOKEN=... list-bills)"; exit 1; \
	fi
	@curl -sS -X GET http://localhost:8080/api/v1/bills \
	  -H "Authorization: Bearer $(JWT_TOKEN)"

approve-bill:
	@if [ -z "$(JWT_TOKEN)" ]; then \
	  echo "Missing JWT_TOKEN (e.g., make JWT_TOKEN=... approve-bill BILL_ID=...)"; exit 1; \
	fi
	@if [ -z "$(BILL_ID)" ]; then \
	  echo "Missing BILL_ID (e.g., make BILL_ID=... approve-bill JWT_TOKEN=...)"; exit 1; \
	fi
	@curl -sS -X POST http://localhost:8080/api/v1/bills/$(BILL_ID)/approve \
	  -H "Authorization: Bearer $(JWT_TOKEN)"

audits:
	@if [ -z "$(JWT_TOKEN)" ]; then \
	  echo "Missing JWT_TOKEN (e.g., make JWT_TOKEN=... audits BILL_ID=...)"; exit 1; \
	fi
	@if [ -z "$(BILL_ID)" ]; then \
	  echo "Missing BILL_ID (e.g., make BILL_ID=... audits JWT_TOKEN=...)"; exit 1; \
	fi
	@curl -sS -X GET http://localhost:8080/api/v1/bills/$(BILL_ID)/audits \
	  -H "Authorization: Bearer $(JWT_TOKEN)"
