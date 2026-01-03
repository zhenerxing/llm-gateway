# LLM Gateway
A gateway service for Large Language Models.

目录结构
llm-gateway/
  cmd/
    gateway/
      main.go
  internal/
    handlers/
      healthz.go
      version.go
    middleware/
      requestid.go
      logging.go
    server/
      server.go
      router.go
  pkg/
    logger/
      logger.go
  .github/
    workflows/
      ci.yml
  go.mod
  go.sum
  Makefile
  README.md