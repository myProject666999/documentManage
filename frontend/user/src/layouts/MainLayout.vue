<template>
  <div class="main-layout">
    <el-header class="header">
      <div class="logo" @click="$router.push('/')">
        <el-icon><Document /></el-icon>
        <span>档案管理系统</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="nav-menu"
        mode="horizontal"
        background-color="transparent"
        text-color="#fff"
        active-text-color="#ffd04b"
        router
      >
        <el-menu-item index="/">首页</el-menu-item>
        <el-menu-item index="/documents">档案信息</el-menu-item>
        <el-menu-item index="/announcements">公告</el-menu-item>
        <el-menu-item index="/news">新闻资讯</el-menu-item>
      </el-menu>
      <div class="user-actions">
        <template v-if="isLoggedIn">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="32">{{ userInfo?.real_name?.charAt(0) || userInfo?.username?.charAt(0) }}</el-avatar>
              <span>{{ userInfo?.real_name || userInfo?.username }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                <el-dropdown-item command="favorites">我的收藏</el-dropdown-item>
                <el-dropdown-item command="changePassword">修改密码</el-dropdown-item>
                <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
        <template v-else>
          <el-button type="primary" link @click="$router.push('/login')">登录</el-button>
          <el-button type="primary" @click="$router.push('/register')">注册</el-button>
        </template>
      </div>
    </el-header>
    <el-main class="main">
      <router-view />
    </el-main>
    <el-footer class="footer">
      <p>© 2024 档案管理系统 版权所有</p>
    </el-footer>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { Document } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)
const isLoggedIn = computed(() => userStore.isLoggedIn)
const userInfo = computed(() => userStore.userInfo)

onMounted(() => {
  if (isLoggedIn.value) {
    userStore.getUserInfo()
  }
})

const handleCommand = (command) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'favorites':
      router.push('/favorites')
      break
    case 'changePassword':
      router.push('/change-password')
      break
    case 'logout':
      userStore.logout()
      router.push('/')
      break
  }
}
</script>

<style scoped>
.main-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  padding: 0 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #fff;
  font-size: 20px;
  font-weight: bold;
  cursor: pointer;
}

.nav-menu {
  margin-left: 40px;
  border-bottom: none;
  flex: 1;
}

.nav-menu .el-menu-item {
  color: #fff;
}

.user-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #fff;
  cursor: pointer;
}

.main {
  flex: 1;
  padding: 20px;
  background-color: #f5f5f5;
}

.footer {
  background-color: #333;
  color: #fff;
  text-align: center;
  padding: 20px;
}
</style>
