# 项目说明文档

## 一、项目结构-MVC

```
go-tpl/
├── cmd/               # 程序启动命令
│   └── run.go         
├── config/            # 配置文件
│   └── config.yml     
├── server/            # web业务
│   ├── api/           # api接口层，处理 HTTP 请求响应、路由定义 
│   ├── infra/         # 基础组件与通用工具，支撑各业务模块复用
│   │   ├── conf/      
│   │   ├── errs/      # 自定义错误码
│   │   ├── global/    # 全局共享资源（如配置实例、日志对象）
│   │   └── app.go     # 应用上下文容器，管理全局实例的依赖注入
│   ├── domain/        # 模型层，定义业务实体
│   │   ├── entity/    # 业务实体，映射核心业务对象
│   │   └── repo/      # 实体仓储接口，规范数据操作
│   ├── service/       # 服务层，编排业务流程
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

- **定义实体**：在 domain/entity 新增业务实体，映射数据库对象
- **定义仓储接口**：在 domain/repo 新增仓储接口，定义数据操作规范
- **实现服务**：在 service 层编排业务流程
- **暴露接口**：在 api 层新增 HTTP 接口，解析请求、调用 service 方法，返回响应