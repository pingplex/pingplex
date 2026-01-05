.PHONY: db-migrate db-seed

db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	go run . migrate up

db-seed: ## Seed development data
	@echo "Seeding development data..."
	go run . seed
