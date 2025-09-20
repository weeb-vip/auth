# Auth Service

Go-based GraphQL authentication service using gqlgen for GraphQL code generation, GORM for database operations, and JWT for authentication. The service provides user authentication, password reset, and token management functionality with HTTP-only cookie support.

## Features

- JWT token authentication with rotating keys
- HTTP-only cookies for secure token storage
- GraphQL API with Apollo Federation support
- Password reset functionality
- User management and credentials handling
- CORS configuration for cross-origin requests
- Docker support for development and production

## Quick Start

### Prerequisites

- Go 1.23+
- Docker and Docker Compose
- MySQL (for local development without Docker)

### Development Commands

```bash
# Build and run locally
go run cmd/cli/main.go server

# Hot reload development (requires Air)
air

# Build production binary
go build -o auth cmd/cli/main.go

# Run tests
go test ./...

# Run linter
golangci-lint run

# Generate GraphQL code
make gql

# Generate mocks for testing
make mocks

# Database migrations
go run cmd/cli/main.go db migrate

# Create new migration
make create-migration name=migration_name
```

## Docker Testing

### Option 1: External Services (docker-compose.yaml)

Use this when you want to run the auth service locally and only use Docker for dependencies:

```bash
# Start MySQL and other dependencies
docker-compose up -d

# Run migrations
go run cmd/cli/main.go db migrate

# Start auth service locally
go run cmd/cli/main.go server
```

The service will be available at `http://localhost:3001/graphql`

### Option 2: Full Docker Environment (docker-compose.local.yaml)

Use this for complete containerized testing that builds the Go service:

```bash
# Build and start all services including auth
docker-compose -f docker-compose.local.yaml up --build -d

# Run migrations in the container
docker-compose -f docker-compose.local.yaml exec auth go run cmd/cli/main.go db migrate

# View logs
docker-compose -f docker-compose.local.yaml logs -f auth
```

The service will be available at `http://localhost:3001/graphql`

## Testing Authentication

### Login Test

Test the login mutation with cookie support:

```bash
curl -X POST http://localhost:3001/graphql \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -c cookies.txt \
  -d '{
    "query": "mutation CreateSession($input: LoginInput) { CreateSession(input: $input) { id Credentials { token refresh_token } } }",
    "variables": {
      "input": {
        "username": "admin@floretos.com",
        "password": "@pfelor@nge1!"
      }
    }
  }'
```

### Verify Cookies

Check that HTTP-only cookies were set:

```bash
# View saved cookies
cat cookies.txt

# Test authenticated request using cookies
curl -X POST http://localhost:3001/graphql \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -b cookies.txt \
  -d '{
    "query": "query { me { id username } }"
  }'
```

### Refresh Token Test

Test token refresh functionality:

```bash
curl -X POST http://localhost:3001/graphql \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -b cookies.txt \
  -c cookies.txt \
  -d '{
    "query": "mutation RefreshToken($refreshToken: String!) { RefreshToken(refreshToken: $refreshToken) { id Credentials { token refresh_token } } }",
    "variables": {
      "refreshToken": "your_refresh_token_here"
    }
  }'
```

## Federation Testing

When testing with Apollo Router at `https://gateway.weeb.vip`:

```bash
# Test login through gateway
curl -X POST https://gateway.weeb.vip/graphql \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -H "apollographql-client-name: test-client" \
  -H "apollographql-client-version: 1.0.0" \
  -H "x-remote-ip: 127.0.0.1" \
  -H "x-user-agent: curl/test" \
  -H "origin: https://weeb.vip" \
  -c cookies.txt \
  -d '{
    "query": "mutation CreateSession($input: LoginInput) { CreateSession(input: $input) { id Credentials { token refresh_token } } }",
    "variables": {
      "input": {
        "username": "admin@floretos.com",
        "password": "@pfelor@nge1!"
      }
    }
  }'
```

## Architecture

### Core Structure

- `cmd/cli/main.go` - CLI entry point with server and database commands
- `server.go` - HTTP server setup with Chi router, CORS, and GraphQL handler
- `graph/` - GraphQL schema files and generated resolvers
- `internal/` - Core business logic organized by domain

### Key Components

- **JWT Management** (`internal/jwt/`) - Token generation and validation with rotating keys
- **Services** (`internal/services/`) - Business logic for credentials, users, etc.
- **Resolvers** (`internal/resolvers/`) - GraphQL resolver implementations with cookie support
- **Database** (`internal/db/`) - GORM-based database operations and migrations
- **Configuration** (`config/`) - Environment-based configuration management

### Cookie Implementation

The service sets HTTP-only cookies for both access and refresh tokens:

- **Access Token Cookie**: `access_token`, 1-hour expiration, HttpOnly, SameSite=None
- **Refresh Token Cookie**: `refresh_token`, 7-day expiration, HttpOnly, SameSite=None
- **Domain Configuration**: Configurable via `CookieDomain` in config files
- **Backwards Compatibility**: Tokens still returned in GraphQL response body

### Security Features

- Rotating JWT signing keys with configurable duration
- HTTP-only cookies prevent XSS token theft
- CORS configuration for cross-origin cookie support
- bcrypt password hashing
- Configurable cookie domains for federation

## Configuration

### Environment Files

- `config/config.dev.json` - Development configuration
- `config/config.docker.json` - Docker environment configuration

### Key Settings

```json
{
  "app": {
    "port": 3001,
    "environment": "development",
    "cookie_domain": "localhost"
  },
  "db": {
    "host": "localhost",
    "port": 3306,
    "user": "root",
    "password": "password",
    "dbname": "auth"
  }
}
```

## Troubleshooting

### Database Connection Issues

```bash
# Check MySQL is running
docker-compose ps

# View database logs
docker-compose logs mysql

# Connect to database directly
docker-compose exec mysql mysql -u root -p auth
```

### Service Not Starting

```bash
# Check service logs
docker-compose -f docker-compose.local.yaml logs auth

# Verify migrations ran
docker-compose -f docker-compose.local.yaml exec auth go run cmd/cli/main.go db status
```

### Cookie Issues

- Ensure `CookieDomain` is correctly configured for your environment
- For localhost testing, use `"localhost"` as cookie domain
- For production federation, use `".weeb.vip"` (note: Go HTTP lib strips leading dot)
- Check CORS configuration allows credentials for cross-origin requests

### Apollo Router Issues

- Verify router configuration accepts new CORS policy format
- Check that federation schema is properly published
- Ensure service discovery is working correctly

## Development Workflow

1. **Start Dependencies**: `docker-compose up -d`
2. **Run Migrations**: `go run cmd/cli/main.go db migrate`
3. **Start Service**: `air` (for hot reload) or `go run cmd/cli/main.go server`
4. **Run Tests**: `go test ./...`
5. **Lint Code**: `golangci-lint run`
6. **Generate Code**: `make gql` and `make mocks` as needed

## Production Deployment

The service is designed for containerized deployment with:

- Multi-stage Docker builds for optimized images
- Health checks and graceful shutdown
- Environment-based configuration
- Database migration support
- Federation-ready GraphQL schema

For production deployment, ensure:

- Proper cookie domain configuration (`.weeb.vip`)
- HTTPS-only cookies (`Secure: true`)
- Production database credentials
- Monitoring and logging setup