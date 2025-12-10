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
go test -bench=. ./...           # Run benchmarks
```

### Code Quality
```bash
go fmt ./...                     # Format code
go vet ./...                      # Vet for issues
go mod tidy                       # Clean dependencies
go generate ./logic               # Generate Wire dependencies
wire ./logic/                     # Alternative Wire generation
go mod download                   # Download dependencies
go mod graph                      # View dependency graph
```

### Docker & Deployment
```bash
docker build -t go-tpl .                    # Build Docker image
docker run -p 8080:8080 go-tpl              # Run with Docker
docker-compose up -d                        # Run with Docker Compose
```

## Code Style Guidelines

### Project Structure
- Follow clean architecture: infra (infrastructure) → logic (business) → web (presentation)
- Domain-driven design: organize by business domains (user, role, permission, auth)
- Use absolute imports: "go-tpl/infra/config"
- Keep layers separated with clear interfaces

### Import Organization
```go
import (
    // Standard library
    "context"
    "time"

    // Third-party libraries
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    // Project modules
    "go-tpl/infra/config"
    "go-tpl/logic/shared"
)
```

### Naming Conventions
- **Packages**: lowercase, single word when possible (`user`, `auth`, `config`)
- **Types**: PascalCase for exported (`UserService`, `CreateUserReq`)
- **Functions**: PascalCase for exported, camelCase for unexported (`CreateUser`, `validateInput`)
- **Variables**: camelCase (`userID`, `accessToken`)
- **Constants**: PascalCase or UPPER_SNAKE_CASE (`MaxRetryCount`, `DEFAULT_TIMEOUT`)
- **Errors**: `Err` prefix for exported errors (`ErrUserExists`, `ErrInvalidToken`)

### Service Layer Patterns
```go
// Service structure
type UserService struct {
    db    *gorm.DB
    redis *redis.Client
}

// Method signature pattern
func (s *UserService) Create(ctx context.Context, req *CreateUserReq) (*User, error) {
    // 1. Validate input
    if err := s.validateCreateRequest(req); err != nil {
        return nil, err
    }

    // 2. Process business logic
    user := &User{
        Username: req.Username,
        Email:    req.Email,
    }

    // 3. Persist data
    if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
        return nil, err
    }

    return user, nil
}
```

### HTTP Handler Patterns
```go
// Handler structure
type UserHandler struct {
    userSvc *user.Service
}

// Handler method pattern
func (h *UserHandler) Create(c *gin.Context) {
    // 1. Bind and validate request
    var req types.CreateUserReq
    if err := c.ShouldBindJSON(&req); err != nil {
        base.FailWithValidation(c, err)
        return
    }

    // 2. Call service
    user, err := h.userSvc.Create(c.Request.Context(), req)
    if err != nil {
        base.FailWithError(c, err)
        return
    }

    // 3. Return response
    base.OKWithData(c, user)
}
```

### Database Patterns
```go
// Model definition
type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
    Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
    Password  string         `gorm:"size:255;not null" json:"-"`
    Status    int            `gorm:"default:1" json:"status"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Query patterns
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*User, error) {
    var user User
    err := r.db.WithContext(ctx).First(&user, id).Error
    return &user, err
}

func (r *UserRepository) FindWithPagination(ctx context.Context, req *ListRequest) ([]*User, int64, error) {
    var users []*User
    var total int64

    query := r.db.WithContext(ctx).Model(&User{})

    // Apply filters
    if req.Username != "" {
        query = query.Where("username LIKE ?", "%"+req.Username+"%")
    }

    // Count total
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Get paginated results
    offset := (req.Page - 1) * req.PageSize
    err := query.Offset(offset).Limit(req.PageSize).Find(&users).Error

    return users, total, err
}
```

### Error Handling Patterns
```go
// Centralized error definitions (logic/shared/errors.go)
var (
    ErrUserNotFound     = NewBusinessError(2001, "用户不存在")
    ErrUserExists       = NewBusinessError(2002, "用户名已存在")
    ErrInvalidPassword  = NewBusinessError(2003, "密码格式错误")
)

// Service error handling
func (s *UserService) GetByID(ctx context.Context, id uint) (*User, error) {
    var user User
    err := s.db.WithContext(ctx).First(&user, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, shared.ErrUserNotFound
        }
        return nil, err
    }
    return &user, nil
}

// Handler error handling
func (h *UserHandler) Get(c *gin.Context) {
    id := c.GetUint("id")
    user, err := h.userSvc.GetByID(c.Request.Context(), id)
    if err != nil {
        base.FailWithError(c, err)  // Automatically handles business errors
        return
    }
    base.OKWithData(c, user)
}
```

### Configuration Patterns
```go
// Configuration structure (infra/config/config.go)
type Config struct {
    Database DBConfig     `mapstructure:"database"`
    Redis    RedisConfig  `mapstructure:"redis"`
    JWT      JWTConfig    `mapstructure:"jwt"`
    Server   ServerConfig `mapstructure:"server"`
}

// Access configuration
func Serve() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }

    // Use configuration
    db := infra.InitDB(cfg.Database)
}
```

### Logging Patterns
```go
// Request-scoped logging (web layer)
func (h *UserHandler) Create(c *gin.Context) {
    logx.Info(c, "Creating user", "username", req.Username)

    user, err := h.userSvc.Create(c.Request.Context(), req)
    if err != nil {
        logx.Error(c, "Failed to create user", "error", err, "username", req.Username)
        base.FailWithError(c, err)
        return
    }

    logx.Info(c, "User created successfully", "user_id", user.ID)
    base.OKWithData(c, user)
}

// Infrastructure logging (service layer)
func (s *UserService) Create(ctx context.Context, req *CreateUserReq) (*User, error) {
    logger := zerolog.Ctx(ctx)
    logger.Info().Str("username", req.Username).Msg("Creating user in service")

    // Service logic...

    return user, nil
}
```

