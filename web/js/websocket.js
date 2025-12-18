/**
 * 原生 WebSocket 连接管理
 * 负责与im-demo后端建立WebSocket连接，处理消息收发
 */
class WebSocketManager {
    constructor() {
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectInterval = 3000; // 3秒
        this.heartbeatInterval = null;
        this.isConnected = false;
        this.userId = null;
        
        // 绑定事件处理器
        this.onMessageCallback = null;
        this.onConnectCallback = null;
        this.onDisconnectCallback = null;
        
        // 用于处理请求-响应模式的回调
        this.pendingResponses = new Map(); // messageId -> {resolve, reject, timeout}
    }

    /**
     * 连接到原生 WebSocket 服务器
     * @param {string|number} userId 用户ID
     */
    connect(userId) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            console.log('WebSocket 已连接');
            return;
        }

        // 确保userId是数字类型
        this.userId = typeof userId === 'string' ? parseInt(userId, 10) : userId;
        
        if (isNaN(this.userId)) {
            console.error('无效的用户ID:', userId);
            return;
        }
        
        // 构建 WebSocket URL
        const wsUrl = window.getWebSocketUrl(this.userId);
        // 替换 http/https 为 ws/wss
        const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const finalUrl = wsUrl.replace(/^https?:/, wsProtocol);
        
        console.log('正在连接 WebSocket:', finalUrl);
        console.log('用户ID:', this.userId);

        try {
            // 创建原生 WebSocket 连接
            this.ws = new WebSocket(finalUrl);
            
            // 设置事件处理器
            this.setupEventHandlers();
        } catch (error) {
            console.error('WebSocket 连接失败:', error);
            this.scheduleReconnect();
        }
    }

    /**
     * 设置 WebSocket 事件处理器
     */
    setupEventHandlers() {
        this.ws.onopen = () => {
            console.log('=== WebSocket 连接成功 ===');
            console.log('用户ID:', this.userId);
            
            this.isConnected = true;
            this.reconnectAttempts = 0;
            
            // 开始心跳
            this.startHeartbeat();
            
            if (this.onConnectCallback) {
                this.onConnectCallback();
            }
            
            console.log('=== WebSocket 连接初始化完成 ===');
        };
        
        this.ws.onmessage = (event) => {
            console.log('=== 收到消息 ===');
            console.log('原始消息:', event.data);
            
            try {
                const data = JSON.parse(event.data);
                console.log('解析后的消息数据:', data);
                this.handleMessage(data);
            } catch (error) {
                console.error('解析消息失败:', error);
                console.error('原始消息体:', event.data);
            }
        };
        
        this.ws.onerror = (error) => {
            console.error('WebSocket 错误:', error);
        };
        
        this.ws.onclose = (event) => {
            console.log('WebSocket 连接断开:', event.code, event.reason);
            this.isConnected = false;
            this.stopHeartbeat();
            
            if (this.onDisconnectCallback) {
                this.onDisconnectCallback();
            }
            
            // 尝试重连
            if (event.code !== 1000) { // 1000 是正常关闭
                this.scheduleReconnect();
            }
        };
    }

    /**
     * 处理收到的消息
     * @param {Object} data 消息数据
     */
    handleMessage(data) {
        console.log('收到消息:', data);
        
        // 检查是否是响应消息
        if (data.success !== undefined && data.messageID) {
            console.log('识别为响应消息');
            this.handleResponse(data);
            return;
        }
        
        // 处理各种类型的消息
        if (this.onMessageCallback) {
            this.onMessageCallback(data);
        }
    }

    /**
     * 发送聊天消息
     * @param {string} content 消息内容
     */
    sendChatMessage(content) {
        if (this.isConnected && this.userId) {
            const messageId = this.generateMessageId();
            const message = {
                messageId: messageId,
                type: 'CHAT',
                content: content,
                sender: this.userId,
                timestamp: new Date().toISOString()  // 改为 ISO 8601 格式
            };
            
            console.log('发送的消息对象:', message);
            console.log('发送的消息ID:', messageId);
            
            // 返回Promise，等待服务器响应
            return this.sendWithResponse(message, messageId);
        } else {
            console.error('无法发送消息: WebSocket未连接或用户未登录');
            return Promise.reject(new Error('WebSocket未连接或用户未登录'));
        }
    }

    /**
     * 发送心跳消息
     */
    sendHeartbeat() {
        if (this.isConnected) {
            const message = {
                type: 'HEARTBEAT',
                userId: this.userId,
                timestamp: new Date().toISOString()  // 改为 ISO 8601 格式
            };
            this.send(message);
        }
    }

    /**
     * 发送消息到服务器
     * @param {Object} message 消息对象
     */
    send(message) {
        console.log('=== 开始发送消息到服务器 ===');
        console.log('消息内容:', message);
        console.log('WebSocket 状态:', {
            exists: !!this.ws,
            readyState: this.ws ? this.ws.readyState : 'N/A'
        });
        
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            try {
                this.ws.send(JSON.stringify(message));
                console.log('消息已发送');
                console.log('=== 消息发送完成 ===');
            } catch (error) {
                console.error('发送消息失败:', error);
            }
        } else {
            console.error('WebSocket 未连接，无法发送消息');
            console.error('连接状态:', this.ws ? this.ws.readyState : 'N/A');
        }
    }

    /**
     * 开始心跳
     */
    startHeartbeat() {
        this.heartbeatInterval = setInterval(() => {
            this.sendHeartbeat();
        }, 30000); // 30秒发送一次心跳
    }

    /**
     * 停止心跳
     */
    stopHeartbeat() {
        if (this.heartbeatInterval) {
            clearInterval(this.heartbeatInterval);
            this.heartbeatInterval = null;
        }
    }

    /**
     * 生成消息ID
     * @returns {string} 唯一的消息ID
     */
    generateMessageId() {
        return 'msg_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
    }

    /**
     * 发送消息并等待响应
     * @param {Object} message 消息对象
     * @param {string} messageId 消息ID
     * @returns {Promise<Object>} 返回Promise，resolve时包含响应结果
     */
    sendWithResponse(message, messageId) {
        return new Promise((resolve, reject) => {
            // 设置超时（30秒）
            const timeout = setTimeout(() => {
                this.pendingResponses.delete(messageId);
                reject(new Error('请求超时'));
            }, 30000);
            
            // 存储Promise的resolve和reject函数
            this.pendingResponses.set(messageId, {
                resolve: resolve,
                reject: reject,
                timeout: timeout
            });
            
            // 发送消息
            try {
                this.send(message);
                console.log('消息已发送，等待响应: messageId=', messageId);
            } catch (error) {
                this.pendingResponses.delete(messageId);
                clearTimeout(timeout);
                reject(error);
            }
        });
    }

    /**
     * 处理响应消息
     * @param {Object} response 响应数据
     */
    handleResponse(response) {
        const { messageID, success, message } = response;
        console.log('=== 收到服务器响应 ===');
        console.log('响应的消息ID:', messageID);
        console.log('响应状态:', success);
        console.log('响应消息:', message);
        
        // 查找对应的Promise
        const pendingResponse = this.pendingResponses.get(messageID);
        if (pendingResponse) {
            console.log('找到匹配的待响应消息，开始处理...');
            // 清除超时
            clearTimeout(pendingResponse.timeout);
            this.pendingResponses.delete(messageID);
            
            // 根据状态决定resolve还是reject
            if (success) {
                console.log('消息发送成功，调用resolve');
                pendingResponse.resolve({
                    success: true,
                    message: message,
                    messageId: messageID
                });
            } else {
                console.log('消息发送失败，调用reject');
                pendingResponse.reject({
                    success: false,
                    message: message,
                    messageId: messageID
                });
            }
        } else {
            console.warn('未找到匹配的待响应消息ID:', messageID);
        }
    }

    /**
     * 安排重连
     */
    scheduleReconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            console.log(`尝试重连 (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);
            
            setTimeout(() => {
                if (this.userId) {
                    this.connect(this.userId);
                }
            }, this.reconnectInterval);
        } else {
            console.error('达到最大重连次数，停止重连');
        }
    }

    /**
     * 断开连接
     */
    disconnect() {
        this.stopHeartbeat();
        
        if (this.ws) {
            this.ws.close(1000, '用户主动断开');
            this.ws = null;
        }
        
        this.isConnected = false;
        this.userId = null;
    }

    /**
     * 设置消息回调
     * @param {Function} callback 回调函数
     */
    onMessage(callback) {
        this.onMessageCallback = callback;
    }

    /**
     * 设置连接回调
     * @param {Function} callback 回调函数
     */
    onConnect(callback) {
        this.onConnectCallback = callback;
    }

    /**
     * 设置断开连接回调
     * @param {Function} callback 回调函数
     */
    onDisconnect(callback) {
        this.onDisconnectCallback = callback;
    }

    /**
     * 获取连接状态
     * @returns {boolean} 是否已连接
     */
    getConnectionStatus() {
        return this.ws && this.ws.readyState === WebSocket.OPEN;
    }
}

// 创建全局WebSocket管理器实例
window.wsManager = new WebSocketManager();
