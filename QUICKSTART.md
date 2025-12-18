# 🚀 快速入门指南

## 1️⃣ 配置 API Key（必须！）

编辑以下文件，替换 `YOUR_DASHSCOPE_API_KEY` 为你的通义千问 API Key：

```bash
# 方式 1：手动编辑
vim configs/question-classifier.yaml
vim configs/assistant.yaml
vim configs/general-chat.yaml
vim configs/knowledge-rag.yaml

# 方式 2：批量替换（Mac/Linux）
find configs -name "*.yaml" -exec sed -i '' 's/YOUR_DASHSCOPE_API_KEY/sk-xxx/g' {} \;
```

## 2️⃣ 启动 Redis

```bash
# macOS (Homebrew)
brew services start redis

# 或者直接启动
redis-server

# 测试
redis-cli ping  # 应该返回 PONG
```

## 3️⃣ 启动服务

**方式 A：一键启动（推荐）**

```bash
./start.sh
```

**方式 B：分别启动（便于调试）**

```bash
# 终端 1
go run cmd/im-demo/main.go

# 终端 2
go run cmd/question-classifier/main.go

# 终端 3
go run cmd/assistant/main.go

# 终端 4
go run cmd/general-chat/main.go

# 终端 5
go run cmd/knowledge-rag/main.go
```

## 4️⃣ 测试服务

```bash
# 测试各服务健康状态
curl http://localhost:11005/api/health  # im-demo
curl http://localhost:11001/api/health  # question-classifier
curl http://localhost:11002/api/health  # assistant
curl http://localhost:11003/api/health  # general-chat
curl http://localhost:11004/api/health  # knowledge-rag
```

## 5️⃣ 打开前端

1. 复制 Java 版本的前端到当前目录（可选）：

```bash
cp -r ../supportbot/customer-service-client ./frontend
```

2. 打开 `frontend/index.html` 或直接用 Java 版本的前端

3. 输入任意用户名登录（如 `test`）

4. 开始对话！

## 6️⃣ 测试对话

试试这些问题：

- "查询商品信息" → 路由到 assistant
- "我的订单在哪里" → 路由到 assistant
- "如何使用这个产品" → 路由到 knowledge-rag
- "你好" → 路由到 general-chat

## 7️⃣ 停止服务

```bash
./stop.sh
```

或者手动杀死进程：

```bash
kill $(cat logs/*.pid)
```

## 🐛 常见问题

### 问题 1：端口被占用

```bash
# 查看端口占用
lsof -i :11001

# 杀死进程
kill -9 <PID>
```

### 问题 2：Redis 连接失败

```bash
# 检查 Redis 是否运行
redis-cli ping

# 如果没运行
brew services start redis
```

### 问题 3：API Key 未配置

如果看到 `401 Unauthorized` 错误，检查配置文件中的 API Key 是否正确。

### 问题 4：依赖下载失败

```bash
# 设置代理（如果需要）
export GOPROXY=https://goproxy.cn,direct

# 重新下载
go mod tidy
```

### 问题 5：前端连接不上

检查 WebSocket 连接：

```bash
# 浏览器控制台
ws://localhost:11005/ws?uid=123
```

如果连接失败，检查 im-demo 是否启动。

## 📊 查看日志

```bash
# 实时查看
tail -f logs/im-demo.log
tail -f logs/question-classifier.log

# 查看所有日志
cat logs/*.log
```

## 🎯 下一步

- 阅读 `README.md` 了解架构
- 修改分类规则（`cmd/question-classifier/main.go`）
- 添加新的 Agent
- 集成真实的业务 API

---

**祝你使用愉快！🎉**

