# Go-Tpl - Go Web Application Template

[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)
[![Gin](https://img.shields.io/badge/Gin-HTTP%20Framework-green.svg)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/GORM-ORM-orange.svg)](https://gorm.io/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

基于清洁架构原则构建的生产级 Go Web 应用模板，具备用户管理、基于角色的访问控制（RBAC）和全面的监控功能。

## ✨ 特性

- **清洁架构**: 关注点分离，分层清晰（infra、logic、web）
- **用户管理**: 完整的 CRUD 操作和角色分配功能
- **基于角色的访问控制 (RBAC)**: 灵活的权限系统
- **数据库**: GORM with MySQL 支持和事务管理
- **缓存**: Redis 集成以优化性能
- **身份认证**: 基于 JWT 的身份认证，安全的令牌处理
- **日志记录**: 基于 Zap 的结构化日志和请求上下文
- **监控**: Prometheus 指标和 pprof 性能分析
- **配置管理**: 基于 YAML 的配置，支持环境变量覆盖
- **Docker**: 生产就绪的 Docker 配置
- **API 文档**: 完整的 REST API 和标准化响应

## 🏗️ 架构

### 目录结构

```
go-tpl/
├── cmd/                    # 应用程序入口点
│   └── run.go             # 主应用程序引导
├── config/                # 配置文件
│   └── conf.yml           # 主配置文件
├── infra/                 # 基础设施层
│   ├── config/            # 配置管理
│   ├── dbs/               # 数据库设置和事务
│   ├── jwt/               # JWT 令牌处理
│   ├── logging/           # 日志工具和中间件
│   └── monitor/           # 监控和分析
├── logic/                 # 业务逻辑层
│   ├── user/              # 用户域逻辑
│   ├── role/              # 角色域逻辑
│   ├── permission/        # 权限域逻辑
│   └── shared/            # 共享工具和常量
├── web/                   # Web 层
│   ├── base/              # 基础工具和响应
│   ├── middleware/        # HTTP 中间件
│   ├── rest/              # REST API 处理器
│   │   ├── user/          # 用户 API 端点
│   │   ├── role/          # 角色 API 端点
│   │   └── permission/    # 权限 API 端点
│   ├── types/             # 请求/响应类型
│   └── router.go          # 路由配置
├── Dockerfile             # Docker 配置
├── go.mod                 # Go 模块文件
└── README.md              # 本文件
```

### 核心组件

- **主入口**: `cmd/run.go` - 应用程序引导和初始化
- **配置**: `infra/config/config.go` - 基于 YAML 的配置，支持环境变量覆盖
- **数据库**: GORM with MySQL 驱动，通过 `infra.DB` 访问
- **缓存**: Redis 客户端，通过 `infra.Redis` 访问
- **日志**: 基于 Zap 的结构化日志和自定义中间件
- **HTTP 框架**: Gin 和自定义响应工具

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

### 构建命令
```bash
# 构建应用程序
go build -o myapp .

# 直接运行
go run .

# 运行特定模块
go run cmd/run.go
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

# 运行测试
go test ./...

# 运行测试并生成覆盖率报告
go test -cover ./...
```

### 监控与调试

应用程序包含内置的监控端点：

- **指标**: `http://localhost:8080/metrics` (Prometheus 格式)
- **性能分析**: `http://localhost:8080/debug/pprof/` (Go pprof)

### 添加新功能

1. **新域**: 创建 `logic/{domain}/` 和 `web/rest/{domain}/` 目录
2. **服务**: 在 `logic/{domain}/service.go` 中添加服务并在 `logic/init.go` 中初始化
3. **HTTP 层**: 在 `web/rest/{domain}/` 中添加 `handler.go` 和 `route.go`
4. **注册**: 在 `web/router.go` 中添加域注册

## 📦 依赖

主要外部依赖：
- `github.com/gin-gonic/gin` - HTTP 框架
- `gorm.io/gorm` - ORM
- `github.com/redis/go-redis/v9` - Redis 客户端
- `go.uber.org/zap` - 结构化日志
- `github.com/goccy/go-yaml` - YAML 解析
- `github.com/prometheus/client_golang` - 指标
- `github.com/golang-jwt/jwt/v5` - JWT 身份认证
- `golang.org/x/crypto` - 加密函数

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