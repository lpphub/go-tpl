# 项目说明文档

## 一、项目结构

```
go-tpl/
├── cmd/               # 程序启动命令
│   └── run.go         
├── config/            # 配置文件
│   └── config.yml     
├── server/            # web服务业务逻辑区，封装业务相关的核心代码
│   ├── api/           # 接口层，处理 HTTP 请求响应、路由定义
│   ├── infra/         # 基础组件与通用工具，支撑各业务模块复用
│   │   ├── conf/      
│   │   ├── errs/      # 自定义错误码
│   │   ├── global/    # 全局共享资源（如配置实例、日志对象）
│   │   └── app.go     # 应用上下文容器，管理全局实例的依赖注入
│   ├── domain/        # 领域层，定义业务实体与领域规则
│   │   ├── entity/    # 业务实体，映射核心业务对象
│   │   └── repo/      # 领域仓储接口，规范数据操作
│   ├── service/       # 应用服务层，编排业务流程、协调领域对象交互
│   └── server.go      # 
├── ext/               # 扩展工具包，可复用的第三方库封装或基础工具
│   ├── fnutil/          
│   ├── pagination/         
│   └── logext/          
├── .dockerignore      
├── .gitignore         
├── Dockerfile         
├── go.mod             
└── main.go          
```

## 二、新增业务功能流程

- **定义实体**：在 domain/entity 新增业务实体（如 BiddingCustomRule），补充属性与领域实体规则方法。
- **定义仓储接口**：扩展仓储，在 domain/repo 新增仓储接口（如 BiddingCustomRuleRepo），定义数据操作规范。
- **实现服务**：在 service 层编排业务流程
- **暴露接口**：在 api 层新增 HTTP 接口，解析请求、调用 service 方法，返回响应。