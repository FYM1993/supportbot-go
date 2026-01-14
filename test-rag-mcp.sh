#!/bin/bash

# RAG + Function Calling 功能测试脚本

echo "🚀 测试 RAG + Function Calling 功能"
echo "======================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. 测试 knowledge-rag 服务
echo -e "${BLUE}[1] 测试 knowledge-rag 服务${NC}"
echo "-----------------------------------"

echo -e "${YELLOW}1.1 检查服务状态${NC}"
curl -s http://localhost:11003/api/health | jq '.'
echo ""

echo -e "${YELLOW}1.2 查看知识库统计${NC}"
curl -s http://localhost:11003/api/knowledge/stats | jq '.'
echo ""

echo -e "${YELLOW}1.3 测试知识检索（退货政策）${NC}"
curl -s -X GET 'http://localhost:11003/api/knowledge/search?q=退货政策' | jq '.results[] | {score: .Score, content: .Document.Content[:50]}'
echo ""

echo -e "${YELLOW}1.4 测试知识检索（优惠券）${NC}"
curl -s -X GET 'http://localhost:11003/api/knowledge/search?q=优惠券怎么用' | jq '.results[] | {score: .Score, content: .Document.Content[:50]}'
echo ""

echo -e "${GREEN}✅ knowledge-rag 服务测试完成${NC}"
echo ""
echo ""

# 2. 测试 assistant 服务（Function Calling）
echo -e "${BLUE}[2] 测试 assistant 服务（Function Calling）${NC}"
echo "-----------------------------------"

echo -e "${YELLOW}2.1 检查服务状态${NC}"
curl -s http://localhost:11002/api/health | jq '.'
echo ""

echo -e "${YELLOW}2.2 查看可用工具${NC}"
curl -s http://localhost:11002/api/tools | jq '.tools[] | {name: .name, description: .description}'
echo ""

echo -e "${GREEN}✅ assistant 服务测试完成${NC}"
echo ""
echo ""

# 3. 集成测试（通过 question-classifier）
echo -e "${BLUE}[3] 集成测试${NC}"
echo "-----------------------------------"

echo -e "${YELLOW}3.1 测试知识库问答（RAG）${NC}"
echo "提问：退货需要什么条件？"
curl -s 'http://localhost:11001/api/classify?question=退货需要什么条件？&uid=1001' | jq '.'
echo ""
sleep 2

echo -e "${YELLOW}3.2 测试商品查询（Function Calling）${NC}"
echo "提问：查询商品30001的信息"
curl -s 'http://localhost:11001/api/classify?question=查询商品30001的信息&uid=1002' | jq '.'
echo ""
sleep 2

echo -e "${YELLOW}3.3 测试订单查询（Function Calling）${NC}"
echo "提问：查询订单20240101001"
curl -s 'http://localhost:11001/api/classify?question=查询订单20240101001&uid=1003' | jq '.'
echo ""

echo -e "${GREEN}✅ 集成测试完成${NC}"
echo ""
echo ""

# 4. 总结
echo "======================================"
echo -e "${GREEN}🎉 所有测试完成！${NC}"
echo ""
echo "查看详细日志："
echo "  • knowledge-rag: logs/knowledge-rag.log"
echo "  • assistant: logs/assistant.log"
echo "  • question-classifier: logs/question-classifier.log"
echo ""
echo "查看 WebSocket 实时消息（im-demo 日志）："
echo "  • im-demo: logs/im-demo.log"
echo ""
echo "完整文档："
echo "  • cat RAG_MCP_README.md"
echo ""

