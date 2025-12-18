.PHONY: help build run clean test deps

# 默认目标
help:
	@echo "SupportBot-Go Makefile"
	@echo ""
	@echo "可用命令："
	@echo "  make deps      - 下载依赖"
	@echo "  make build     - 编译所有服务"
	@echo "  make run       - 运行所有服务"
	@echo "  make clean     - 清理编译文件"
	@echo "  make test      - 运行测试"
	@echo ""

# 下载依赖
deps:
	@echo "📦 下载依赖..."
	go mod tidy
	go mod download

# 编译所有服务
build:
	@echo "🔨 编译服务..."
	@mkdir -p bin
	go build -o bin/im-demo cmd/im-demo/main.go
	go build -o bin/question-classifier cmd/question-classifier/main.go
	go build -o bin/assistant cmd/assistant/main.go
	go build -o bin/general-chat cmd/general-chat/main.go
	go build -o bin/knowledge-rag cmd/knowledge-rag/main.go
	@echo "✅ 编译完成！二进制文件在 bin/ 目录"

# 运行所有服务
run:
	@./start.sh

# 停止所有服务
stop:
	@./stop.sh

# 清理
clean:
	@echo "🧹 清理..."
	rm -rf bin/
	rm -rf logs/
	@echo "✅ 清理完成"

# 测试
test:
	@echo "🧪 运行测试..."
	go test -v ./...

# 格式化代码
fmt:
	@echo "💅 格式化代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "🔍 代码检查..."
	golangci-lint run

# 安装开发工具
install-tools:
	@echo "🛠️  安装开发工具..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 快速启动（编译 + 运行）
quick: build
	@./start.sh

