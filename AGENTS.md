# AGENTS.md

## Development Commands

### Build & Run
```bash
go build -o myapp .
go run .
go run cmd/run.go
```

### Testing
```bash
go test ./...                    # Run all tests
go test -cover ./...             # Run with coverage
go test ./infra/dbs/...          # Run specific package tests
go test -run TestSpecificFunc    # Run single test function
```

### Code Quality
```bash
go fmt ./...                     # Format code
go vet ./...                      # Vet for issues
go mod tidy                       # Clean dependencies
go generate ./logic               # Generate Wire dependencies
wire ./logic/                     # Alternative Wire generation
```

## Code Style Guidelines

### Imports
- Group imports: stdlib, third-party, project modules
- Use absolute imports: "go-tpl/infra/logging"
- No unused imports

### Naming Conventions
- PascalCase for exported types, functions, constants
- camelCase for unexported
- Use descriptive names: `UserService` not `UsrSvc`
- Error variables: `ErrUserExists` not `UserErr`

### Error Handling
- Use centralized errors from `logic/shared/errors.go`
- Return errors as last return value
- Log errors with context: `logging.Errorf(c, "message: %v", err)`
- Use `base.FailWithError(c, err)` for API responses

### Patterns
- Follow 3-layer architecture: infra → logic → web
- Use Wire for dependency injection
- Handlers: bind request → call service → return response
- Services: accept context, return domain models
- Models: use GORM annotations, include soft delete patterns