# 客服客户端配置说明

## 🚀 快速配置

### 1. 修改服务地址
编辑 `js/config.js` 文件：

```javascript
const IM_SERVICE_CONFIG = {
    // 修改为你的im-demo服务实际IP地址
    host: '192.168.1.100', // 替换为实际IP
    
    // 服务端口（通常不需要修改）
    port: '11005',
    
    // WebSocket端点（通常不需要修改）
    wsEndpoint: '/ws'
};
```

### 2. 常见IP地址示例

#### 本地开发
```javascript
host: 'localhost'  // 或 '127.0.0.1'
```

#### 局域网访问
```javascript
host: '192.168.1.100'  // 你的实际局域网IP
host: '10.0.0.50'      // 另一个局域网IP示例
```

#### 公网访问
```javascript
host: 'your-domain.com'  // 你的域名
host: '203.0.113.1'      // 公网IP地址
```

## 🔧 详细配置说明

### 配置文件位置
- **主配置文件**: `js/config.js`
- **影响范围**: 所有页面都会使用此配置

### 配置项说明

| 配置项 | 说明 | 默认值 | 是否必填 |
|--------|------|--------|----------|
| `host` | im-demo服务IP地址 | `localhost` | ✅ 是 |
| `port` | im-demo服务端口 | `11005` | ❌ 否 |
| `wsEndpoint` | WebSocket端点 | `/ws` | ❌ 否 |

### 修改后生效
1. 保存配置文件
2. 刷新浏览器页面
3. 重新连接WebSocket

## 🌐 网络配置检查

### 1. 检查im-demo服务状态
```bash
# 检查服务是否运行
curl http://YOUR_IP:11005/health

# 或者浏览器访问
http://YOUR_IP:11005/health
```

### 2. 检查网络连通性
```bash
# Windows
ping YOUR_IP

# Linux/Mac
ping -c 4 YOUR_IP
```

### 3. 检查防火墙设置
确保端口11005在防火墙中开放

## 🚨 常见问题

### 1. 连接失败
- 检查IP地址是否正确
- 确认im-demo服务已启动
- 检查网络连通性

### 2. 跨域问题
- 确保使用HTTP服务器访问，不要直接打开HTML文件
- 检查im-demo的CORS配置

### 3. 端口被占用
- 修改配置文件中的端口号
- 同时修改im-demo服务的端口配置

## 📝 配置示例

### 开发环境
```javascript
host: 'localhost'
port: '11005'
```

### 测试环境
```javascript
host: '192.168.1.100'
port: '11005'
```

### 生产环境
```javascript
host: 'your-domain.com'
port: '11005'
```

## 🔍 调试方法

1. **打开浏览器开发者工具** (F12)
2. **查看Console标签页**
3. **检查WebSocket连接状态**
4. **查看网络请求**

## 📞 技术支持

如果遇到配置问题，请检查：
1. 配置文件语法是否正确
2. 服务地址是否可访问
3. 网络连接是否正常
4. 浏览器控制台是否有错误信息
