# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a production-grade Go web application template built with Gin framework that follows clean architecture principles and domain-driven design. The project features a complete RBAC (Role-Based Access Control) system with user, role, and permission management, following a three-layer architecture: infrastructure (infra), business logic (logic), and web layer (web). The application includes JWT-based authentication with access/refresh token pairs, comprehensive monitoring, and Docker deployment support.

## Tech Stack

- **Go** - Programming language (1.25+)
- **Gin** - HTTP web framework for REST APIs
- **GORM** - ORM with MySQL driver
- **Redis** - Caching and session management
- **JWT** - Authentication tokens with access/refresh token pairs
- **Wire** - Compile-time dependency injection
- **Viper** - Configuration management with environment variable support
- **Zerolog** - High-performance structured logging
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
├── logger/            # Logging infrastructure
│   └── logx/          # Custom logging utilities
└── monitor/           # Monitoring and metrics

logic/                 # Business logic layer
├── auth/              # Authentication service (register, login, refresh)
├── user/              # User domain (model, service)
├── role/              # Role domain
├── permission/        # Permission domain
├── shared/            # Shared utilities (consts, errors, pagination)
├── init.go            # Logic layer initialization
├── wire.go            # Wire dependency injection setup
└── wire_gen.go        # Generated wire code

web/                   # Web layer
├── app.go             # Application setup with graceful shutdown
├── base/              # Base utilities (response rendering)
├── middleware/        # HTTP middleware
├── rest/              # REST API handlers
│   ├── user/          # User API endpoints
│   ├── role/          # Role API endpoints
│   └── permission/    # Permission API endpoints
├── types/             # Request/response types
└── router.go          # Route configuration

scripts/               # Utility scripts
├── docker-compose.yaml    # Docker Compose configuration
└── Dockerfile          # Multi-stage Docker build
```

### Key Components
- **Main Entry**: `cmd/run.go` - Application bootstrap following 4-step initialization
- **Configuration**: `infra/config/config.go` - Viper-based config with YAML and environment variable support
- **Database**: GORM with MySQL driver, accessible via `infra.DB` with transaction support
- **Caching**: Redis client, accessible via `infra.RDB`
- **Logging**: Zerolog-based structured logging with custom logx utilities and middleware
- **Dependency Injection**: Wire-based compile-time DI for clean dependencies
- **Monitoring**: Prometheus metrics and pprof profiling support
- **Authentication**: JWT-based auth with access/refresh token pairs and bcrypt password hashing
- **Application Setup**: `web/app.go` - Application initialization with graceful shutdown

### Initialization Flow
1. `infra.Init()` - Load config, setup logging, initialize DB/RDB/JWT
2. `logic.Init()` - Initialize service layer with Wire dependency injection
3. `web.SetupRouter()` - Setup HTTP routes, middleware, and handlers
4. `monitor.SetupMetrics(app)` - Start Prometheus metrics collection
5. Start HTTP server on port 8080

## Domain Models

### Authentication Service
- User registration with automatic token generation
- User login with credential validation
- Token refresh mechanism with refresh token rotation
- Token pair generation (access + refresh tokens)

### User Management
- Complete CRUD operations with status management (active/inactive)
- Password encryption using bcrypt
- Role assignment and RBAC integration
- Login validation and credential checking

### Role System
- Hierarchical role management
- Permission assignment capabilities
- Status management
- User-role many-to-many relationships

### Permission System
- Granular permission structure with module organization
- Resource-action based permissions
- Integration with role-based access control
- Dynamic permission checking

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

Configuration is managed by Viper and loaded from `config/conf.yml` with environment variable overrides:

### Database Configuration
```yaml
database:
  host: 127.0.0.1
  port: 3306
  dbname: app_db
  user: root
  password: 123456
```
Environment variables: `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_USER`, `DATABASE_PASSWORD`, `DATABASE_NAME`

### Redis Configuration
```yaml
redis:
  host: 127.0.0.1
  port: 6379
  password: ""
  db: 0
```
Environment variables: `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`

### JWT Configuration
```yaml
jwt:
  secret: your-secret-key
  expire_time: 7200        # access_token expiry (seconds)
  refresh_expire_time: 604800  # refresh_token expiry (seconds)
```
Environment variables: `JWT_SECRET`, `JWT_EXPIRE_TIME`, `JWT_REFRESH_EXPIRE_TIME`

### Server Configuration
```yaml
server:
  port: 8080
  mode: debug  # debug, release, test
```
Environment variables: `SERVER_PORT`, `SERVER_MODE`

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
- `github.com/gin-gonic/gin` - HTTP web framework
- `gorm.io/gorm` - ORM library
- `gorm.io/driver/mysql` - MySQL database driver
- `github.com/redis/go-redis/v9` - Redis client
- `github.com/rs/zerolog` - Structured logging

### Authentication & Security
- `github.com/golang-jwt/jwt/v5` - JWT token handling with access/refresh tokens
- `golang.org/x/crypto` - Cryptographic functions (bcrypt)

### Dependency Injection & Configuration
- `github.com/google/wire` - Compile-time dependency injection
- `github.com/spf13/viper` - Configuration management with environment variable support
- `github.com/goccy/go-yaml` - YAML parsing

### Monitoring & Utilities
- `github.com/prometheus/client_golang` - Metrics collection
- `gopkg.in/natefinch/lumberjack.v2` - Log rotation

### Testing & Development
- `github.com/stretchr/testify` - Testing framework
- `github.com/DATA-DOG/go-sqlmock` - Database mocking
- `github.com/felixge/fgprof` - Continuous profiling

## Security Features

### Authentication & Authorization
- JWT-based authentication with access token (2h default) and refresh token (7d default)
- bcrypt password hashing with configurable cost factor
- Role-based access control (RBAC) with fine-grained permissions
- Middleware for protected routes with token validation

### Data Protection
- Soft delete patterns for sensitive data
- Input validation on all API endpoints
- Environment variable support for sensitive configuration
- CORS middleware configuration
- SQL injection prevention through GORM ORM

## Monitoring & Observability

### Metrics Collection
- Prometheus metrics at `/metrics` endpoint
- Custom business metrics tracking
- HTTP request metrics (count, duration, status codes)
- Application performance monitoring

### Logging
- Structured logging with Zerolog
- Request-scoped logging with correlation IDs
- Log levels: debug, info, warn, error
- Log rotation with Lumberjack
- Context preservation across service calls

### Profiling
- pprof profiling support at `/debug/pprof/`
- CPU, memory, goroutine, and block profiling
- Continuous profiling with fgprof
- Production-safe profiling endpoints