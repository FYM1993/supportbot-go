# SupportBot-Go - 企业级客服 Agent 系统（Golang 版）

> 基于 Golang + Gin + WebSocket + 通义千问的多 Agent 客服系统

## 📖 项目简介

这是一个完全用 Golang 重写的企业级 Agent 客服系统，原 Java 版本的架构保持一致，但使用 Go 的高性能并发特性重新实现。

### 核心特性

- ✅ **WebSocket 实时通信**：基于 gorilla/websocket 的长连接管理
- ✅ **多 Agent 协作**：问题分类 + 3 个专业 Agent
- ✅ **异步回调架构**：高性能非阻塞处理
- ✅ **LLM 集成**：通义千问（DashScope）API
- ✅ **会话管理**：双向索引 + 心跳检测
- ✅ **并发安全**：sync.RWMutex + Goroutine

## 🏗️ 系统架构

```
┌─────────────┐
│   前端页面   │ (customer-service-client)
└──────┬──────┘
       │ WebSocket
       ▼
┌─────────────────────────────────┐
│      im-demo (11005)            │ WebSocket 网关
│  - 连接管理 - 消息推送          │
└──────┬──────────────────────────┘
       │ HTTP API
       ▼
┌─────────────────────────────────┐
│  question-classifier (11001)    │ 问题分类路由
│  - 意图识别 - Agent 路由        │
└──────┬──────────────────────────┘
       │ 路由到不同 Agent
       ├──────┬──────┬──────┐
       ▼      ▼      ▼      ▼
  ┌────────┐ ┌────┐ ┌─────┐
  │assistant│ │chat│ │ rag │
  │ (11002) │ │(11003)│(11004)│
  └────────┘ └────┘ └─────┘
       │      │      │
       └──────┴──────┘
              │ 回调
              ▼
        回到 im-demo 推送
```

## 🚀 快速开始

### 环境要求

- Go 1.21+
- Redis 6.0+
- 通义千问 API Key

### 安装依赖

```bash
cd supportbot-go
go mod tidy
```

### 配置 API Key

编辑配置文件，替换 `YOUR_DASHSCOPE_API_KEY`：

```bash
# configs/question-classifier.yaml
# configs/assistant.yaml
# configs/general-chat.yaml
# configs/knowledge-rag.yaml
```

### 启动 Redis

```bash
# macOS (Homebrew)
brew services start redis

# 或者直接启动
redis-server
```

### 启动服务

**方式 1：分别启动（推荐，便于调试）**

```bash
# 终端 1：启动 im-demo
go run cmd/im-demo/main.go

# 终端 2：启动 question-classifier
go run cmd/question-classifier/main.go

# 终端 3：启动 assistant
go run cmd/assistant/main.go

# 终端 4：启动 general-chat
go run cmd/general-chat/main.go

# 终端 5：启动 knowledge-rag
go run cmd/knowledge-rag/main.go
```

**方式 2：编译后启动**

```bash
# 编译所有服务
go build -o bin/im-demo cmd/im-demo/main.go
go build -o bin/question-classifier cmd/question-classifier/main.go
go build -o bin/assistant cmd/assistant/main.go
go build -o bin/general-chat cmd/general-chat/main.go
go build -o bin/knowledge-rag cmd/knowledge-rag/main.go

# 启动
./bin/im-demo &
./bin/question-classifier &
./bin/assistant &
./bin/general-chat &
./bin/knowledge-rag &
```

### 访问前端

1. 打开原项目的前端：`customer-service-client/index.html`
2. 输入任意用户名登录
3. 开始对话！

## 📁 项目结构

```
supportbot-go/
├── cmd/                      # 各服务的入口程序
│   ├── im-demo/
│   ├── question-classifier/
│   ├── assistant/
│   ├── general-chat/
│   └── knowledge-rag/
├── internal/                 # 内部代码
│   ├── model/               # 数据模型
│   ├── service/             # 业务服务
│   ├── handler/             # HTTP/WebSocket 处理器
│   ├── middleware/          # 中间件
│   ├── client/              # 外部客户端（LLM）
│   └── config/              # 配置管理
├── pkg/                     # 可复用的包
│   ├── logger/              # 日志工具
│   └── redis/               # Redis 工具
├── configs/                 # 配置文件
│   ├── im-demo.yaml
│   ├── question-classifier.yaml
│   ├── assistant.yaml
│   ├── general-chat.yaml
│   └── knowledge-rag.yaml
├── go.mod
└── README.md
```

## 🎯 核心服务说明

### 1. im-demo（WebSocket 网关）

- **端口**：11005
- **职责**：
  - WebSocket 连接管理
  - 用户会话管理（双向索引）
  - 心跳检测（30s 检查，60s 超时）
  - 消息推送

**关键代码**：
```go
// 会话管理（并发安全）
type SessionService struct {
    userSessions  map[int64]*UserSession  // userId -> session
    sessionToUser map[string]int64        // sessionId -> userId
    mu            sync.RWMutex            // 读写锁
}

// 心跳检测（Goroutine）
go s.heartbeatChecker()
```

### 2. question-classifier（问题分类）

- **端口**：11001
- **职责**：
  - 调用通义千问进行问题分类
  - 路由到对应的 Agent
  - 对话历史管理（Redis）

**分类类型**：
- `product.inquiry`: 商品咨询 → assistant
- `order.status`: 订单查询 → assistant
- `knowledge.query`: 知识库查询 → knowledge-rag
- `general-chat`: 通用对话 → general-chat

### 3. assistant（业务助手）

- **端口**：11002
- **职责**：
  - 调用业务 API（商品、订单、工单）
  - 使用 LLM 生成友好回复

### 4. general-chat（通用对话）

- **端口**：11003
- **职责**：闲聊、问候等一般性对话

### 5. knowledge-rag（知识库）

- **端口**：11004
- **职责**：
  - 知识库检索（简化版，未集成 Milvus）
  - RAG 问答

## 🔧 性能优化建议

1. **连接池**：HTTP 客户端使用连接池
2. **Redis 连接复用**：使用 go-redis 的连接池
3. **Goroutine 池**：限制并发数量（如 `ants` 库）
4. **消息批量推送**：多个消息合并推送
5. **监控指标**：Prometheus + Grafana

## 📚 学习路径

1. **第一天**：理解架构，运行起来
2. **第二天**：深入 WebSocket 和会话管理
3. **第三天**：学习 LLM 集成和 Agent 路由
4. **第四天**：优化性能，准备面试

## 🤝 贡献

欢迎提 Issue 和 PR！

## 📄 License

MIT

---

**祝你面试顺利！🎉**

如有问题，请提 Issue 或联系作者。

