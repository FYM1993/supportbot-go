/**
 * 客服客户端配置文件
 * 在这里配置后端服务地址和端口
 */

// im-demo服务配置
const IM_SERVICE_CONFIG = {
    // 服务IP地址 - 修改为你的实际IP地址
    host: 'localhost', // 例如: '192.168.1.100', '10.0.0.50'
    
    // 服务端口
    port: '11005',
    
    // WebSocket端点
    wsEndpoint: '/ws'
};

// 获取完整的WebSocket URL
function getWebSocketUrl(uid = null) {
    console.log('getWebSocketUrl 被调用，参数 uid:', uid, '类型:', typeof uid);
    const protocol = window.location.protocol === 'https:' ? 'https:' : 'http:';
    let url = `${protocol}//${IM_SERVICE_CONFIG.host}:${IM_SERVICE_CONFIG.port}${IM_SERVICE_CONFIG.wsEndpoint}`;
    
    // 如果提供了uid，添加到URL参数中
    if (uid) {
        url += `?uid=${uid}`;
        console.log('添加uid参数后的URL:', url);
    } else {
        console.log('没有提供uid参数，使用基础URL:', url);
    }
    
    return url;
}

// 导出配置
window.IM_SERVICE_CONFIG = IM_SERVICE_CONFIG;
window.getWebSocketUrl = getWebSocketUrl;

// 为了兼容性，也导出为 config
window.config = {
    host: IM_SERVICE_CONFIG.host,
    port: IM_SERVICE_CONFIG.port,
    wsEndpoint: IM_SERVICE_CONFIG.wsEndpoint
};
