<template>
  <div class="announcements">
    <el-card shadow="hover">
      <template #header>
        <span class="card-title">公告列表</span>
      </template>

      <el-table :data="announcements" style="width: 100%">
        <el-table-column prop="title" label="公告标题" min-width="300">
          <template #default="{ row }">
            <el-link type="primary" @click="$router.push(`/announcements/${row.id}`)">
              <span v-if="row.is_top" style="color: #e6a23c;">[置顶]</span>
              {{ row.title }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="发布时间" width="200">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        style="margin-top: 20px; justify-content: flex-end;"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { announcementApi } from '@/api'

const announcements = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const fetchAnnouncements = async () => {
  try {
    const res = await announcementApi.getList({
      page: page.value,
      page_size: pageSize.value
    })
    if (res.code === 200) {
      announcements.value = res.data.list
      total.value = res.data.total
    }
  } catch (error) {
    console.error(error)
  }
}

const handleSizeChange = (val) => {
  pageSize.value = val
  fetchAnnouncements()
}

const handleCurrentChange = (val) => {
  page.value = val
  fetchAnnouncements()
}

onMounted(() => {
  fetchAnnouncements()
})
</script>

<style scoped>
.announcements {
  max-width: 1000px;
  margin: 0 auto;
}

.card-title {
  font-weight: bold;
  font-size: 16px;
}
</style>
