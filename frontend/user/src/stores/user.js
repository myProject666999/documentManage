import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userApi } from '@/api'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref(null)
  const token = ref(localStorage.getItem('token') || '')

  const isLoggedIn = computed(() => !!token.value)

  async function login(username, password) {
    const res = await userApi.login({ username, password })
    if (res.code === 200) {
      token.value = res.data.token
      userInfo.value = res.data.user
      localStorage.setItem('token', res.data.token)
    }
    return res
  }

  async function register(data) {
    return await userApi.register(data)
  }

  async function getUserInfo() {
    if (!token.value) return
    const res = await userApi.getCurrentUser()
    if (res.code === 200) {
      userInfo.value = res.data
    }
    return res
  }

  async function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  async function updateProfile(data) {
    const res = await userApi.updateProfile(data)
    if (res.code === 200 && userInfo.value) {
      userInfo.value = { ...userInfo.value, ...data }
    }
    return res
  }

  async function changePassword(data) {
    return await userApi.changePassword(data)
  }

  return {
    userInfo,
    token,
    isLoggedIn,
    login,
    register,
    getUserInfo,
    logout,
    updateProfile,
    changePassword
  }
})
