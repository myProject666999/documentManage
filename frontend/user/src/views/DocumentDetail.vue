<template>
  <div class="document-detail">
    <el-card shadow="hover" v-loading="loading">
      <template #header>
        <div class="card-header">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item :to="{ path: '/documents' }">档案信息</el-breadcrumb-item>
            <el-breadcrumb-item>档案详情</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
      </template>

      <div v-if="document" class="detail-content">
        <h1 class="title">{{ document.title }}</h1>
        <div class="meta">
          <span v-if="document.document_number">档案编号：{{ document.document_number }}</span>
          <span v-if="document.document_type">
            <el-tag>{{ document.document_type.name }}</el-tag>
          </span>
          <span v-if="document.author">作者：{{ document.author }}</span>
          <span>浏览量：{{ document.view_count }}</span>
          <span>发布时间：{{ formatDate(document.created_at) }}</span>
        </div>
        <div class="actions">
          <el-button type="primary" @click="addFavorite">
            <el-icon><Star /></el-icon>
            收藏
          </el-button>
          <el-button v-if="document.file_url" type="success">
            <el-icon><Download /></el-icon>
            下载附件
          </el-button>
        </div>
        <el-divider />
        <div class="content" v-html="document.content"></div>
        <el-divider />
        <div class="additional-info">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-descriptions title="基本信息" :column="1" border>
                <el-descriptions-item label="档案编号">{{ document.document_number || '-' }}</el-descriptions-item>
                <el-descriptions-item label="档案类型">{{ document.document_type?.name || '-' }}</el-descriptions-item>
                <el-descriptions-item label="所属部门">{{ document.department?.name || '-' }}</el-descriptions-item>
              </el-descriptions>
            </el-col>
            <el-col :span="12">
              <el-descriptions title="其他信息" :column="1" border>
                <el-descriptions-item label="作者">{{ document.author || '-' }}</el-descriptions-item>
                <el-descriptions-item label="关键词">{{ document.keyword || '-' }}</el-descriptions-item>
                <el-descriptions-item label="发布时间">{{ formatDate(document.created_at) }}</el-descriptions-item>
              </el-descriptions>
            </el-col>
          </el-row>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { documentApi, favoriteApi } from '@/api'
import { Star, Download } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const document = ref(null)
const loading = ref(false)

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const fetchDocument = async () => {
  loading.value = true
  try {
    const res = await documentApi.getDetail(route.params.id)
    if (res.code === 200) {
      document.value = res.data
    } else {
      ElMessage.error(res.message || '获取档案详情失败')
      router.push('/documents')
    }
  } catch (error) {
    console.error(error)
    router.push('/documents')
  } finally {
    loading.value = false
  }
}

const addFavorite = async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  try {
    const res = await favoriteApi.add({ document_id: document.value.id })
    if (res.code === 200) {
      ElMessage.success('收藏成功')
    } else {
      ElMessage.error(res.message || '收藏失败')
    }
  } catch (error) {
    console.error(error)
  }
}

onMounted(() => {
  fetchDocument()
})
</script>

<style scoped>
.document-detail {
  max-width: 1200px;
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
  flex-wrap: wrap;
  gap: 20px;
  color: #666;
  font-size: 14px;
  margin-bottom: 20px;
}

.actions {
  margin-bottom: 20px;
}

.content {
  font-size: 16px;
  line-height: 2;
  color: #333;
  white-space: pre-wrap;
}

.additional-info {
  margin-top: 20px;
}
</style>
