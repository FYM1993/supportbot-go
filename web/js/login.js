/**
 * 用户登录功能管理
 * 处理登录弹窗的显示、表单提交和用户状态管理
 */
class LoginManager {
    constructor() {
        this.isLoggedIn = false;
        this.currentUser = null;
        
        // DOM元素
        this.loginModal = document.getElementById('loginModal');
        this.loginBtn = document.getElementById('loginBtn');
        this.closeLoginModal = document.getElementById('closeLoginModal');
        this.loginForm = document.getElementById('loginForm');
        this.loginStatus = document.getElementById('loginStatus');
        
        this.initEventListeners();
        this.checkLoginStatus();
    }

    /**
     * 初始化事件监听器
     */
    initEventListeners() {
        // 登录按钮点击事件
        this.loginBtn.addEventListener('click', () => {
            this.showLoginModal();
        });

        // 关闭登录弹窗
        this.closeLoginModal.addEventListener('click', () => {
            this.hideLoginModal();
        });

        // 点击弹窗外部关闭
        this.loginModal.addEventListener('click', (e) => {
            if (e.target === this.loginModal) {
                this.hideLoginModal();
            }
        });

        // 登录表单提交
        this.loginForm.addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleLogin();
        });

        // ESC键关闭弹窗
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && this.loginModal.style.display === 'block') {
                this.hideLoginModal();
            }
        });
    }

    /**
     * 显示登录弹窗
     */
    showLoginModal() {
        this.loginModal.style.display = 'block';
        
        // 聚焦到用户名输入框
        setTimeout(() => {
            document.getElementById('username').focus();
        }, 100);
    }

    /**
     * 隐藏登录弹窗
     */
    hideLoginModal() {
        this.loginModal.style.display = 'none';
        
        // 清空表单
        this.loginForm.reset();
    }

    /**
     * 处理登录逻辑
     */
    async handleLogin() {
        const username = document.getElementById('username').value.trim();
        const password = document.getElementById('password').value.trim();

        if (!username) {
            alert('请输入用户名');
            return;
        }

        try {
            // 显示加载状态
            this.setLoadingState(true);

            // 模拟登录API调用
            const success = await this.performLogin(username, password);

            if (success) {
                // 登录成功
                this.loginSuccess(username);
                this.hideLoginModal();
            } else {
                // 登录失败
                alert('用户名或密码错误，请重试');
            }
        } catch (error) {
            console.error('登录失败:', error);
            alert('登录失败，请重试');
        } finally {
            this.setLoadingState(false);
        }
    }

    /**
     * 执行登录
     * @param {string} username 用户名
     * @param {string} password 密码
     * @returns {Promise<boolean>} 登录是否成功
     */
    async performLogin(username, password) {
        try {
            console.log('=== 开始调用 im-demo 登录接口 ===');
            const loginUrl = `http://${window.IM_SERVICE_CONFIG.host}:${window.IM_SERVICE_CONFIG.port}/api/user/login`;
            console.log('请求URL:', loginUrl);
            console.log('请求数据:', { username: username });
            
            const response = await fetch(loginUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ username: username })
            });
            
            console.log('响应状态:', response.status);
            console.log('响应头:', Object.fromEntries(response.headers.entries()));
            
            const result = await response.json();
            console.log('响应数据:', result);
            console.log('=== im-demo 登录接口调用完成 ===');
            
            if (result.success) {
                // 使用 im-demo 返回的真实 UID
                this.currentUser = {
                    id: result.uid,  // 使用真实的 UID
                    username: username,
                    loginTime: new Date()
                };
                console.log('登录成功，用户信息:', this.currentUser);
                return true;
            } else {
                console.error('im-demo 登录失败:', result.message);
                return false;
            }
        } catch (error) {
            console.error('调用 im-demo 登录接口异常:', error);
            return false;
        }
    }

    /**
     * 登录成功处理
     * @param {string} username 用户名
     */
    loginSuccess(username) {
        this.isLoggedIn = true;
        

        // 保存登录状态到本地存储
        this.saveLoginState();

        // 更新界面显示
        this.updateLoginStatus();

        console.log('=== 登录成功后开始设置相关组件 ===');
        console.log('当前用户信息:', this.currentUser);
        
        // 设置聊天管理器的当前用户
        if (window.chatManager) {
            console.log('设置聊天管理器的当前用户...');
            window.chatManager.setCurrentUser(this.currentUser.id);
        } else {
            console.warn('聊天管理器未找到');
        }

        // 立即建立WebSocket连接
        if (window.wsManager) {
            console.log('开始建立WebSocket连接...');
            window.wsManager.connect(this.currentUser.id);
        } else {
            console.warn('WebSocket管理器未找到');
        }
        
        console.log('=== 登录成功后的组件设置完成 ===');

        // 显示欢迎消息
        this.showWelcomeMessage(username);

        console.log('用户登录成功:', this.currentUser);
    }

    /**
     * 生成用户ID
     * @param {string} username 用户名
     * @returns {string} 用户ID
     */
    generateUserId(username) {
        return `user_${username}_${Date.now()}`;
    }

    /**
     * 保存登录状态
     */
    saveLoginState() {
        if (this.currentUser) {
            localStorage.setItem('loginState', JSON.stringify({
                user: this.currentUser,
                timestamp: Date.now()
            }));
        }
    }

    /**
     * 检查登录状态
     */
    checkLoginStatus() {
        const savedState = localStorage.getItem('loginState');
        if (savedState) {
            try {
                const state = JSON.parse(savedState);
                const now = Date.now();
                const expireTime = 24 * 60 * 60 * 1000; // 24小时过期

                if (now - state.timestamp < expireTime) {
                    // 登录状态未过期，恢复登录
                    this.currentUser = state.user;
                    this.isLoggedIn = true;
                    this.updateLoginStatus();

                    // 设置聊天管理器的当前用户
                    if (window.chatManager) {
                        window.chatManager.setCurrentUser(this.currentUser.id);
                    }
                } else {
                    // 登录状态已过期，清除
                    this.clearLoginState();
                }
            } catch (error) {
                console.error('解析登录状态失败:', error);
                this.clearLoginState();
            }
        }
    }

    /**
     * 清除登录状态
     */
    clearLoginState() {
        localStorage.removeItem('loginState');
        this.isLoggedIn = false;
        this.currentUser = null;
        this.updateLoginStatus();
    }

    /**
     * 更新登录状态显示
     */
    updateLoginStatus() {
        if (this.isLoggedIn && this.currentUser) {
            this.loginStatus.textContent = `您好, ${this.currentUser.username}`;
            this.loginBtn.textContent = '退出登录';
            this.loginBtn.onclick = () => this.logout();
        } else {
            this.loginStatus.textContent = '您好,请登录';
            this.loginBtn.textContent = '免费注册';
            this.loginBtn.onclick = () => this.showLoginModal();
        }
    }

    /**
     * 显示欢迎消息
     * @param {string} username 用户名
     */
    showWelcomeMessage(username) {
        // 可以在这里显示一个欢迎通知
        console.log(`欢迎回来，${username}！`);
    }

    /**
     * 退出登录
     */
    logout() {
        if (confirm('确定要退出登录吗？')) {
            // 断开WebSocket连接
            if (window.wsManager) {
                window.wsManager.disconnect();
            }

            // 关闭聊天窗口
            if (window.chatManager && window.chatManager.getChatStatus()) {
                window.chatManager.closeChat();
            }

            // 清除登录状态
            this.clearLoginState();

            // 显示退出成功消息
            alert('已成功退出登录');
        }
    }

    /**
     * 设置加载状态
     * @param {boolean} isLoading 是否正在加载
     */
    setLoadingState(isLoading) {
        const submitBtn = document.querySelector('.submit-btn');
        if (submitBtn) {
            if (isLoading) {
                submitBtn.textContent = '登录中...';
                submitBtn.disabled = true;
            } else {
                submitBtn.textContent = '登录';
                submitBtn.disabled = false;
            }
        }
    }

    /**
     * 获取当前用户信息
     * @returns {Object|null} 当前用户信息
     */
    getCurrentUser() {
        return this.currentUser;
    }

    /**
     * 检查是否已登录
     * @returns {boolean} 是否已登录
     */
    isUserLoggedIn() {
        return this.isLoggedIn;
    }

    /**
     * 获取用户ID
     * @returns {string|null} 用户ID
     */
    getUserId() {
        return this.currentUser ? this.currentUser.id : null;
    }
}

// 创建全局登录管理器实例
window.loginManager = new LoginManager();
