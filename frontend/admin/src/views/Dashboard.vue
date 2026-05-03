<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-info">
              <div class="stat-value">{{ stats.documents }}</div>
              <div class="stat-label">档案总数</div>
            </div>
            <div class="stat-icon" style="background: linear-gradient(135deg, #409eff, #66b1ff);">
              <el-icon size="32"><Folder /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-info">
              <div class="stat-value">{{ stats.users }}</div>
              <div class="stat-label">用户总数</div>
            </div>
            <div class="stat-icon" style="background: linear-gradient(135deg, #67c23a, #85ce61);">
              <el-icon size="32"><User /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-info">
              <div class="stat-value">{{ stats.news }}</div>
              <div class="stat-label">新闻总数</div>
            </div>
            <div class="stat-icon" style="background: linear-gradient(135deg, #e6a23c, #ebb563);">
              <el-icon size="32"><News /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-info">
              <div class="stat-value">{{ stats.banners }}</div>
              <div class="stat-label">轮播图</div>
            </div>
            <div class="stat-icon" style="background: linear-gradient(135deg, #f56c6c, #f78989);">
              <el-icon size="32"><Picture /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <span class="card-title">欢迎使用档案管理系统</span>
          </template>
          <div class="welcome-content">
            <h3>您好，{{ userInfo?.real_name || userInfo?.username }}，欢迎回来！</h3>
            <p>当前角色：<el-tag :type="roleTagType">{{ roleText }}</el-tag></p>
            <el-divider />
            <div class="quick-actions">
              <h4>快捷操作</h4>
              <el-space>
                <el-button type="primary" @click="$router.push('/documents')">档案管理</el-button>
                <el-button type="success" @click="$router.push('/users')" v-if="isSuperAdmin">用户管理</el-button>
                <el-button type="warning" @click="$router.push('/news')" v-if="isSuperAdmin">新闻管理</el-button>
                <el-button type="danger" @click="$router.push('/banners')" v-if="isSuperAdmin">轮播图管理</el-button>
              </el-space>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <span class="card-title">系统信息</span>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="系统名称">档案管理系统</el-descriptions-item>
            <el-descriptions-item label="当前版本">v1.0.0</el-descriptions-item>
            <el-descriptions-item label="技术栈">Gin + Gorm + Vue3 + Element Plus</el-descriptions-item>
            <el-descriptions-item label="当前时间">{{ currentTime }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { Folder, User, News, Picture } from '@element-plus/icons-vue'

const userStore = useUserStore()
const userInfo = computed(() => userStore.userInfo)
const isSuperAdmin = computed(() => userStore.isSuperAdmin)

const stats = ref({
  documents: 0,
  users: 0,
  news: 0,
  banners: 0
})

const currentTime = ref('')
let timer = null

const roleTagType = computed(() => {
  const role = userStore.userInfo?.role
  switch (role) {
    case 'super_admin':
      return 'danger'
    case 'admin':
      return 'warning'
    default:
      return ''
  }
})

const roleText = computed(() => {
  const role = userStore.userInfo?.role
  switch (role) {
    case 'super_admin':
      return '超级管理员'
    case 'admin':
      return '管理员'
    default:
      return '用户'
  }
})

const updateTime = () => {
  const now = new Date()
  currentTime.value = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')} ${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}:${String(now.getSeconds()).padStart(2, '0')}`
}

onMounted(() => {
  updateTime()
  timer = setInterval(updateTime, 1000)
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style scoped>
.dashboard {
  max-width: 1400px;
  margin: 0 auto;
}

.stat-card {
  cursor: pointer;
  transition: transform 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
}

.stat-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-info .stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #303133;
}

.stat-info .stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 5px;
}

.stat-icon {
  width: 70px;
  height: 70px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.card-title {
  font-weight: bold;
  font-size: 16px;
}

.welcome-content h3 {
  margin-bottom: 15px;
  color: #303133;
}

.welcome-content h4 {
  margin-bottom: 15px;
  color: #606266;
}

.quick-actions {
  margin-top: 15px;
}
</style>
