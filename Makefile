.PHONY: run test lint
run:
	go run ./cmd/gateway

test:
	## Go语言的测试工具 ./...是通配符
	go test ./...

lint:
	## golangci-lint 是 Go 社区最流行的静态代码分析工具./...是通配符
	golangci-lint run ./...