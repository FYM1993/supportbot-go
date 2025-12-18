@echo off
chcp 65001 >nul
echo 启动客服客户端HTTP服务器...
echo.
echo 请确保已安装Python 3
echo 如果没有安装Python，请先安装Python 3
echo.
echo 正在检查Python安装...
python --version >nul 2>&1
if errorlevel 1 (
    echo 错误：未找到Python，请先安装Python 3
    echo 下载地址：https://www.python.org/downloads/
    pause
    exit /b 1
)
echo Python检查通过
echo.
echo 服务器将在 http://localhost:8081 启动
echo 按 Ctrl+C 停止服务器
echo.

cd /d "%~dp0"
echo 正在启动服务器...
python -m http.server 8081

echo.
echo 服务器已停止
pause