### Testing Patterns
```go
// Service unit test
func TestUserService_Create(t *testing.T) {
    // Setup
    db, mock := sqlmock.New()
    gormDB, _ := gorm.Open(mysql.New(mysql.Config{
        Conn: db,
    }), &gorm.Config{})

    userSvc := user.NewService(gormDB, nil)

    // Test case
    req := &types.CreateUserReq{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }

    // Mock expectations
    mock.ExpectBegin()
    mock.ExpectExec("INSERT INTO").
        WithArgs(req.Username, req.Email, sqlmock.AnyArg()).
        WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectCommit()

    // Execute
    user, err := userSvc.Create(context.Background(), req)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, req.Username, user.Username)
    assert.NoError(t, mock.ExpectationsWereMet())
}

// Integration test with test database
func TestUserHandler_Create_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    // Setup test dependencies
    userSvc := user.NewService(db, nil)
    userHandler := &UserHandler{userSvc: userSvc}

    // Setup Gin test context
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.POST("/users", userHandler.Create)

    // Test request
    req := map[string]interface{}{
        "username": "testuser",
        "email":    "test@example.com",
        "password": "password123",
    }
    reqBody, _ := json.Marshal(req)

    // Execute request
    w := httptest.NewRecorder()
    httpReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
    router.ServeHTTP(w, httpReq)

    // Assert response
    assert.Equal(t, 200, w.Code)

    var response base.Response
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, float64(0), response["code"])
}
```

### Wire Dependency Injection Patterns
```go
// Wire provider sets (logic/wire.go)
var ProviderSet = wire.NewSet(
    // Infrastructure providers
    config.ProviderSet,
    dbs.ProviderSet,
    jwt.ProviderSet,

    // Service providers
    user.ProviderSet,
    auth.ProviderSet,

    // Web providers
    web.ProviderSet,
)

// Service provider example (logic/user/wire.go)
var ProviderSet = wire.NewSet(
    NewService,
    NewRepository,
    wire.Bind(new(ServiceInterface), new(*Service)),
    wire.Bind(new(RepositoryInterface), new(*Repository)),
)
```

### JWT Authentication Patterns
```go
// Token generation and validation (infra/jwt/jwt.go)
func GenerateTokenPair(userID uint) (*TokenPair, error) {
    // Generate access token (short-lived)
    accessToken, err := generateToken(userID, AccessTokenType, accessExpireTime)
    if err != nil {
        return nil, err
    }

    // Generate refresh token (long-lived)
    refreshToken, err := generateToken(userID, RefreshTokenType, refreshExpireTime)
    if err != nil {
        return nil, err
    }

    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}

// Authentication middleware (web/middleware/auth.go)
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            base.FailWithCode(c, 1000, "no token")
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := jwt.ParseToken(tokenString)
        if err != nil {
            base.FailWithCode(c, 1001, "invalid token")
            c.Abort()
            return
        }

        // Check token type
        if claims.Type != jwt.AccessTokenType {
            base.FailWithCode(c, 1001, "invalid token type")
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

### Performance Optimization Patterns
```go
// Database query optimization
func (r *UserRepository) FindWithRoles(ctx context.Context, userID uint) (*User, error) {
    var user User
    err := r.db.WithContext(ctx).
        Preload("Roles", "status = ?", 1).  // Only active roles
        Preload("Roles.Permissions").       // Eager load permissions
        First(&user, userID).Error
    return &user, err
}

// Redis caching pattern
func (s *UserService) GetByID(ctx context.Context, id uint) (*User, error) {
    // Try cache first
    cacheKey := fmt.Sprintf("user:%d", id)
    cached, err := s.redis.Get(ctx, cacheKey).Result()
    if err == nil {
        var user User
        json.Unmarshal([]byte(cached), &user)
        return &user, nil
    }

    // Cache miss, get from database
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Update cache
    userJSON, _ := json.Marshal(user)
    s.redis.Set(ctx, cacheKey, userJSON, 5*time.Minute)

    return user, nil
}
```

## Best Practices

### Security
1. **Never log passwords or sensitive data**
2. **Use parameterized queries** (handled automatically by GORM)
3. **Validate all input** at the handler layer
4. **Use HTTPS in production** for all API endpoints
5. **Implement rate limiting** for sensitive endpoints
6. **Store secrets in environment variables**, not in code

### Performance
1. **Use database transactions** for multi-step operations
2. **Implement proper database indexing** for frequently queried fields
3. **Use Redis caching** for frequently accessed data
4. **Optimize SQL queries** to avoid N+1 problems
5. **Use connection pooling** (handled automatically by GORM)
6. **Monitor performance metrics** with Prometheus

### Maintainability
1. **Keep functions small** and focused on single responsibility
2. **Write comprehensive tests** for business logic
3. **Use consistent error handling** throughout the application
4. **Document complex business logic** with clear comments
5. **Follow semantic versioning** for releases
6. **Keep dependencies up to date** regularly

### Monitoring & Debugging
1. **Use structured logging** with context information
2. **Implement health checks** for external dependencies
3. **Monitor key business metrics** (user registrations, logins, etc.)
4. **Set up alerts** for error rates and response times
5. **Use distributed tracing** for complex request flows
6. **Regular review of logs** for unusual patterns