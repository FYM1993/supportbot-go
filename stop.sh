#!/bin/bash

# SupportBot-Go 停止脚本

echo "🛑 停止 SupportBot-Go 服务..."
echo ""

# 读取 PID 并杀死进程
for service in im-demo question-classifier assistant general-chat knowledge-rag; do
    pid_file="logs/${service}.pid"
    if [ -f "$pid_file" ]; then
        pid=$(cat "$pid_file")
        if kill -0 "$pid" 2>/dev/null; then
            echo "停止 $service (PID: $pid)..."
            kill "$pid"
            rm "$pid_file"
        else
            echo "$service 已经停止"
            rm "$pid_file"
        fi
    else
        echo "$service 的 PID 文件不存在"
    fi
done

echo ""
echo "✅ 所有服务已停止"

