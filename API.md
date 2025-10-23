# API 接口文档

## 用户管理 API

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

## 角色管理 API

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

## 权限管理 API

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

## 响应状态码说明

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

### HTTP状态码
所有API响应都返回HTTP 200状态码，业务状态通过响应体中的`code`字段区分。

## 标准响应格式

所有API响应都遵循统一的JSON格式：
```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

- `code`: 业务状态码（0表示成功，其他表示各种业务错误）
- `msg`: 响应消息
- `data`: 响应数据（成功时返回具体数据，失败时为null）

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

## 状态值说明

- `1`: 正常/启用
- `0`: 禁用

## 注意事项

1. 密码在创建和更新时会自动进行bcrypt加密
2. 删除操作为软删除，数据不会物理删除
3. 所有时间字段使用ISO 8601格式
4. 分页查询使用基于索引的分页方式
5. 用户密码字段在返回时会隐藏（不返回到前端）