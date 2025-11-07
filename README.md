# Go-Tpl - Production-Grade Go Web Application Template

[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)
[![Gin](https://img.shields.io/badge/Gin-HTTP%20Framework-green.svg)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/GORM-ORM-orange.svg)](https://gorm.io/)
[![Redis](https://img.shields.io/badge/Redis-Cache-red.svg)](https://redis.io/)
[![JWT](https://img.shields.io/badge/JWT-Authentication-yellow.svg)](https://jwt.io/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

基于清洁架构原则和领域驱动设计构建的生产级 Go Web 应用模板，具备完整的用户管理、基于角色的访问控制（RBAC）、依赖注入、监控和日志功能。

## ✨ 特性

- **清洁架构**: 严格分层架构，关注点分离（infra、logic、web）
- **领域驱动设计**: 按业务域组织代码（user、role、permission）
- **依赖注入**: 使用 Wire 进行编译时依赖注入
- **用户管理**: 完整的 CRUD 操作、状态管理、角色分配
- **基于角色的访问控制 (RBAC)**: 灵活的权限系统，支持细粒度权限控制
- **数据库**: GORM with MySQL 驱动，完整的事务支持
- **缓存**: Redis 集成，支持缓存和会话管理
- **身份认证**: 基于 JWT 的身份认证，bcrypt 密码加密
- **日志记录**: 基于 Zap 的高性能结构化日志和请求上下文
- **监控**: Prometheus 指标收集和 pprof 性能分析
- **配置管理**: 基于 YAML 的配置系统，支持环境变量覆盖
- **Docker**: 生产就绪的 Docker 配置和部署
- **API 文档**: 完整的 REST API 和标准化响应格式
- **测试支持**: 集成测试框架和数据库模拟

## 🏗️ 架构

### 目录结构

```
go-tpl/
├── cmd/                    # 应用程序入口点
│   └── run.go             # 主应用程序引导和初始化
├── config/                # 配置文件
│   └── conf.yml           # 主配置文件（YAML格式）
├── infra/                 # 基础设施层
│   ├── config/            # 配置管理和环境变量处理
│   ├── dbs/               # 数据库连接、事务管理
│   │   ├── db.go          # MySQL/GORM 初始化
│   │   ├── transaction.go # 事务管理
│   │   └── transaction_test.go
│   ├── jwt/               # JWT 令牌生成和验证
│   ├── logging/           # 日志基础设施
│   │   └── logx/          # 自定义日志工具和中间件
│   └── monitor/           # 监控和性能分析
│       └── monitor.go     # Prometheus 指标设置
├── logic/                 # 业务逻辑层
│   ├── user/              # 用户域
│   │   ├── model.go       # 用户数据模型
│   │   └── service.go     # 用户业务逻辑
│   ├── role/              # 角色域
│   ├── permission/        # 权限域
│   ├── shared/            # 共享工具
│   │   ├── consts.go      # 应用常量
│   │   ├── errors.go      # 错误处理
│   │   └── pagination.go # 分页工具
│   ├── init.go            # 逻辑层初始化
│   ├── wire.go            # Wire 依赖注入配置
│   └── wire_gen.go        # 生成的 Wire 代码
├── web/                   # Web 层
│   ├── base/              # 基础工具
│   │   └── render.go      # 标准化响应渲染
│   ├── middleware/        # HTTP 中间件
│   ├── rest/              # REST API 处理器
│   │   ├── user/          # 用户 API 端点
│   │   ├── role/          # 角色 API 端点
│   │   └── permission/    # 权限 API 端点
│   ├── types/             # 请求/响应类型定义
│   └── router.go          # 路由配置和注册
├── scripts/               # 实用脚本
├── Dockerfile             # Docker 配置
├── go.mod                 # Go 模块依赖
├── go.sum                 # 依赖校验和
├── CLAUDE.md              # Claude Code 开发指南
└── README.md              # 本文件
```

### 核心组件

- **主入口**: `cmd/run.go` - 应用程序引导，遵循 4 步初始化流程
- **配置管理**: `infra/config/config.go` - 基于 YAML 的配置，支持环境变量覆盖
- **数据库层**: GORM with MySQL 驱动，通过 `infra.DB` 访问，支持事务
- **缓存层**: Redis 客户端，通过 `infra.Redis` 访问
- **日志系统**: 基于 Zap 的高性能结构化日志和自定义 logx 工具
- **依赖注入**: Wire 编译时依赖注入，清晰的依赖流向
- **监控系统**: Prometheus 指标收集和 pprof 性能分析
- **身份认证**: JWT-based 认证，bcrypt 密码加密
- **HTTP 框架**: Gin Web 框架和自定义响应工具

### 初始化流程

1. **基础设施初始化** (`infra.Init()`)
   - 加载配置文件和环境变量
   - 设置日志系统
   - 初始化数据库连接和 Redis
   - 配置 JWT 和其他基础设施组件

2. **业务逻辑初始化** (`logic.Init()`)
   - 使用 Wire 进行依赖注入
   - 初始化各域服务 (User, Role, Permission)
   - 设置业务逻辑层依赖关系

3. **Web 层设置** (`web.SetupRouter()`)
   - 配置 HTTP 路由和中间件
   - 注册各域 API 端点
   - 设置请求处理器

4. **监控启动** (`monitor.SetupMetrics()`)
   - 启动 Prometheus 指标收集
   - 配置性能分析端点
   - 开始 HTTP 服务器 (端口 8080)

### 技术栈

#### 核心框架
- **Go 1.25** - 编程语言
- **Gin v1.11.0** - HTTP Web 框架
- **GORM v1.31.0** - ORM 库

#### 数据存储
- **MySQL** - 主数据库
- **Redis v9.16.0** - 缓存和会话存储
- **GORM MySQL Driver v1.6.0** - MySQL 数据库驱动

#### 身份认证与安全
- **JWT v5.3.0** - 令牌认证
- **bcrypt** - 密码哈希

#### 依赖注入与配置
- **Wire v0.7.0** - 编译时依赖注入
- **goccy/go-yaml v1.18.0** - YAML 配置解析

#### 日志与监控
- **Zap v1.27.0** - 高性能结构化日志
- **Prometheus v1.23.2** - 指标收集
- **fgprof v0.9.5** - 性能分析

#### 测试与开发工具
- **testify v1.11.1** - 测试框架
- **go-sqlmock v1.5.2** - 数据库模拟
- **pprof** - 性能分析

## 🚀 快速开始

### 前置要求

- Go 1.25+
- MySQL 8.0+
- Redis 6.0+
- Docker（可选）

### 安装

1. **克隆仓库**
```bash
git clone <repository-url>
cd go-tpl
```

2. **安装依赖**
```bash
go mod tidy
```

3. **配置应用程序**
```bash
cp config/conf.yml.example config/conf.yml
# 编辑 config/conf.yml 设置您的数据库和 Redis 配置
```

4. **运行应用程序**
```bash
# 直接运行
go run .

# 或者构建后运行
go build -o myapp .
./myapp
```

应用程序将在 `http://localhost:8080` 启动

### Docker 设置

1. **构建 Docker 镜像**
```bash
docker build -t go-tpl .
```

2. **使用 Docker 运行**
```bash
docker run -p 8080:8080 -v $(pwd)/config:/app/config go-tpl
```

## ⚙️ 配置

配置从 `config/conf.yml` 加载，并支持环境变量覆盖：

### 数据库配置
```yaml
database:
  host: 127.0.0.1
  port: 3306
  dbname: app_db
  user: root
  password: 123456
```

环境变量：
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`

### Redis 配置
```yaml
redis:
  host: 127.0.0.1
  port: 6379
  password: 123456
  db: 0
```

环境变量：
- `REDIS_HOST`
- `REDIS_PORT`
- `REDIS_PASSWORD`
- `REDIS_DB`

### JWT 配置
```yaml
jwt:
  secret: your-secret-key-change-in-production
  expire_time: 86400  # 24小时，单位：秒
```

环境变量：
- `JWT_SECRET`
- `JWT_EXPIRE_TIME`

### 服务器配置
```yaml
server:
  port: 8080
  mode: debug  # debug, release, test
```

环境变量：
- `SERVER_PORT`
- `SERVER_MODE`

## 📚 API 文档

### 基础 URL
```
http://localhost:8080/api
```

### 标准响应格式
所有 API 响应都遵循统一的 JSON 格式：
```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

- `code`: 业务状态码（0 表示成功，其他表示各种业务错误）
- `msg`: 响应消息
- `data`: 响应数据（成功时返回具体数据，失败时为 null）

### 身份认证头
对于受保护的端点，在 Authorization 头中包含 JWT 令牌：
```
Authorization: Bearer <your-jwt-token>
```

## 👥 用户管理 API

### 1. 获取用户列表
- **URL**: `POST /api/user/list`
- **Method**: `POST`
- **Body**:
```json
{
  "page": 1,
  "page_size": 10,
  "username": "用户名(可选)",
  "email": "邮箱(可选)",
  "status": 1
}
```
- **Response**:
```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "total": 100,
    "list": [
      {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com",
        "status": 1,
        "created_at": "2025-10-22T10:00:00Z",
        "updated_at": "2025-10-22T10:00:00Z"
      }
    ]
  }
}
```

### 2. 获取单个用户
- **URL**: `GET /api/user/{id}`
- **Method**: `GET`
- **Response**:
```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "status": 1,
    "created_at": "2025-10-22T10:00:00Z",
    "updated_at": "2025-10-22T10:00:00Z"
  }
}
```

### 3. 创建用户
- **URL**: `POST /api/user`
- **Method**: `POST`
- **Body**:
```json
{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "password123"
}
```

### 4. 更新用户
- **URL**: `PUT /api/user/{id}`
- **Method**: `PUT`
- **Body**:
```json
{
  "username": "updateduser",
  "email": "updated@example.com",
  "password": "newpassword",
  "status": 1
}
```

### 5. 删除用户
- **URL**: `DELETE /api/user/{id}`
- **Method**: `DELETE`

### 6. 更新用户状态
- **URL**: `PUT /api/user/{id}/status`
- **Method**: `PUT`
- **Body**:
```json
{
  "status": 0
}
```

### 7. 获取用户角色
- **URL**: `GET /api/user/{id}/roles`
- **Method**: `GET`
- **Response**:
```json
{
  "code": 0,
  "msg": "ok",
  "data": [1, 2, 3]
}
```

### 8. 分配用户角色
- **URL**: `PUT /api/user/{id}/roles`
- **Method**: `PUT`
- **Body**:
```json
{
  "role_ids": [1, 2, 3]
}
```

## 🎭 角色管理 API

### 1. 获取角色列表
- **URL**: `POST /api/role/list`
- **Method**: `POST`
- **Body**:
```json
{
  "page": 1,
  "page_size": 10,
  "name": "角色名(可选)",
  "status": 1
}
```

### 2. 获取单个角色
- **URL**: `GET /api/role/{id}`
- **Method**: `GET`

### 3. 创建角色
- **URL**: `POST /api/role`
- **Method**: `POST`
- **Body**:
```json
{
  "name": "新角色",
  "description": "角色描述"
}
```

### 4. 更新角色
- **URL**: `PUT /api/role/{id}`
- **Method**: `PUT`
- **Body**:
```json
{
  "name": "更新的角色名",
  "description": "更新的描述",
  "status": 1
}
```

### 5. 删除角色
- **URL**: `DELETE /api/role/{id}`
- **Method**: `DELETE`

### 6. 更新角色状态
- **URL**: `PUT /api/role/{id}/status`
- **Method**: `PUT`
- **Body**:
```json
{
  "status": 0
}
```

### 7. 获取角色权限
- **URL**: `GET /api/role/{id}/permissions`
- **Method**: `GET`

### 8. 分配角色权限
- **URL**: `PUT /api/role/{id}/permissions`
- **Method**: `PUT`
- **Body**:
```json
{
  "permission_ids": [1, 2, 3, 4]
}
```

### 9. 获取角色用户
- **URL**: `GET /api/role/{id}/users`
- **Method**: `GET`

## 🔐 权限管理 API

### 1. 获取权限列表
- **URL**: `POST /api/permission/list`
- **Method**: `POST`
- **Body**:
```json
{
  "page": 1,
  "page_size": 10,
  "code": "权限代码(可选)",
  "name": "权限名(可选)",
  "module": "模块名(可选)",
  "status": 1
}
```

### 2. 获取单个权限
- **URL**: `GET /api/permission/{id}`
- **Method**: `GET`

### 3. 创建权限
- **URL**: `POST /api/permission`
- **Method**: `POST`
- **Body**:
```json
{
  "code": "user:create",
  "name": "创建用户",
  "description": "允许创建新用户",
  "module": "user"
}
```

### 4. 更新权限
- **URL**: `PUT /api/permission/{id}`
- **Method**: `PUT`
- **Body**:
```json
{
  "code": "user:update",
  "name": "更新的权限名",
  "description": "更新的描述",
  "module": "user",
  "status": 1
}
```

### 5. 删除权限
- **URL**: `DELETE /api/permission/{id}`
- **Method**: `DELETE`

### 6. 更新权限状态
- **URL**: `PUT /api/permission/{id}/status`
- **Method**: `PUT`
- **Body**:
```json
{
  "status": 0
}
```

### 7. 获取所有模块
- **URL**: `GET /api/permission/modules`
- **Method**: `GET`
- **Response**:
```json
{
  "code": 0,
  "msg": "ok",
  "data": ["user", "role", "permission", "system"]
}
```

### 8. 获取权限角色
- **URL**: `GET /api/permission/{id}/roles`
- **Method**: `GET`

## 🔧 状态码与错误处理

### 业务错误码
- `0`: 成功
- `2001`: 用户不存在
- `2002`: 用户名已存在
- `2003`: 邮箱已存在
- `2004`: 密码格式错误
- `3001`: 角色不存在
- `3002`: 角色名已存在
- `3003`: 角色正在使用中
- `4001`: 权限不存在
- `4002`: 权限代码已存在
- `5001`: 参数错误
- `5002`: 记录不存在
- `500`: 服务器内部错误（通用系统错误）

### HTTP 状态码
所有 API 响应都返回 HTTP 200 状态码。业务状态通过响应体中的 `code` 字段区分。

### 状态值
- `1`: 正常/启用
- `0`: 禁用

### 错误响应示例

用户不存在：
```json
{
  "code": 2001,
  "msg": "用户不存在",
  "data": null
}
```

用户名已存在：
```json
{
  "code": 2002,
  "msg": "用户名已存在",
  "data": null
}
```

## 📝 重要说明

1. **密码安全**: 密码在创建和更新时会自动进行 bcrypt 加密
2. **软删除**: 删除操作为软删除，数据不会物理删除
3. **时间格式**: 所有时间字段使用 ISO 8601 格式
4. **分页**: 分页查询使用基于索引的分页方式
5. **密码隐藏**: 用户密码字段在返回时会隐藏（不返回到前端）
6. **JWT 令牌**: 对于受保护的端点，在 Authorization 头中包含 JWT 令牌
7. **环境变量**: 可以使用环境变量覆盖配置

## 🛠️ 开发

### 构建和运行
```bash
# 构建应用程序
go build -o myapp .

# 直接运行
go run .

# 运行特定模块
go run cmd/run.go

# 生成 Wire 依赖
go generate ./logic

# 使用 Wire 生成依赖
wire ./logic/
```

### 测试
```bash
# 运行所有测试
go test ./...

# 运行测试并生成覆盖率报告
go test -cover ./...

# 运行特定包的测试
go test ./infra/dbs/...
go test ./logic/user/...

# 运行基准测试
go test -bench=. ./...
```

### 开发工具
```bash
# 格式化代码
go fmt ./...

# 检查潜在问题
go vet ./...

# 获取依赖
go mod tidy

# 下载依赖
go mod download

# 更新依赖
go get -u ./...

# 查看依赖图
go mod graph
```

### 监控与调试

应用程序包含内置的监控端点：

- **指标**: `http://localhost:8080/metrics` (Prometheus 格式)
- **性能分析**: `http://localhost:8080/debug/pprof/` (Go pprof)

## 🏛️ 领域模型

### 用户管理 (User Domain)
- **完整的 CRUD 操作**: 创建、读取、更新、删除用户
- **状态管理**: 支持用户启用/禁用状态
- **密码安全**: 使用 bcrypt 进行密码哈希
- **角色关联**: 用户可以分配多个角色
- **数据验证**: 用户名、邮箱格式验证

### 角色系统 (Role Domain)
- **角色管理**: 创建、更新、删除角色
- **权限分配**: 角色可以拥有多个权限
- **层级关系**: 支持角色间的层级关系
- **状态控制**: 角色启用/禁用功能

### 权限系统 (Permission Domain)
- **细粒度权限**: 基于 资源-操作 的权限模型
- **模块化组织**: 权限按业务模块分组
- **动态权限**: 支持运行时权限检查
- **权限继承**: 通过角色继承权限

### 共享组件 (Shared Components)
- **常量定义**: 统一的业务常量
- **错误处理**: 标准化的错误码和消息
- **分页工具**: 通用的分页查询支持

## 🔧 添加新功能

### 1. 创建新域
```bash
# 创建领域目录
mkdir -p logic/{newdomain}
mkdir -p web/rest/{newdomain}
```

### 2. 实现数据模型
在 `logic/{newdomain}/model.go` 中定义：
```go
type NewDomain struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"size:100;not null" json:"name"`
    Status    int       `gorm:"default:1" json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

### 3. 实现服务层
在 `logic/{newdomain}/service.go` 中：
```go
type NewDomainSvc struct {
    db *gorm.DB
    redis *redis.Client
}

// 实现业务逻辑方法
func (s *NewDomainSvc) Create(req *CreateNewDomainReq) error {
    // 业务逻辑实现
}
```

### 4. 配置依赖注入
在 `logic/wire.go` 中添加：
```go
var NewDomainSet = wire.NewSet(
    // 添加服务提供者
    wire.Bind((*NewDomainInterface)(nil), (*NewDomainSvc)(nil)),
)
```

### 5. 生成 Wire 代码
```bash
go generate ./logic
# 或
wire ./logic/
```

### 6. 实现 HTTP 处理器
在 `web/rest/{newdomain}/handler.go` 中实现 API 处理器

### 7. 配置路由
在 `web/rest/{newdomain}/route.go` 中定义路由，并在 `web/router.go` 中注册

### 8. 添加测试
为每个层次添加单元测试和集成测试

## 📦 依赖

### 核心依赖
- `github.com/gin-gonic/gin v1.11.0` - HTTP Web 框架
- `gorm.io/gorm v1.31.0` - ORM 库
- `gorm.io/driver/mysql v1.6.0` - MySQL 数据库驱动
- `github.com/redis/go-redis/v9 v9.16.0` - Redis 客户端
- `go.uber.org/zap v1.27.0` - 高性能结构化日志

### 身份认证与安全
- `github.com/golang-jwt/jwt/v5 v5.3.0` - JWT 令牌处理
- `golang.org/x/crypto v0.43.0` - 加密函数 (bcrypt)

### 依赖注入与配置
- `github.com/google/wire v0.7.0` - 编译时依赖注入
- `github.com/goccy/go-yaml v1.18.0` - YAML 配置解析

### 监控与日志
- `github.com/prometheus/client_golang v1.23.2` - Prometheus 指标收集
- `gopkg.in/natefinch/lumberjack.v2 v2.2.1` - 日志轮转

### 测试工具
- `github.com/stretchr/testify v1.11.1` - 测试框架
- `github.com/DATA-DOG/go-sqlmock v1.5.2` - 数据库模拟
- `go.uber.org/mock v0.6.0` - Mock 生成工具

### 性能分析
- `github.com/felixge/fgprof v0.9.5` - 连续性能分析
- `github.com/google/pprof` - Go 性能分析工具

### 间接依赖
项目还包含多个间接依赖，用于支持：
- JSON 序列化/反序列化 (`github.com/goccy/go-json`)
- HTTP/2 和 QUIC 支持
- 数据验证 (`github.com/go-playground/validator/v10`)
- 国际化支持 (`github.com/go-playground/locales`)
- 配置文件解析 (`github.com/pelletier/go-toml/v2`)
- 等

所有依赖都在 `go.mod` 文件中定义，确保版本锁定和构建可重现性。

## 🤝 贡献

1. Fork 仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

## 📄 许可证

本项目基于 MIT 许可证 - 详情请查看 [LICENSE](LICENSE) 文件。

## 🆘 支持

如有任何问题或疑问，请在仓库中创建 issue。

---

**使用 ❤️ 和清洁架构原则构建的 Go 项目**