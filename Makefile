up:
	@echo "Running in Development mode..."
	@docker compose -f ./build/compose/compose.yaml up -d
	@echo "Development mode completed."

down:
	@echo "Running in Development mode..."
	@docker compose -f ./build/compose/compose.yaml down --volumes
	@echo "Development mode completed."


# Make migrations
migrate-up:
	@echo "Making migrations..."
	@migrate -path ./services/user/internal/database/migrations -database "postgres://root:password@localhost:5432/user-db?sslmode=disable" -verbose up
	@echo "Migrations completed."

# Delete migrations
migrate-down:
	@echo "Deleting migrations..."
	@migrate -path ./services/user/internal/database/migrations -database "postgres://root:password@localhost:5432/user-db?sslmode=disable" -verbose down

# SQLC generate
sqlc-gen:
	@echo "Generating SQLC..."
	@sqlc generate -f ./services/user/internal/database/sqlc/sqlc.yaml
	@echo "SQLC generation completed."