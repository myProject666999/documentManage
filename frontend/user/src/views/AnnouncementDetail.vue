<template>
  <div class="announcement-detail">
    <el-card shadow="hover" v-loading="loading">
      <template #header>
        <div class="card-header">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item :to="{ path: '/announcements' }">公告</el-breadcrumb-item>
            <el-breadcrumb-item>公告详情</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
      </template>

      <div v-if="announcement" class="detail-content">
        <h1 class="title">{{ announcement.title }}</h1>
        <div class="meta">
          <span v-if="announcement.is_top"><el-tag type="warning">置顶</el-tag></span>
          <span>发布时间：{{ formatDate(announcement.created_at) }}</span>
        </div>
        <el-divider />
        <div class="content" v-html="announcement.content"></div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { announcementApi } from '@/api'

const route = useRoute()
const router = useRouter()
const announcement = ref(null)
const loading = ref(false)

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const fetchAnnouncement = async () => {
  loading.value = true
  try {
    const res = await announcementApi.getDetail(route.params.id)
    if (res.code === 200) {
      announcement.value = res.data
    } else {
      ElMessage.error(res.message || '获取公告详情失败')
      router.push('/announcements')
    }
  } catch (error) {
    console.error(error)
    router.push('/announcements')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchAnnouncement()
})
</script>

<style scoped>
.announcement-detail {
  max-width: 1000px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.detail-content {
  padding: 10px;
}

.title {
  font-size: 24px;
  color: #333;
  margin-bottom: 20px;
  text-align: center;
}

.meta {
  display: flex;
  gap: 20px;
  color: #666;
  font-size: 14px;
  margin-bottom: 20px;
  justify-content: center;
}

.content {
  font-size: 16px;
  line-height: 2;
  color: #333;
  white-space: pre-wrap;
}
</style>
