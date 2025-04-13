include .env
# Variables
DB_DRIVER=postgres
DIR=internal/infrastructure/storage/postgres/migrations
MIGRATE=migrate

# Default target
.DEFAULT_GOAL := help

# Colors
RESET=\033[0m
GREEN=\033[32m
YELLOW=\033[33m
BLUE=\033[34m
CYAN=\033[36m
BOLD=\033[1m

# Help
.PHONY: help
help: ## Show available commands
	@printf "${CYAN}${BOLD}Available Commands:${RESET}\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "${YELLOW}%-20s${RESET} %s\n", $$1, $$2}'

# Create a new migration
.PHONY: create-migration
create-migration: ## Create a new migration file
	@printf "${BLUE}Creating a new migration file...${RESET}\n"
	@read -p "Enter migration name: " name; \
	$(MIGRATE) create -ext sql -dir $(DIR) -seq "$$name" && \
	printf "${GREEN}Migration file created successfully!${RESET}\n"

# Migrate up
.PHONY: migrate-up
migrate-up: ## Apply all up migrations
	@printf "${BLUE}Applying all up migrations...${RESET}\n"
	$(MIGRATE) -path $(DIR) -database "$(DB_URL)" up && \
	printf "${GREEN}All migrations applied successfully!${RESET}\n"

# Migrate down
.PHONY: migrate-down
migrate-down: ## Rollback the last migration
	@printf "${BLUE}Rolling back the last migration...${RESET}\n"
	$(MIGRATE) -path $(DIR) -database "$(DB_URL)" down 1 && \
	printf "${GREEN}Last migration rolled back successfully!${RESET}\n"

# Migrate force
.PHONY: migrate-force
migrate-force: ## Force the migration version
	@printf "${BLUE}Forcing the migration version...${RESET}\n"
	@read -p "Enter version to force: " version; \
	$(MIGRATE) -path $(DIR) -database "$(DB_URL)" force "$$version" && \
	printf "${GREEN}Migration version forced to %s successfully!${RESET}\n" "$$version"

# Clean
.PHONY: clean
clean: ## Clean all compiled binaries
	@printf "${BLUE}Cleaning up...${RESET}\n"
	rm -rf ./bin && \
	printf "${GREEN}Cleanup complete!${RESET}\n"

# Generate swagger docs
.PHONY: swagger
swagger: ## Generate swagger docs
	@printf "${BLUE}Generating swagger docs...${RESET}\n"
	@swag init -g app/main.go -d cmd,internal && swag fmt && \
	printf "${GREEN}Swagger docs generated successfully!${RESET}\n"