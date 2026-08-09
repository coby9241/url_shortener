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
	@export GOBIN=$${GOBIN:-$$HOME/go/bin}; \
	echo "GOBIN after export: $$GOBIN"; \
	export PATH=$$GOBIN:$$PATH; \
	if ! which mockery >/dev/null 2>&1; then \
	    echo "mockery not found. Installing..."; \
	    go install github.com/vektra/mockery/v2@v2.53.6; \
	fi; \
	cd server && $$GOBIN/mockery
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