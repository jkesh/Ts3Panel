import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '../router'

const api = axios.create({
    baseURL: 'http://localhost:8080', // 后端地址
    timeout: 5000
})

// 请求拦截器：注入 Token
api.interceptors.request.use(config => {
    const token = localStorage.getItem('token')
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})

// 响应拦截器：处理错误
api.interceptors.response.use(
    response => response,
    error => {
        if (error.response && error.response.status === 401) {
            ElMessage.error('登录过期，请重新登录')
            localStorage.removeItem('token')
            router.push('/login')
        } else {
            ElMessage.error(error.response?.data?.error || '网络错误')
        }
        return Promise.reject(error)
    }
)

export default api