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
	@sqlc generate -f ./services/storage/internal/database/sqlc/sqlc.yaml
	@echo "SQLC generation completed."


# Generate Go code from Protocol Buffers
grpc-gen:
	@echo "Generating Go code from Protocol Buffers..."
	@protoc --go_out=paths=source_relative:./common/proto_go --go-grpc_out=paths=source_relative:./common/proto_go --proto_path=./common/proto ./common/proto/*.proto
	@echo "Go code generation completed."

# Clean generated Go code
grpc-clean:
	@echo "Cleaning generated Go code..."
	@rm -f ./common/proto_go/*.go
	@echo "Clean completed."