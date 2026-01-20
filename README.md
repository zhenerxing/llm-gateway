# LLM Gateway
A gateway service for Large Language Models.

## 功能特性

## 环境要求

- **Go**：Go **1.23+**（推荐使用 `toolchain` 指定的 **Go 1.24.11**）
  - 本项目 `go.mod`：
    - `go 1.23.0`
    - `toolchain go1.24.11`
- **网络**：需要能访问 Go Module 镜像/源码仓库以下载依赖（如 `proxy.golang.org` 或你配置的 `GOPROXY`）
- **操作系统**：macOS / Linux / Windows 均可（只要安装对应版本 Go）
- **依赖服务**：无（当前依赖仅为 Gin、Zap 等 Go 包；不强制要求数据库/Redis）

> 说明：若你的 Go 版本低于 1.23，可能无法通过构建；使用 Go 1.24.11 可以保证与 `toolchain` 配置一致。

## 如何运行

### 1、获取代码

```bash
git clone https://github.com/zhenerxing/llm-gateway.git
cd llm-gateway
```

### 2、命令行
```bash
# 下载依赖
go mod download

# 本地运行
make run

# 测试（等价于：go test ./...）
make test

# 静态检查（需要预先安装 golangci-lint）
make lint
```

## 目录结构

```text
llm-gateway/
├─ cmd/
│  └─ gateway/
│     └─ main.go              # 程序入口：初始化配置/日志并启动 HTTP 服务
├─ internal/
│  ├─ handlers/
│  │  ├─ healthz.go           # 健康检查接口（health check）
│  │  └─ version.go           # 版本信息接口（build/version 信息）
│  ├─ middleware/
│  │  ├─ requestid.go         # RequestID 中间件：为每个请求注入/传递请求 ID
│  │  └─ logging.go           # 访问日志中间件：请求/响应日志记录
│  └─ server/
│     ├─ server.go            # HTTP Server 生命周期管理（启动/关闭等）
│     └─ router.go            # 路由注册与 Gin 引擎初始化
├─ pkg/
│  └─ logger/
│     └─ logger.go            # 日志封装：基于 zap 的统一日志初始化与使用
├─ .github/
│  └─ workflows/
│     └─ ci.yml               # GitHub Actions CI：构建/测试/检查流程
├─ go.mod                     # Go module 定义与依赖版本
├─ go.sum                     # 依赖校验和（由 Go 自动维护）
├─ Makefile                   # 常用命令封装（build/test/lint/run 等）
└─ README.md                  # 项目文档

创建key：
curl -sS -X POST 'http://localhost:8080/admin/keys' \
  -H 'Content-Type: application/json' \
  -H 'X-Admin-Token: dev-admin-token' \
  -d '{"tenant_id":"acme","quota_daily_requests":1000,"quota_daily_tokens":200000}'

列出key：
curl -sS 'http://localhost:8080/admin/keys' \
  -H 'X-Admin-Token: dev-admin-token'
export API_KEY='...'

调用：
curl -sS -X POST 'http://localhost:8080/chat' \
  -H 'Content-Type: application/json' \
  -H "X-API-Key: ${API_KEY}" \
  -d '{"message":"hi"}'