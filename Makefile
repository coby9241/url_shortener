.PHONY: help generate-mocks test build run

# Help target
help:
	@echo "Available targets:"
	@echo "  generate-mocks   Generate mocks using mockery"
	@echo "  test             Run unit tests"
	@echo "  build            Build the application"
	@echo "  run              Run the application"
	@echo "  help             Show this help message"

# Generate mocks using mockery
generate-mocks:
	@echo "Generating mocks..."
	@which mockery >/dev/null 2>&1 || { \
		echo "mockery not found. Installing..."; \
		go install github.com/vektra/mockery/v2@v2.28.2; \
	}
	cd server && mockery --dir=./repositories --output=./mocks --all --keep-header --filename=mock_{interface}.go
	@echo "Mocks generated successfully in ./server/mocks/"

# Run tests
test:
	cd server && go test ./... -v


# Build the application
build:
	cd server && go build -o url_shortener ./server/main.go

# Run the application
run:
	cd server && go run ./server/main.go