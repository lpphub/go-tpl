# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a production-grade Go web application template built with Gin framework that follows clean architecture principles and domain-driven design. The project features a complete RBAC (Role-Based Access Control) system with user, role, and permission management, following a three-layer architecture: infrastructure (infra), business logic (logic), and web layer (web).

## Tech Stack

- **Go 1.25** - Programming language
- **Gin** - HTTP web framework for REST APIs
- **GORM** - ORM with MySQL driver
- **RDB** - Caching and session management
- **JWT** - Authentication tokens (golang-jwt/jwt/v5)
- **Wire** - Compile-time dependency injection
- **Zap** - High-performance structured logging
- **Prometheus** - Metrics collection
- **goccy/go-yaml** - YAML configuration parsing

## Architecture

### Directory Structure
```
cmd/                    # Application entry points
├── run.go             # Main application bootstrap

config/                # Configuration files
├── conf.yml           # Main YAML configuration

infra/                 # Infrastructure layer
├── config/            # Configuration management
├── dbs/               # Database and RDB setup
├── jwt/               # JWT token handling
├── logging/           # Logging infrastructure
│   └── logx/          # Custom logging utilities
└── monitor/           # Monitoring and metrics

logic/                 # Business logic layer
├── user/              # User domain (model, service)
├── role/              # Role domain
├── permission/        # Permission domain
├── shared/            # Shared utilities (consts, errors, pagination)
├── init.go            # Logic layer initialization
├── wire.go            # Wire dependency injection setup
└── wire_gen.go        # Generated wire code

web/                   # Web layer
├── base/              # Base utilities (response rendering)
├── middleware/        # HTTP middleware
├── rest/              # REST API handlers
│   ├── user/          # User API endpoints
│   ├── role/          # Role API endpoints
│   └── permission/    # Permission API endpoints
├── types/             # Request/response types
└── router.go          # Route configuration

scripts/               # Utility scripts
```

### Key Components
- **Main Entry**: `cmd/run.go` - Application bootstrap following 4-step initialization
- **Configuration**: `infra/config/config.go` - YAML-based config with environment variable overrides
- **Database**: GORM with MySQL driver, accessible via `infra.DB` with transaction support
- **Caching**: RDB client, accessible via `infra.RDB`
- **Logging**: Zap-based structured logging with custom logx utilities and middleware
- **Dependency Injection**: Wire-based compile-time DI for clean dependencies
- **Monitoring**: Prometheus metrics and pprof profiling support
- **Authentication**: JWT-based auth with bcrypt password hashing

### Initialization Flow
1. `infra.Init()` - Load config, setup logging, initialize DB/RDB/JWT
2. `logic.Init()` - Initialize service layer with Wire dependency injection
3. `web.SetupRouter()` - Setup HTTP routes, middleware, and handlers
4. `monitor.SetupMetrics(app)` - Start Prometheus metrics collection
5. Start HTTP server on port 8080

## Domain Models

### User Management
- Complete CRUD operations with status management (active/inactive)
- Password encryption using bcrypt
- Role assignment and RBAC integration
- JWT-based authentication

### Role System
- Hierarchical role management
- Permission assignment capabilities
- Status management

### Permission System
- Granular permission structure with module organization
- Resource-action based permissions
- Integration with role-based access control

## Development Commands

### Build and Run
```bash
# Build the application
go build -o myapp .

# Run directly
go run .

# Run with specific module
go run cmd/run.go

# Generate Wire dependencies
go generate ./logic
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./infra/dbs/...
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

# Run Wire to generate dependencies
wire ./logic/
```

### Docker
```bash
# Build Docker image
docker build -t go-tpl .

# Run with Docker
docker run -p 8080:8080 go-tpl
```

## Configuration

Configuration is loaded from `config/conf.yml` with environment variable overrides:

### Database Configuration
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- Default: MySQL on 127.0.0.1:3306

### RDB Configuration
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`
- Default: RDB on 127.0.0.1:6379

### JWT Configuration
- `JWT_SECRET`, `JWT_EXPIRE_TIME`
- Default: 24-hour expiration

### Server Configuration
- `SERVER_PORT`, `SERVER_MODE`
- Default: Port 8080, debug mode

## Key Patterns

### Dependency Injection with Wire
- Services are initialized in `logic/init.go` using Wire
- Clear dependency flow: infrastructure → logic → web
- Compile-time dependency verification

### Service Layer Pattern
- Each domain has dedicated service with business logic
- Services accept infrastructure dependencies (DB, RDB)
- Shared utilities in `logic/shared/`

### HTTP Handler Pattern
- Handlers organized by domain in `web/rest/{domain}/`
- Each domain has `route.go` for registration and `handler.go` for request handling
- Standardized JSON responses: `base.OK()`, `base.OKWithData()`, `base.Error()`

### Error Handling
- Centralized error definitions in `logic/shared/errors.go`
- Business error codes with user-friendly messages
- Consistent error response format across APIs

### Database Pattern
- Models in `logic/{domain}/model.go` with GORM annotations
- Transaction support via `infra/dbs/transaction.go`
- Global DB instance accessible via `infra.DB`

### Logging Pattern
- Request-scoped logging: `logx.Info(c, message)`
- Infrastructure logging: standard zap logger
- Structured logging with context preservation

## Adding New Features

1. **New Domain**:
   - Create `logic/{domain}/model.go` and `service.go`
   - Create `web/rest/{domain}/handler.go` and `route.go`
   - Add to Wire configuration in `logic/wire.go`

2. **Service Registration**:
   - Add service struct to `logic/init.go`
   - Include in Wire provider set
   - Generate dependencies with `go generate ./logic`

3. **API Registration**:
   - Add routes in `web/rest/{domain}/route.go`
   - Register domain in `web/router.go`

4. **Database Models**:
   - Define models with GORM annotations
   - Include soft delete patterns where applicable
   - Add validation tags

## Dependencies

### Core Dependencies
- `github.com/gin-gonic/gin v1.11.0` - HTTP web framework
- `gorm.io/gorm v1.31.0` - ORM library
- `gorm.io/driver/mysql v1.6.0` - MySQL database driver
- `github.com/redis/go-redis/v9 v9.16.0` - RDB client
- `go.uber.org/zap v1.27.0` - Structured logging

### Authentication & Security
- `github.com/golang-jwt/jwt/v5 v5.3.0` - JWT token handling
- `golang.org/x/crypto v0.43.0` - Cryptographic functions (bcrypt)

### Dependency Injection
- `github.com/google/wire v0.7.0` - Compile-time dependency injection

### Configuration & Utilities
- `github.com/goccy/go-yaml v1.18.0` - YAML parsing
- `gopkg.in/natefinch/lumberjack.v2 v2.2.1` - Log rotation

### Monitoring & Testing
- `github.com/prometheus/client_golang v1.23.2` - Metrics collection
- `github.com/stretchr/testify v1.11.1` - Testing framework
- `github.com/DATA-DOG/go-sqlmock v1.5.2` - Database mocking

### Performance Profiling
- `github.com/felixge/fgprof v0.9.5` - Continuous profiling