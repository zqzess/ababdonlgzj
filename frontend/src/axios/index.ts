/**
 * ========================
 * Created with WebStorm.
 * @Author zqzess
 * @Date 2024/10/14 17:33
 * @File abandonlgzj/index.ts
 * @Version :
 * @Desc :
 * @GitHUb Https://github.com/zqzess
 * ========================
 **/
// src/services/api.ts
import axios from 'axios';

// 获取浏览器当前 URL 的 IP 地址
export const getApiBaseUrl = (port2?) => {
    const { protocol, hostname, port } = window.location;
    // 如果有端口号，则包含在 URL 中
    if (port2 === null) {
        return `${protocol}//${hostname}${port ? `:${port}` : ''}`;
    } else {
        return `${protocol}//${hostname}:${port2}`;
    }

};

// 创建一个 axios 实例
const apiClient = axios.create({
    baseURL: getApiBaseUrl(9091), // 动态设置基础 URL
    timeout: 1000,
});

// 导出 API 客户端
export default apiClient
