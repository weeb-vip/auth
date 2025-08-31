# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based GraphQL authentication service using gqlgen for GraphQL code generation, GORM for database operations, and JWT for authentication. The service provides user authentication, password reset, and token management functionality.

## Development Commands

### Building and Running
- `go run cmd/cli/main.go server` - Start the development server
- `air` - Hot reload development server (using Air)
- `go build -o auth cmd/cli/main.go` - Build production binary

### Testing and Quality
- `go test ./...` - Run all tests
- `golangci-lint run` - Run linter (configured in .golangci.json)
- `pre-commit run --all-files` - Run pre-commit hooks

### Code Generation
- `make gql` - Generate GraphQL resolvers and models using gqlgen
- `make mocks` - Generate mock implementations for testing

### Database
- `go run cmd/cli/main.go db migrate` - Run database migrations
- `make create-migration name=migration_name` - Create new migration file

## Architecture

### Core Structure
- `cmd/cli/main.go` - CLI entry point with server and database commands
- `server.go` - HTTP server setup with Chi router, CORS, and GraphQL handler
- `graph/` - GraphQL schema files and generated resolvers
- `internal/` - Core business logic organized by domain

### Key Internal Modules
- `internal/jwt/` - JWT token generation and validation
- `internal/keypair/` - Cryptographic key management with rotation
- `internal/services/` - Business logic services (credentials, user client, etc.)
- `internal/db/` - Database connection and migration handling
- `internal/entities/` - Database models and entities
- `internal/resolvers/` - GraphQL resolver implementations

### Configuration
- Uses `github.com/jinzhu/configor` for environment-based config
- Configuration files in `config/` directory (config.dev.json, config.docker.json)
- Supports development and Docker environments

### GraphQL Setup
- Schema files: `graph/*.graphqls`
- Generated code: `graph/generated/`
- Resolvers: `graph/` with follow-schema layout
- Federation enabled for microservices architecture

### Database
- GORM with MySQL driver
- Migrations in `internal/migrations/scripts/`
- Database initialization script: `.db_init.sql`

### Security Features
- Rotating JWT signing keys with configurable duration
- Key publishing to external key management service
- Password hashing with bcrypt
- CORS configuration for frontend integration

### Development Tools
- Air for hot reloading (configured in .air.toml)
- Pre-commit hooks for code quality
- Comprehensive linting rules with golangci-lint
- Mock generation for testing dependencies

### Docker Support
- Multi-stage Dockerfile for production builds
- Docker Compose for local development with database
- Separate local and production compose configurations