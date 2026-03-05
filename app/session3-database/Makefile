# Makefile for Session 3 - Database Integration
# Common commands for development

.PHONY: help db-start db-stop db-logs db-shell migrate-up migrate-down run clean test

# Default target
help:
	@echo "Available commands:"
	@echo "  make db-start     - Start PostgreSQL with Docker"
	@echo "  make db-stop      - Stop PostgreSQL"
	@echo "  make db-logs      - Show database logs"
	@echo "  make db-shell     - Open PostgreSQL shell"
	@echo "  make migrate-up   - Run database migrations (up)"
	@echo "  make migrate-down - Rollback database migrations (down)"
	@echo "  make run          - Start the server"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean up containers and volumes"
	@echo "  make setup        - Complete setup (db-start + migrate-up)"

# Database operations
db-start:
	@echo "🚀 Starting PostgreSQL..."
	docker-compose up -d
	@echo "✅ PostgreSQL started on localhost:5432"

db-stop:
	@echo "🛑 Stopping PostgreSQL..."
	docker-compose stop
	@echo "✅ PostgreSQL stopped"

db-logs:
	@echo "📋 Database logs (Ctrl+C to exit):"
	docker-compose logs -f db

db-shell:
	@echo "🐘 Opening PostgreSQL shell..."
	docker-compose exec db psql -U postgres -d mini_asm

# Migrations
migrate-up:
	@echo "⬆️  Running migrations (up)..."
	docker-compose exec db psql -U postgres -d mini_asm -f /docker-entrypoint-initdb.d/001_create_assets.up.sql
	@echo "✅ Migrations applied"

migrate-down:
	@echo "⬇️  Rolling back migrations (down)..."
	docker-compose exec db psql -U postgres -d mini_asm -f /docker-entrypoint-initdb.d/001_create_assets.down.sql
	@echo "✅ Migrations rolled back"

# Application
run:
	@echo "🚀 Starting Mini ASM Server..."
	go run cmd/server/main.go

# Testing
test:
	@echo "🧪 Running tests..."
	go test ./... -v

# Cleanup
clean:
	@echo "🧹 Cleaning up..."
	docker-compose down -v
	@echo "✅ Containers and volumes removed"

# Complete setup
setup: db-start
	@echo "⏳ Waiting for database to be ready..."
	@sleep 5
	@echo "✅ Complete! Run 'make run' to start the server"

# Development helpers
dev: db-start run

# Show database tables
db-tables:
	@echo "📊 Tables in mini_asm database:"
	docker-compose exec db psql -U postgres -d mini_asm -c "\dt"

# Show all assets
db-assets:
	@echo "📦 Assets in database:"
	docker-compose exec db psql -U postgres -d mini_asm -c "SELECT * FROM assets;"

# Reset database (down + up)
db-reset: migrate-down migrate-up
	@echo "♻️  Database reset complete"

# 🎓 TEACHING NOTES:
#
# Makefile benefits:
# - Simple commands: make db-start instead of docker-compose up -d
# - Self-documenting: make help
# - Consistent workflow for all team members
# - Great for CI/CD pipelines
#
# Common workflow:
# 1. make setup      (first time only)
# 2. make run        (start server)
# 3. make db-logs    (debug database)
# 4. make db-shell   (run SQL queries)
# 5. make clean      (cleanup when done)
