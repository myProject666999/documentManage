<template>
  <div class="documents">
    <el-card shadow="never">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.title" placeholder="请输入关键词" clearable @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="档案类型">
          <el-select v-model="searchForm.document_type_id" placeholder="请选择类型" clearable>
            <el-option
              v-for="item in documentTypes"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="hover" style="margin-top: 20px;">
      <el-table :data="documents" style="width: 100%">
        <el-table-column prop="title" label="档案标题" min-width="200">
          <template #default="{ row }">
            <el-link type="primary" @click="$router.push(`/documents/${row.id}`)">{{ row.title }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="document_number" label="档案编号" width="150" />
        <el-table-column prop="document_type.name" label="档案类型" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.document_type">{{ row.document_type.name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="author" label="作者" width="100" />
        <el-table-column prop="view_count" label="浏览量" width="100" />
        <el-table-column prop="created_at" label="发布时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button type="primary" link @click="addFavorite(row.id)">收藏</el-button>
            <el-button type="primary" link @click="$router.push(`/documents/${row.id}`)">详情</el-button>
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
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { documentApi, favoriteApi } from '@/api'

const router = useRouter()
const documents = ref([])
const documentTypes = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const searchForm = reactive({
  title: '',
  document_type_id: ''
})

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const fetchDocuments = async () => {
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (searchForm.title) params.title = searchForm.title
    if (searchForm.document_type_id) params.document_type_id = searchForm.document_type_id

    const res = await documentApi.getList(params)
    if (res.code === 200) {
      documents.value = res.data.list
      total.value = res.data.total
    }
  } catch (error) {
    console.error(error)
  }
}

const fetchDocumentTypes = async () => {
  try {
    const res = await documentApi.getTypes()
    if (res.code === 200) {
      documentTypes.value = res.data
    }
  } catch (error) {
    console.error(error)
  }
}

const handleSearch = () => {
  page.value = 1
  fetchDocuments()
}

const resetForm = () => {
  searchForm.title = ''
  searchForm.document_type_id = ''
  handleSearch()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  fetchDocuments()
}

const handleCurrentChange = (val) => {
  page.value = val
  fetchDocuments()
}

const addFavorite = async (documentId) => {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  try {
    const res = await favoriteApi.add({ document_id: documentId })
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
  fetchDocuments()
  fetchDocumentTypes()
})
</script>

<style scoped>
.documents {
  max-width: 1400px;
  margin: 0 auto;
}

.search-form {
  margin-bottom: 10px;
}
</style>
