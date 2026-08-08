# Mock Generation Guide

## To Generate Mocks

Run the following command in the project root:

```bash
make generate-mocks
```

This will:
1. Check if mockery is installed (install it automatically if missing)
2. Generate mocks for all interfaces in the `server/repositories` directory
3. Place the generated mocks in `server/mocks/` directory

## Manual Mock Generation (Alternative)

If you prefer to run the command manually:

```bash
# Install mockery if not present
go install github.com/vektra/mockery/v2@v2.28.2

# Generate mocks from the server directory
cd server && mockery --dir=./repositories --output=./mocks --all --keep-header --filename=mock_{interface}.go
```

## Running Tests

To run all unit tests:

```bash
make test
```

Or manually:

```bash
go test ./... -v
```

## Troubleshooting

### "mockery: command not found"
- Run `make generate-mocks` which will automatically install mockery
- Or install manually: `go install github.com/vektra/mockery/v2@v2.28.2`

### Import path errors when generating mocks
- Ensure you're running the command from the project root or server directory correctly
- The mockery command should be run from within the server directory:
  ```bash
  cd server && mockery --dir=./repositories --output=./mocks ...
  ```

### Package name errors
- Verify that your `go.mod` file is correctly configured in the server directory
- The module should be defined as `module url_shortener`
- Ensure all imports use the correct path: `url_shortener/models` etc.