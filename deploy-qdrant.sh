#!/bin/bash

# Qdrant 向量数据库部署脚本

echo "🚀 部署 Qdrant 向量数据库..."

# 创建数据目录
mkdir -p ./qdrant_storage

# 启动 Qdrant 容器
docker run -d \
  --name qdrant \
  -p 6333:6333 \
  -p 6334:6334 \
  -v $(pwd)/qdrant_storage:/qdrant/storage \
  qdrant/qdrant:latest

echo "✅ Qdrant 已启动！"
echo "📍 HTTP API: http://localhost:6333"
echo "📍 Web UI: http://localhost:6333/dashboard"
echo ""
echo "测试连接："
sleep 3
curl http://localhost:6333/collections

