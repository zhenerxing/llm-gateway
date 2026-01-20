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
├── Makefile
├── README.md
├── cmd
│   └── gateway
│       └── main.go # 程序入口：初始化配置/日志并启动 HTTP 服务
├── data
│   └── audit.db
├── go.mod
├── go.sum
├── internal
│   ├── apperr
│   │   ├── apperr.go
│   │   ├── codes.go
│   │   └── httpmap.go
│   ├── audit
│   │   ├── model.go
│   │   └── store_sqlite.go
│   ├── auth
│   │   ├── error.go
│   │   ├── inmemory_store.go
│   │   ├── keystore.go
│   │   ├── middleware.go
│   │   └── service.go
│   └── http
│       ├── handler
│       │   ├── admin_key.go
│       │   ├── audit.go
│       │   ├── chat.go
│       │   ├── healthz.go
│       │   ├── healthz_test.go
│       │   ├── version.go
│       │   └── version_test.go
│       ├── middleware
│       │   ├── admin_auth.go
│       │   ├── audit.go
│       │   ├── error_handler.go
│       │   ├── logging.go
│       │   ├── requestid.go
│       │   └── requestid_test.go
│       └── server
│           ├── router.go
│           └── server.go
└── pkg
    ├── logger
    │   └── logger.go
    └── response
        └── response.go

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


## 错误码（Error Codes）

本服务所有失败响应均使用统一结构：

```json
{
  "error": {
    "code": "PLATFORM_REQUEST_INVALID",
    "message": "tenant_id is required",
    "type": "platform",
    "retry_after": 0,
    "details": { "tenant_id": "required" },
    "request_id": "7f6c8a0056d1c86dc0dfa7583e3ccb1f"
  }
}
```

字段说明：

- `code`：稳定的错误码（用于程序判断与告警聚合）
- `message`：面向调用方的简要说明（不保证稳定，勿做强依赖）
- `type`：
  - `platform`：平台自身错误（参数、鉴权、内部错误等）
  - `upstream`：上游/依赖服务错误（模型服务、第三方 API 等）
- `retry_after`：建议重试等待秒数；`0` 表示未设置
- `details`：可选的结构化信息（字段校验错误、上下文信息等）
- `request_id`：请求追踪 ID（用于排查日志）

---

## 平台侧错误（Type = `platform`）

| code | HTTP | 含义 | 客户端建议 |
|---|---:|---|---|
| `PLATFORM_REQUEST_INVALID` | 400 | 请求参数/JSON 不合法、缺字段、格式错误 | 修正请求后重试 |
| `PLATFORM_AUTH_MISSING_API_KEY` | 401 | 缺少 API Key | 补充 API Key 后重试 |
| `PLATFORM_AUTH_INVALID_API_KEY` | 401 | API Key 无效/不存在/已被禁用 | 更换有效 Key |
| `PLATFORM_AUTH_FORBIDDEN` | 403 | 权限不足（例如非管理员访问管理接口） | 换有权限的凭证/联系管理员 |
| `PLATFORM_CONFLICT` | 409 | 资源冲突（例如重复创建） | 修改请求或先查询现有资源 |
| `PLATFORM_DEPENDENCY_UNAVAILABLE` | 503 | 平台依赖不可用（DB、缓存等） | 可稍后重试；如有 `retry_after` 按其等待 |
| `PLATFORM_INTERNAL_ERROR` | 500 | 平台内部错误 | 可重试；持续发生请带 `request_id` 反馈 |

> 说明：`PLATFORM_CONFLICT` 需要在代码常量中定义（如果还未定义请补上）。

---

## 上游侧错误（Type = `upstream`）

| code | HTTP | 含义 | 客户端建议 |
|---|---:|---|---|
| `UPSTREAM_TIMEOUT` | 504 | 上游请求超时 | 可重试；必要时降低超时/换模型 |
| `UPSTREAM_UNAVAILABLE` | 502 | 上游不可用或返回异常 | 可重试；持续发生请切换上游 |

---

## 重试策略建议（通用）

- 遇到 `503/502/504`：建议指数退避重试（如 1s, 2s, 4s…），若响应包含 `retry_after` 优先遵循。
- 遇到 `400/401/403/409`：通常不建议盲目重试，应先修正请求/凭证/权限/冲突条件。

---

## 错误码兼容性承诺

- `code`：稳定字段，新增不破坏既有语义。
- `message/details`：用于人类排查，可能随版本调整；客户端不应依赖具体文案或 details 的 key。

---
