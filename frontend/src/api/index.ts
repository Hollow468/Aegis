import axios from 'axios'
import { ElMessage } from 'element-plus'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    const status = error.response?.status
    const message = error.response?.data?.message || error.message

    if (status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    } else if (status === 429) {
      ElMessage.warning('Rate limit exceeded')
    } else if (status === 503) {
      ElMessage.error('Service unavailable (circuit breaker open)')
    } else {
      ElMessage.error(message || 'Request failed')
    }

    return Promise.reject(error)
  }
)

export default api
