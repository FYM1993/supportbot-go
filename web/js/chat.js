/**
 * 客服聊天功能管理
 * 处理聊天界面的显示、消息发送和接收
 */
class ChatManager {
    constructor() {
        this.isChatOpen = false;
        this.currentUserId = null;
        
        // DOM元素
        this.chatWindow = document.getElementById('chatWindow');
        this.chatMessages = document.getElementById('chatMessages');
        this.messageInput = document.getElementById('messageInput');
        this.sendBtn = document.getElementById('sendBtn');
        this.attachBtn = document.getElementById('attachBtn');
        this.customerServiceBtn = document.getElementById('customerServiceBtn');
        this.closeChatWindow = document.getElementById('closeChatWindow');
        
        this.initEventListeners();
        this.setupWebSocketHandlers();
    }

    /**
     * 初始化事件监听器
     */
    initEventListeners() {
        // 客服按钮点击事件
        this.customerServiceBtn.addEventListener('click', () => {
            this.toggleChat();
        });

        // 关闭聊天窗口
        this.closeChatWindow.addEventListener('click', () => {
            this.closeChat();
        });

        // 发送按钮点击事件
        this.sendBtn.addEventListener('click', () => {
            this.sendMessage();
        });

        // 回车键发送消息
        this.messageInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.sendMessage();
            }
        });

        // 附件按钮点击事件
        this.attachBtn.addEventListener('click', () => {
            this.handleAttachment();
        });

        // 点击聊天窗口外部关闭
        document.addEventListener('click', (e) => {
            if (this.isChatOpen && 
                !this.chatWindow.contains(e.target) && 
                !this.customerServiceBtn.contains(e.target)) {
                this.closeChat();
            }
        });
    }

    /**
     * 设置WebSocket消息处理器
     */
    setupWebSocketHandlers() {
        // 设置消息接收回调
        window.wsManager.onMessage((data) => {
            this.handleIncomingMessage(data);
        });

        // 设置连接状态回调
        window.wsManager.onConnect(() => {
            this.updateConnectionStatus(true);
        });

        window.wsManager.onDisconnect(() => {
            this.updateConnectionStatus(false);
        });
    }

    /**
     * 切换聊天窗口显示状态
     */
    toggleChat() {
        if (this.isChatOpen) {
            this.closeChat();
        } else {
            this.openChat();
        }
    }

    /**
     * 打开聊天窗口
     */
    openChat() {
        if (!this.currentUserId) {
            alert('请先登录后再使用客服功能');
            return;
        }

        this.isChatOpen = true;
        this.chatWindow.style.display = 'flex';
        this.customerServiceBtn.style.display = 'none';
        
        // 清空聊天记录，重新开始
        this.clearChat();
        
        // 确保WebSocket连接已建立
        if (window.wsManager && !window.wsManager.getConnectionStatus()) {
            window.wsManager.connect(this.currentUserId);
        }
        
        // 聚焦到输入框
        setTimeout(() => {
            this.messageInput.focus();
        }, 100);

        // 滚动到底部
        this.scrollToBottom();
    }

    /**
     * 关闭聊天窗口
     */
    closeChat() {
        this.isChatOpen = false;
        this.chatWindow.style.display = 'none';
        this.customerServiceBtn.style.display = 'flex';
        
        // 清空输入框
        this.messageInput.value = '';
    }

    /**
     * 发送消息
     */
    sendMessage() {
        const content = this.messageInput.value.trim();
        if (!content) {
            return;
        }

        console.log('=== 聊天管理器开始发送消息 ===');
        console.log('消息内容:', content);
        console.log('当前用户ID:', this.currentUserId);
        console.log('WebSocket连接状态:', window.wsManager.getConnectionStatus());

        if (!this.currentUserId) {
            console.error('用户未登录，无法发送消息');
            alert('请先登录');
            return;
        }

        if (!window.wsManager.getConnectionStatus()) {
            console.error('WebSocket连接已断开，无法发送消息');
            alert('连接已断开，请刷新页面重试');
            return;
        }

        // 先添加加载状态的消息到界面
        const tempMessageId = 'temp_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
        this.addMessage(content, 'user', false, true, tempMessageId);
        console.log('消息已添加到界面（加载状态）');

        // 发送消息到服务器
        console.log('开始调用 WebSocket 发送消息...');
        window.wsManager.sendChatMessage(content)
            .then((response) => {
                console.log('=== 聊天管理器收到处理结果响应 ===');
                console.log('响应的消息ID:', response.messageId);
                console.log('响应状态码:', response.code);
                console.log('响应消息:', response.message);
                
                // 根据处理结果更新消息状态
                if (response.code === 0) {
                    // 处理成功，移除加载状态
                    this.updateMessageStatus(tempMessageId, 'success');
                } else {
                    // 处理失败，显示重发按钮
                    this.updateMessageStatus(tempMessageId, 'failed');
                }
            })
            .catch((error) => {
                console.error('=== 聊天管理器收到错误响应 ===');
                console.error('错误的消息ID:', error.messageId);
                console.error('错误状态码:', error.code);
                console.error('错误消息:', error.message);
                
                // 处理失败，显示重发按钮
                this.updateMessageStatus(tempMessageId, 'failed');
            });

        // 清空输入框
        this.messageInput.value = '';

        // 滚动到底部
        this.scrollToBottom();
    }

    /**
     * 处理收到的消息
     * @param {Object} data 消息数据
     */
    handleIncomingMessage(data) {
        console.log('=== 聊天管理器处理收到的消息 ===');
        console.log('消息类型:', data.type);
        console.log('消息内容:', data.content);
        console.log('发送者:', data.sender);
        console.log('完整消息数据:', data);
        
        if (data.type === 'CHAT') {
            // 如果是AI回复，显示在界面上
            // 检查sender是否为AI（0或'0'）或系统消息
            if (data.sender === 0 || data.sender === '0' || data.sender === 'ai' || data.sender === 'system') {
                console.log('显示AI消息:', data.content);
                this.addMessage(data.content, 'ai');
            } else {
                console.log('CHAT类型消息，但sender不是AI，不显示:', data.sender);
            }
        } else if (data.type === 'SYSTEM') {
            // 系统消息
            console.log('显示系统消息:', data.content);
            this.addMessage(data.content, 'system');
        } else if (data.type === 'AI_RESPONSE') {
            // AI回复消息
            console.log('显示AI回复消息:', data.content);
            this.addMessage(data.content, 'ai');
        } else {
            console.log('未知消息类型，不处理:', data.type);
        }
        
        console.log('=== 消息处理完成 ===');
    }

    /**
     * 添加消息到聊天界面
     * @param {string} content 消息内容
     * @param {string} type 消息类型 (user/ai/system)
     * @param {boolean} isFailed 是否为失败状态
     * @param {boolean} isLoading 是否为加载中的消息
     * @param {string} messageId 消息ID，用于更新状态
     */
    addMessage(content, type, isFailed = false, isLoading = false, messageId = null) {
        const messageDiv = document.createElement('div');
        messageDiv.className = `message ${type}-message`;
        
        // 如果没有messageId，生成一个临时的
        if (!messageId) {
            messageId = 'temp_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
        }
        
        // 存储messageId到DOM元素中
        messageDiv.setAttribute('data-message-id', messageId);
        
        const messageContent = document.createElement('div');
        messageContent.className = 'message-content';
        messageContent.textContent = content;
        
        // 如果是加载中的用户消息，添加转圈特效
        if (isLoading && type === 'user') {
            const loadingSpinner = document.createElement('div');
            loadingSpinner.className = 'loading-spinner';
            loadingSpinner.innerHTML = '⏳';
            
            const statusContainer = document.createElement('div');
            statusContainer.className = 'message-status';
            statusContainer.appendChild(loadingSpinner);
            
            messageDiv.appendChild(statusContainer);
        }
        
        // 如果是失败的用户消息，添加重发按钮
        if (isFailed && type === 'user') {
            const retryButton = document.createElement('button');
            retryButton.className = 'retry-button';
            retryButton.textContent = '重发';
            retryButton.onclick = () => this.retryMessage(content);
            
            const buttonContainer = document.createElement('div');
            buttonContainer.className = 'retry-button-container';
            buttonContainer.appendChild(retryButton);
            
            messageDiv.appendChild(buttonContainer);
        }
        
        messageDiv.appendChild(messageContent);
        this.chatMessages.appendChild(messageDiv);
        
        // 滚动到底部
        this.scrollToBottom();
        
        return messageDiv;
    }

    /**
     * 更新消息状态
     * @param {string} messageId 消息ID
     * @param {string} status 新状态 ('loading', 'success', 'failed')
     */
    updateMessageStatus(messageId, status) {
        console.log('更新消息状态:', messageId, '状态:', status);
        
        // 首先尝试用传入的messageId查找
        let messageElement = this.chatMessages.querySelector(`[data-message-id="${messageId}"]`);
        
        // 如果没找到，尝试查找最后一个用户消息（通常是刚发送的）
        if (!messageElement) {
            console.log('未找到指定ID的消息元素，尝试查找最后一个用户消息:', messageId);
            const userMessages = this.chatMessages.querySelectorAll('.user-message');
            if (userMessages.length > 0) {
                messageElement = userMessages[userMessages.length - 1];
                console.log('找到最后一个用户消息，使用它来更新状态');
            }
        }
        
        if (!messageElement) {
            console.warn('未找到消息元素:', messageId);
            return;
        }
        
        // 移除现有的状态元素
        const existingStatus = messageElement.querySelector('.message-status');
        if (existingStatus) {
            existingStatus.remove();
        }
        
        // 根据新状态添加相应的元素
        if (status === 'loading') {
            const loadingSpinner = document.createElement('div');
            loadingSpinner.className = 'loading-spinner';
            loadingSpinner.innerHTML = '⏳';
            
            const statusContainer = document.createElement('div');
            statusContainer.className = 'message-status';
            statusContainer.appendChild(loadingSpinner);
            
            messageElement.insertBefore(statusContainer, messageElement.querySelector('.message-content'));
        } else if (status === 'failed') {
            const retryButton = document.createElement('button');
            retryButton.className = 'retry-button';
            retryButton.textContent = '重发';
            retryButton.onclick = () => {
                const content = messageElement.querySelector('.message-content').textContent;
                this.retryMessage(content);
            };
            
            const buttonContainer = document.createElement('div');
            buttonContainer.className = 'retry-button-container';
            buttonContainer.appendChild(retryButton);
            
            messageElement.appendChild(buttonContainer);
        }
        // success状态不需要添加任何元素，保持原样即可
    }

    /**
     * 重发消息
     * @param {string} content 消息内容
     */
    retryMessage(content) {
        console.log('重发消息:', content);
        
        // 清空输入框并设置内容
        this.messageInput.value = content;
        
        // 发送消息
        this.sendMessage();
    }

    /**
     * 处理附件上传
     */
    handleAttachment() {
        // 创建文件输入元素
        const fileInput = document.createElement('input');
        fileInput.type = 'file';
        fileInput.accept = 'image/*';
        
        fileInput.onchange = (e) => {
            const file = e.target.files[0];
            if (file) {
                this.uploadFile(file);
            }
        };
        
        fileInput.click();
    }

    /**
     * 上传文件
     * @param {File} file 文件对象
     */
    uploadFile(file) {
        // 这里可以实现文件上传逻辑
        // 目前只是显示文件名
        const fileName = file.name;
        this.addMessage(`已选择文件: ${fileName}`, 'user');
        
        // 可以在这里调用文件上传API
        console.log('文件上传:', fileName);
    }

    /**
     * 滚动到底部
     */
    scrollToBottom() {
        this.chatMessages.scrollTop = this.chatMessages.scrollHeight;
    }

    /**
     * 设置当前用户ID
     * @param {string|number} userId 用户ID
     */
    setCurrentUser(userId) {
        console.log('=== 聊天管理器设置当前用户 ===');
        console.log('之前的用户ID:', this.currentUserId);
        console.log('新的用户ID:', userId);
        
        // 确保userId是数字类型，如果传入的是字符串则转换为数字
        this.currentUserId = typeof userId === 'string' ? parseInt(userId, 10) : userId;
        
        // 验证userId是否为有效数字
        if (isNaN(this.currentUserId)) {
            console.error('无效的用户ID:', userId);
            this.currentUserId = null;
            return;
        }
        
        console.log('用户ID已更新 (数字类型):', this.currentUserId);
        
        // 如果用户已登录且聊天窗口已打开，检查WebSocket连接状态
        if (this.currentUserId && this.isChatOpen) {
            console.log('用户已登录且聊天窗口已打开，检查WebSocket连接...');
            if (!window.wsManager.getConnectionStatus()) {
                console.log('WebSocket未连接，开始连接...');
                window.wsManager.connect(this.currentUserId);
            } else {
                console.log('WebSocket已连接，无需重新连接');
            }
        } else {
            console.log('用户未登录或聊天窗口未打开:', {
                userId: this.currentUserId,
                isChatOpen: this.isChatOpen
            });
        }
        
        console.log('=== 聊天管理器用户设置完成 ===');
    }

    /**
     * 更新连接状态显示
     * @param {boolean} isConnected 是否已连接
     */
    updateConnectionStatus(isConnected) {
        const statusElement = document.querySelector('.chat-header span');
        if (statusElement) {
            if (isConnected) {
                statusElement.textContent = '在线客服';
                statusElement.style.color = '#4CAF50';
            } else {
                statusElement.textContent = '连接断开';
                statusElement.style.color = '#f44336';
            }
        }
    }

    /**
     * 清空聊天记录
     */
    clearChat() {
        this.chatMessages.innerHTML = '';
        
        // 添加欢迎消息
        const welcomeMessage = document.createElement('div');
        welcomeMessage.className = 'message ai-message';
        welcomeMessage.innerHTML = `
            <div class="message-content">
                您好! 我是智能客服助手, 有什么可以帮您的吗? 😊
            </div>
        `;
        this.chatMessages.appendChild(welcomeMessage);
    }

    /**
     * 获取聊天状态
     * @returns {boolean} 聊天窗口是否打开
     */
    getChatStatus() {
        return this.isChatOpen;
    }

    /**
     * 获取当前用户ID
     * @returns {string|null} 当前用户ID
     */
    getCurrentUser() {
        return this.currentUserId;
    }
}

// 创建全局聊天管理器实例
window.chatManager = new ChatManager();
