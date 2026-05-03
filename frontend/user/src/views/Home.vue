<template>
  <div class="home">
    <el-carousel height="400px" v-if="banners.length > 0">
      <el-carousel-item v-for="item in banners" :key="item.id">
        <div class="carousel-item" :style="{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }">
          <div class="banner-content">
            <h2>{{ item.title }}</h2>
            <p v-if="item.link_url">点击查看详情</p>
          </div>
        </div>
      </el-carousel-item>
    </el-carousel>

    <el-row :gutter="20" class="content-section">
      <el-col :span="16">
        <el-card shadow="hover" class="section-card">
          <template #header>
            <div class="card-header">
              <span>最新档案</span>
              <el-button type="primary" link @click="$router.push('/documents')">查看更多</el-button>
            </div>
          </template>
          <el-table :data="documents" style="width: 100%">
            <el-table-column prop="title" label="档案标题">
              <template #default="{ row }">
                <el-link type="primary" @click="$router.push(`/documents/${row.id}`)">{{ row.title }}</el-link>
              </template>
            </el-table-column>
            <el-table-column prop="document_number" label="档案编号" width="150" />
            <el-table-column prop="author" label="作者" width="120" />
            <el-table-column prop="view_count" label="浏览量" width="100" />
            <el-table-column prop="created_at" label="发布时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card shadow="hover" class="section-card">
          <template #header>
            <div class="card-header">
              <span>最新公告</span>
              <el-button type="primary" link @click="$router.push('/announcements')">查看更多</el-button>
            </div>
          </template>
          <ul class="news-list">
            <li v-for="item in announcements" :key="item.id" @click="$router.push(`/announcements/${item.id}`)">
              <span class="title" :class="{ 'is-top': item.is_top }">{{ item.title }}</span>
              <span class="date">{{ formatDate(item.created_at) }}</span>
            </li>
          </ul>
        </el-card>

        <el-card shadow="hover" class="section-card" style="margin-top: 20px;">
          <template #header>
            <div class="card-header">
              <span>新闻资讯</span>
              <el-button type="primary" link @click="$router.push('/news')">查看更多</el-button>
            </div>
          </template>
          <ul class="news-list">
            <li v-for="item in news" :key="item.id" @click="$router.push(`/news/${item.id}`)">
              <span class="title" :class="{ 'is-top': item.is_top }">{{ item.title }}</span>
              <span class="date">{{ formatDate(item.created_at) }}</span>
            </li>
          </ul>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { bannerApi, documentApi, announcementApi, newsApi } from '@/api'

const banners = ref([])
const documents = ref([])
const announcements = ref([])
const news = ref([])

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

onMounted(async () => {
  try {
    const [bannerRes, docRes, annRes, newsRes] = await Promise.all([
      bannerApi.getList(),
      documentApi.getList({ page: 1, page_size: 5 }),
      announcementApi.getList({ page: 1, page_size: 5 }),
      newsApi.getList({ page: 1, page_size: 5 })
    ])
    
    if (bannerRes.code === 200) banners.value = bannerRes.data
    if (docRes.code === 200) documents.value = docRes.data.list
    if (annRes.code === 200) announcements.value = annRes.data.list
    if (newsRes.code === 200) news.value = newsRes.data.list
  } catch (error) {
    console.error(error)
  }
})
</script>

<style scoped>
.home {
  max-width: 1400px;
  margin: 0 auto;
}

.carousel-item {
  display: flex;
  justify-content: center;
  align-items: center;
  color: #fff;
}

.banner-content {
  text-align: center;
}

.banner-content h2 {
  font-size: 36px;
  margin-bottom: 10px;
}

.banner-content p {
  font-size: 18px;
  opacity: 0.8;
}

.content-section {
  margin-top: 20px;
}

.section-card {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: bold;
  font-size: 16px;
}

.news-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.news-list li {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #eee;
  cursor: pointer;
}

.news-list li:last-child {
  border-bottom: none;
}

.news-list li:hover .title {
  color: #667eea;
}

.news-list .title {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #333;
}

.news-list .title.is-top {
  color: #e6a23c;
}

.news-list .date {
  color: #999;
  font-size: 12px;
  margin-left: 10px;
  flex-shrink: 0;
}
</style>
