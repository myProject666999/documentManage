<template>
  <div class="favorites">
    <el-card shadow="hover">
      <template #header>
        <span class="card-title">我的收藏</span>
      </template>

      <el-empty v-if="favorites.length === 0" description="暂无收藏" />

      <el-table v-else :data="favorites" style="width: 100%">
        <el-table-column prop="document.title" label="档案标题" min-width="300">
          <template #default="{ row }">
            <el-link type="primary" @click="$router.push(`/documents/${row.document_id}`)">{{ row.document?.title }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="document.document_number" label="档案编号" width="150">
          <template #default="{ row }">
            {{ row.document?.document_number }}
          </template>
        </el-table-column>
        <el-table-column prop="document.document_type.name" label="档案类型" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.document?.document_type">{{ row.document.document_type.name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="document.view_count" label="浏览量" width="100">
          <template #default="{ row }">
            {{ row.document?.view_count }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="收藏时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button type="danger" link @click="removeFavorite(row.document_id)">取消收藏</el-button>
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { favoriteApi } from '@/api'

const favorites = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const fetchFavorites = async () => {
  try {
    const res = await favoriteApi.getList({
      page: page.value,
      page_size: pageSize.value
    })
    if (res.code === 200) {
      favorites.value = res.data.list
      total.value = res.data.total
    }
  } catch (error) {
    console.error(error)
  }
}

const removeFavorite = async (documentId) => {
  try {
    await ElMessageBox.confirm('确定要取消收藏吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const res = await favoriteApi.remove(documentId)
    if (res.code === 200) {
      ElMessage.success('取消收藏成功')
      fetchFavorites()
    } else {
      ElMessage.error(res.message || '取消收藏失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

const handleSizeChange = (val) => {
  pageSize.value = val
  fetchFavorites()
}

const handleCurrentChange = (val) => {
  page.value = val
  fetchFavorites()
}

onMounted(() => {
  fetchFavorites()
})
</script>

<style scoped>
.favorites {
  max-width: 1200px;
  margin: 0 auto;
}

.card-title {
  font-weight: bold;
  font-size: 16px;
}
</style>
