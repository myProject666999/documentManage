import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userApi } from '@/api'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref(null)
  const token = ref(localStorage.getItem('token') || '')

  const isLoggedIn = computed(() => !!token.value)
  const isSuperAdmin = computed(() => userInfo.value?.role === 'super_admin')
  const isAdmin = computed(() => userInfo.value?.role === 'admin' || isSuperAdmin.value)

  async function login(username, password) {
    const res = await userApi.login({ username, password })
    if (res.code === 200) {
      const user = res.data.user
      if (user.role !== 'admin' && user.role !== 'super_admin') {
        return { code: 403, message: '您没有后台管理权限' }
      }
      token.value = res.data.token
      userInfo.value = user
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('userRole', user.role)
    }
    return res
  }

  async function getUserInfo() {
    if (!token.value) return
    const res = await userApi.getCurrentUser()
    if (res.code === 200) {
      userInfo.value = res.data
      localStorage.setItem('userRole', res.data.role)
    }
    return res
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userRole')
  }

  return {
    userInfo,
    token,
    isLoggedIn,
    isSuperAdmin,
    isAdmin,
    login,
    getUserInfo,
    logout
  }
})
