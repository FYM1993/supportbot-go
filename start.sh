#!/bin/bash

# SupportBot-Go 一键启动脚本

echo "🚀 启动 SupportBot-Go 服务..."
echo ""

# 检查 Redis 是否运行
if ! redis-cli ping > /dev/null 2>&1; then
    echo "❌ Redis 未运行，正在启动..."
    if command -v brew > /dev/null; then
        brew services start redis
    else
        redis-server &
    fi
    sleep 2
fi

echo "✅ Redis 已就绪"
echo ""

# 创建日志目录
mkdir -p logs

# 启动服务
echo "启动 im-demo (端口 11005)..."
go run cmd/im-demo/main.go > logs/im-demo.log 2>&1 &
IM_DEMO_PID=$!

sleep 1

echo "启动 question-classifier (端口 11001)..."
go run cmd/question-classifier/main.go > logs/question-classifier.log 2>&1 &
CLASSIFIER_PID=$!

sleep 1

echo "启动 assistant (端口 11002)..."
go run cmd/assistant/main.go > logs/assistant.log 2>&1 &
ASSISTANT_PID=$!

sleep 1

echo "启动 general-chat (端口 11003)..."
go run cmd/general-chat/main.go > logs/general-chat.log 2>&1 &
CHAT_PID=$!

sleep 1

echo "启动 knowledge-rag (端口 11004)..."
go run cmd/knowledge-rag/main.go > logs/knowledge-rag.log 2>&1 &
RAG_PID=$!

echo ""
echo "✅ 所有服务已启动！"
echo ""
echo "进程 ID："
echo "  im-demo:              $IM_DEMO_PID"
echo "  question-classifier:  $CLASSIFIER_PID"
echo "  assistant:            $ASSISTANT_PID"
echo "  general-chat:         $CHAT_PID"
echo "  knowledge-rag:        $RAG_PID"
echo ""
echo "日志文件在 logs/ 目录下"
echo ""
echo "停止服务："
echo "  kill $IM_DEMO_PID $CLASSIFIER_PID $ASSISTANT_PID $CHAT_PID $RAG_PID"
echo ""
echo "或者运行: ./stop.sh"
echo ""
echo "查看日志："
echo "  tail -f logs/im-demo.log"
echo ""

# 保存 PID 到文件
echo "$IM_DEMO_PID" > logs/im-demo.pid
echo "$CLASSIFIER_PID" > logs/question-classifier.pid
echo "$ASSISTANT_PID" > logs/assistant.pid
echo "$CHAT_PID" > logs/general-chat.pid
echo "$RAG_PID" > logs/knowledge-rag.pid

echo "🎉 启动完成！打开前端页面开始使用吧！"

