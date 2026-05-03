import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

const request = axios.create({
  baseURL: '/api',
  timeout: 30000
})

request.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

request.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          ElMessage.error('登录已过期，请重新登录')
          localStorage.removeItem('token')
          router.push('/login')
          break
        case 403:
          ElMessage.error('没有权限访问')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器错误')
          break
        default:
          ElMessage.error(error.response.data?.message || '请求失败')
      }
    } else {
      ElMessage.error('网络错误，请检查网络连接')
    }
    return Promise.reject(error)
  }
)

export const userApi = {
  login: (data) => request.post('/login', data),
  register: (data) => request.post('/register', data),
  getCurrentUser: () => request.get('/user/info'),
  updateProfile: (data) => request.put('/user/profile', data),
  changePassword: (data) => request.put('/user/password', data)
}

export const documentApi = {
  getList: (params) => request.get('/documents', { params }),
  getDetail: (id) => request.get(`/documents/${id}`),
  getTypes: () => request.get('/document-types')
}

export const announcementApi = {
  getList: (params) => request.get('/announcements', { params }),
  getDetail: (id) => request.get(`/announcements/${id}`)
}

export const newsApi = {
  getList: (params) => request.get('/news', { params }),
  getDetail: (id) => request.get(`/news/${id}`)
}

export const bannerApi = {
  getList: () => request.get('/banners')
}

export const favoriteApi = {
  getList: (params) => request.get('/favorites', { params }),
  add: (data) => request.post('/favorites', data),
  remove: (documentId) => request.delete(`/favorites/${documentId}`)
}

export const departmentApi = {
  getList: () => request.get('/departments')
}

export default request
