# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go web application template built with Gin framework that follows a clean architecture pattern. The project is structured in distinct layers: infrastructure (infra), business logic (logic), and web layer (web).

## Architecture

### Directory Structure
- `cmd/` - Application entry points
- `infra/` - Infrastructure layer (database, redis, config, logging, monitoring)
- `logic/` - Business logic layer (services, models, business rules)
- `web/` - Web layer (HTTP handlers, routing, middleware)
- `config/` - Configuration files

### Key Components
- **Main Entry**: `cmd/run.go` - Application bootstrap and initialization
- **Configuration**: `infra/config/config.go` - YAML-based config with environment variable overrides
- **Database**: GORM with MySQL driver, accessible via `infra.DB`
- **Caching**: Redis client, accessible via `infra.Redis`
- **Logging**: Zap-based structured logging with custom middleware
- **HTTP Framework**: Gin with custom response utilities

### Initialization Flow
1. `infra.Init()` - Load config, setup logging, initialize DB/Redis
2. `logic.Init()` - Initialize service layer with database/redis dependencies
3. `web.SetupRouter()` - Setup HTTP routes and middleware
4. Start monitoring and HTTP server on port 8080

## Development Commands

### Build and Run
```bash
# Build the application
go build -o myapp .

# Run directly
go run .

# Run with specific module
go run cmd/run.go
```

### Docker
```bash
# Build Docker image
docker build -t go-tpl .

# The application expects a `config/` directory with conf.yml
```

### Development Tools
```bash
# Format code
go fmt ./...

# Vet for potential issues
go vet ./...

# Get dependencies
go mod tidy

# Download dependencies
go mod download
```

## Configuration

Configuration is loaded from `config/conf.yml` with environment variable overrides:
- Database: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- Redis: `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`
- JWT: `JWT_SECRET`, `JWT_EXPIRE_TIME`
- Server: `SERVER_PORT`, `SERVER_MODE`

## Key Patterns

### Service Layer
Services are initialized in `logic/init.go` and depend on infra components:
- `UserSvc` - User business logic with DB and Redis
- `RoleSvc` - Role management with DB only

### HTTP Handlers
- Handlers are organized by domain in `web/rest/{domain}/`
- Each domain has `route.go` for route registration and `handler.go` for request handling
- Use `base.OK()`, `base.OKWithData()` for standardized responses

### Logging
- Use `logx.Info(c, message)` for request-scoped logging
- Infrastructure uses standard zap logger

### Database Models
- Models are defined in `logic/{domain}/model.go`
- Use GORM annotations and standard patterns

## Adding New Features

1. **New Domain**: Create `logic/{domain}/` and `web/rest/{domain}/` directories
2. **Service**: Add service in `logic/{domain}/service.go` and initialize in `logic/init.go`
3. **HTTP Layer**: Add `handler.go` and `route.go` in `web/rest/{domain}/`
4. **Register**: Add domain registration in `web/router.go`

## Dependencies

Key external dependencies:
- `github.com/gin-gonic/gin` - HTTP framework
- `gorm.io/gorm` - ORM
- `github.com/redis/go-redis/v9` - Redis client
- `go.uber.org/zap` - Structured logging
- `github.com/goccy/go-yaml` - YAML parsing
- `github.com/prometheus/client_golang` - Metrics