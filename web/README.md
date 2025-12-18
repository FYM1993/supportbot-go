# 客服客户端 (Customer Service Client)

这是一个完整的电商客服前端客户端，提供电商页面展示和智能客服聊天功能。

## 功能特性

### 🛒 电商页面功能
- **商品展示**: 响应式商品网格布局，支持商品卡片展示
- **分类导航**: 左侧商品分类导航，热门分类标签
- **搜索功能**: 商品搜索框，支持关键词搜索
- **购物车**: 购物车图标和数量显示
- **响应式设计**: 支持桌面端和移动端访问

### 💬 客服聊天功能
- **实时聊天**: WebSocket连接，支持实时消息收发
- **智能客服**: 与AI助手进行智能对话
- **文件上传**: 支持图片等附件上传
- **会话管理**: 自动重连、心跳检测、连接状态显示
- **消息历史**: 聊天记录保存和显示

### 🔐 用户登录系统
- **登录弹窗**: 美观的登录界面
- **状态管理**: 登录状态持久化，支持自动登录
- **用户验证**: 模拟用户验证逻辑
- **会话管理**: 登录状态检查和过期处理

## 技术架构

### 前端技术栈
- **HTML5**: 语义化标签，现代化结构
- **CSS3**: Flexbox/Grid布局，响应式设计，动画效果
- **JavaScript ES6+**: 模块化设计，面向对象编程
- **STOMP over SockJS**: 实时双向通信，与Spring WebSocket兼容

### 核心模块
1. **WebSocketManager**: WebSocket连接管理
2. **ChatManager**: 聊天功能管理
3. **LoginManager**: 用户登录管理

## 使用方法

### 1. 启动后端服务
确保im-demo后端服务已启动并监听在11005端口。

### 2. 启动前端服务
**重要：不要直接打开HTML文件，必须使用HTTP服务器！**

#### 方法一：使用启动脚本（推荐）
```bash
# Windows
start-server.bat

# Linux/Mac
chmod +x start-server.sh
./start-server.sh
```

#### 方法二：手动启动
```bash
cd customer-service-client

# Python 3
python -m http.server 8080

# 或者 Python 2
python -m SimpleHTTPServer 8080

# 或者 Node.js
npx http-server -p 8080
```

### 3. 访问页面
在浏览器中访问：`http://localhost:8080/index.html`

### 3. 用户登录
- 点击右上角"免费注册"按钮
- 输入任意用户名即可登录（密码可选）
- 系统会自动生成用户ID并建立WebSocket连接

### 4. 使用客服功能
- 登录后，点击右下角"联系客服"按钮
- 在聊天窗口中与AI助手对话
- 支持文本消息和图片附件

## 配置说明

### WebSocket连接配置
在 `js/websocket.js` 中修改连接配置：
```javascript
const port = '11005'; // im-demo服务端口
```

### 登录验证配置
后端没有登录验证，前端只验证用户名不为空：
```javascript
if (username && username.trim().length > 0) {
    resolve(true);
} else {
    resolve(false);
}
```

## 文件结构

```
customer-service-client/
├── index.html              # 主页面
├── css/
│   └── style.css          # 样式文件
├── js/
│   ├── websocket.js       # WebSocket管理
│   ├── chat.js            # 聊天功能
│   └── login.js           # 登录功能
└── README.md              # 说明文档
```

## 浏览器兼容性

- Chrome 60+
- Firefox 55+
- Safari 12+
- Edge 79+

## 开发说明

### 添加新功能
1. 在相应的JavaScript文件中添加新方法
2. 在HTML中添加对应的UI元素
3. 在CSS中添加样式定义

### 自定义样式
- 主色调: `#e74c3c` (红色)
- 辅助色: `#8b4513` (棕色)
- 字体: Microsoft YaHei, Arial, sans-serif

### 扩展功能建议
- 商品详情页面
- 购物车功能实现
- 订单管理系统
- 用户注册功能
- 支付集成
- 多语言支持

## 注意事项

1. **WebSocket连接**: 确保im-demo后端服务正常运行
2. **端口配置**: 检查WebSocket连接端口是否正确
3. **HTTPS环境**: 在生产环境中使用HTTPS时，WebSocket会自动使用WSS协议
4. **浏览器控制台**: 查看控制台日志了解连接状态和错误信息

## 故障排除

### 常见问题
1. **WebSocket连接失败**: 检查后端服务是否启动，端口是否正确
2. **登录失败**: 确认使用正确的测试账号
3. **聊天无响应**: 检查WebSocket连接状态和网络连接

### 调试方法
1. 打开浏览器开发者工具
2. 查看Console标签页的日志信息
3. 检查Network标签页的WebSocket连接状态

## 许可证

本项目仅供学习和演示使用。
