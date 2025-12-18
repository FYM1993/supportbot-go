#!/bin/bash

echo "启动客服客户端HTTP服务器..."
echo ""
echo "请确保已安装Python 3"
echo "如果没有安装Python，请先安装Python 3"
echo ""
echo "服务器将在 http://localhost:8080 启动"
echo "按 Ctrl+C 停止服务器"
echo ""

python3 -m http.server 8080
