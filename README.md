# URL Shortener

A modern URL shortener application built with a Go backend and React frontend, using PostgreSQL for data storage.

## Table of Contents
- [Deployment](#deployment)
- [Requirements](#requirements)
- [Getting Started](#getting-started)
  - [Using Docker Compose (Recommended)](#using-docker-compose-recommended)
  - [Local Development (Without Docker)](#local-development-without-docker)
    - [Backend Setup](#backend-setup)
    - [Frontend Setup](#frontend-setup)
- [Mock Generation](#mock-generation)
- [Running Tests](#running-tests)
- [Screenshots](#screenshots)
- [Further Help](#further-help)

## Deployment

The application is deployed on Render with the following URL:
- Web UI: https://url-shortener-ui-wibi.onrender.com

The deployment consists of:
- Static frontend hosted on Render's static site service
- Backend API running in a Docker container on Render
- PostgreSQL database provided by Render

## Requirements
- Docker and Docker Compose
- (Optional for local development) Go 1.25+ and Node.js 18+

## Getting Started

### Using Docker Compose (Recommended)
This is the easiest way to get the application running with all dependencies.

1. Clone the repository:
   ```bash
   git clone https://github.com/coby9241/url_shortener
   cd url_shortener
   ```

2. Create a `.env` file in the `server` directory (if not already present) with the following variables:
   ```dotenv
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=postgres
   POSTGRES_DB=url_shortener
   DB_HOST=postgres
   DB_PORT=5432
   DB_SSLMODE=disable
   ```
   Note: The docker-compose file already sets up a PostgreSQL service, so the backend will connect to it automatically.

3. Start the application:
   ```bash
   docker-compose up --build
   ```

4. Once the containers are healthy, access the application:
   - Frontend: http://localhost:3001
   - Backend API: http://localhost:8080


Open up the page at http://localhost:3001. Then key in the URL you wish to generate the redirect for and press "Shorten URL".

![Original URL](docs/images/url.png)

You will get a shortened URL generated after a short while. Copy the URL and paste it into the URL bar. Then you will be redirected to the original URL.

![Shortened URL](docs/images/shortened_url.png)

### Local Development (Without Docker)
If you prefer to run the services directly on your machine for development:

#### Backend Setup
1. Navigate to the server directory:
   ```bash
   cd server
   ```

2. Install Go dependencies:
   ```bash
   go mod download
   ```

3. Create a `.env` file in the server directory with the following (adjust as needed):
   ```dotenv
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=postgres
   POSTGRES_DB=url_shortener
   DB_HOST=localhost
   DB_PORT=5432
   DB_SSLMODE=disable
   ```

4. Ensure you have a PostgreSQL instance running on localhost:5432 (or update the host/port in .env).

5. Run the backend:
   ```bash
   go run main.go
   ```
   The backend will be available at http://localhost:8080

#### Frontend Setup
1. Navigate to the client directory:
   ```bash
   cd client
   ```

2. Install Node.js dependencies:
   ```bash
   npm ci
   ```

3. Start the frontend development server:
   ```bash
   npm start
   ```
   The frontend will be available at http://localhost:3001 (proxy set to backend via package.json setup in Docker, but for local dev you may need to adjust)

   Note: For local development without Docker, you may need to set the `REACT_APP_API_URL` environment variable in the client to point to your backend (e.g., `http://localhost:8080`). You can do this by creating a `.env` file in the client directory:
   ```
   REACT_APP_API_URL=http://localhost:8080
   ```

## Mock Generation
This project uses [mockery](https://github.com/vektra/mockery) to generate mocks for interfaces in the `server/repositories` directory.

### To Generate Mocks
Run the following command in the project root:
```bash
make generate-mocks
```

This will:
1. Check if mockery is installed (install it automatically if missing)
2. Generate mocks for all interfaces in the `server/repositories` directory
3. Place the generated mocks in `server/mocks/` directory

### Manual Mock Generation (Alternative)
If you prefer to run the command manually:

```bash
# Install mockery if not present
go install github.com/vektra/mockery/v2@v2.28.2

# Generate mocks from the server directory
cd server && mockery --dir=./repositories --output=./mocks --all --keep-header --filename=mock_{interface}.go
```

### Running Tests
To run all unit tests:
```bash
make test
```

Or manually:
```bash
go test ./... -v
```

### Troubleshooting
- **"mockery: command not found"**
  - Run `make generate-mocks` which will automatically install mockery
  - Or install manually: `go install github.com/vektra/mockery/v2@v2.28.2`

- **Import path errors when generating mocks**
  - Ensure you're running the command from the project root or server directory correctly
  - The mockery command should be run from within the server directory:
    ```bash
    cd server && mockery --dir=./repositories --output=./mocks ...
    ```

- **Package name errors**
  - Verify that your `go.mod` file is correctly configured in the server directory
  - The module should be defined as `module url_shortener`
  - Ensure all imports use the correct path: `url_shortener/models` etc.

