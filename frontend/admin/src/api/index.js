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
          localStorage.removeItem('userRole')
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
  getCurrentUser: () => request.get('/user/info')
}

export const documentApi = {
  getList: (params) => request.get('/admin/documents', { params }),
  getDetail: (id) => request.get(`/admin/documents/${id}`),
  create: (data) => request.post('/admin/documents', data),
  update: (id, data) => request.put(`/admin/documents/${id}`, data),
  delete: (id) => request.delete(`/admin/documents/${id}`),
  getTypes: () => request.get('/document-types')
}

export const adminUserApi = {
  getList: (params) => request.get('/super-admin/users', { params }),
  getDetail: (id) => request.get(`/super-admin/users/${id}`),
  create: (data) => request.post('/super-admin/users', data),
  update: (id, data) => request.put(`/super-admin/users/${id}`, data),
  delete: (id) => request.delete(`/super-admin/users/${id}`),
  resetPassword: (id, data) => request.post(`/super-admin/users/${id}/reset-password`, data)
}

export const departmentApi = {
  getList: (params) => request.get('/super-admin/departments', { params }),
  create: (data) => request.post('/super-admin/departments', data),
  update: (id, data) => request.put(`/super-admin/departments/${id}`, data),
  delete: (id) => request.delete(`/super-admin/departments/${id}`)
}

export const documentTypeApi = {
  getList: (params) => request.get('/super-admin/document-types', { params }),
  create: (data) => request.post('/super-admin/document-types', data),
  update: (id, data) => request.put(`/super-admin/document-types/${id}`, data),
  delete: (id) => request.delete(`/super-admin/document-types/${id}`)
}

export const newsApi = {
  getList: (params) => request.get('/super-admin/news', { params }),
  getDetail: (id) => request.get(`/super-admin/news/${id}`),
  create: (data) => request.post('/super-admin/news', data),
  update: (id, data) => request.put(`/super-admin/news/${id}`, data),
  delete: (id) => request.delete(`/super-admin/news/${id}`)
}

export const bannerApi = {
  getList: (params) => request.get('/super-admin/banners', { params }),
  create: (data) => request.post('/super-admin/banners', data),
  update: (id, data) => request.put(`/super-admin/banners/${id}`, data),
  delete: (id) => request.delete(`/super-admin/banners/${id}`)
}

export const announcementApi = {
  getList: (params) => request.get('/super-admin/announcements', { params }),
  create: (data) => request.post('/super-admin/announcements', data),
  update: (id, data) => request.put(`/super-admin/announcements/${id}`, data),
  delete: (id) => request.delete(`/super-admin/announcements/${id}`)
}

export default request
