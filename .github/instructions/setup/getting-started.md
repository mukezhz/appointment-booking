# Getting Started 

This guide will help you set up and run the project locally.

## Prerequisites

- Go 1.20 or higher
- Docker and Docker Compose
- Python 3.x (for development tools)
- [Atlas CLI](https://atlasgo.io) for migrations

## Quick Start

1. Clone and set up environment:
   ```bash
   # Clone repository
   git clone <repository-url>
   cd appointment-booking

   # Copy and configure environment
   cp .env.example .env
   # Edit .env with your settings
   ```

2. Set up development tools:
   ```bash
   # Install Git hooks and linters
   make lint-setup
   ```

3. Start the application:
   ```bash
   # Using Docker (recommended)
   docker-compose up -d

   # Or locally
   go run main.go app:serve
   ```

4. View available commands:
   ```bash
   go run main.go -help
   ```

## Development Workflow

1. **Database Migrations**
   - View status: `make migrate-status`
   - Create migration: `make migrate-diff`
   - Apply migrations: `make migrate-apply`
   - Rollback: `make migrate-down`

2. **Testing**
   - Run all tests: `go test ./... -v`
   - Run with Ginkgo: `ginkgo -v --cover -r ./domain/... ./pkg/...`
   - Generate coverage: `make test-coverage`

3. **Code Quality**
   - Run linter: `make lint`
   - Format code: `make fmt`

## Additional Resources

- [Architecture Overview](../architecture/01-project-overview.md)
- [Development Guide](../architecture/03-feature-development.md)
- [Testing Strategy](../testing/01-strategy.md)
- [Error Handling Guide](../error-guide.md)
