<template>
  <div class="documents">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <el-form :inline="true" :model="searchForm">
            <el-form-item label="关键词">
              <el-input v-model="searchForm.title" placeholder="请输入关键词" clearable />
            </el-form-item>
            <el-form-item label="状态">
              <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
                <el-option label="启用" :value="1" />
                <el-option label="禁用" :value="0" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSearch">搜索</el-button>
              <el-button @click="resetForm">重置</el-button>
            </el-form-item>
          </el-form>
          <el-button type="primary" @click="handleAdd">新增档案</el-button>
        </div>
      </template>

      <el-table :data="list" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="档案标题" min-width="200" show-overflow-tooltip />
        <el-table-column prop="document_number" label="档案编号" width="150" />
        <el-table-column prop="document_type.name" label="档案类型" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.document_type" size="small">{{ row.document_type.name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="author" label="作者" width="100" />
        <el-table-column prop="view_count" label="浏览量" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchList"
        @current-change="fetchList"
      />
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑档案' : '新增档案'"
      width="700px"
      :close-on-click-modal="false"
    >
      <el-form
        :model="form"
        :rules="rules"
        ref="formRef"
        label-width="100px"
        style="max-width: 600px;"
      >
        <el-form-item label="档案标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入档案标题" />
        </el-form-item>
        <el-form-item label="档案编号" prop="document_number">
          <el-input v-model="form.document_number" placeholder="请输入档案编号" />
        </el-form-item>
        <el-form-item label="档案类型" prop="document_type_id">
          <el-select v-model="form.document_type_id" placeholder="请选择档案类型" style="width: 100%">
            <el-option
              v-for="item in documentTypes"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="所属部门" prop="department_id">
          <el-select v-model="form.department_id" placeholder="请选择所属部门" clearable style="width: 100%">
            <el-option
              v-for="item in departments"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="作者" prop="author">
          <el-input v-model="form.author" placeholder="请输入作者" />
        </el-form-item>
        <el-form-item label="关键词" prop="keyword">
          <el-input v-model="form.keyword" placeholder="请输入关键词，多个用逗号分隔" />
        </el-form-item>
        <el-form-item label="附件链接" prop="file_url">
          <el-input v-model="form.file_url" placeholder="请输入附件链接" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="5"
            placeholder="请输入档案内容"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { documentApi, departmentApi } from '@/api'

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)

const list = ref([])
const documentTypes = ref([])
const departments = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const searchForm = reactive({
  title: '',
  status: ''
})

const form = reactive({
  id: null,
  title: '',
  document_number: '',
  document_type_id: null,
  department_id: null,
  author: '',
  keyword: '',
  file_url: '',
  content: '',
  status: 1
})

const rules = {
  title: [{ required: true, message: '请输入档案标题', trigger: 'blur' }],
  document_type_id: [{ required: true, message: '请选择档案类型', trigger: 'change' }]
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const fetchList = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (searchForm.title) params.title = searchForm.title
    if (searchForm.status !== '') params.status = searchForm.status

    const res = await documentApi.getList(params)
    if (res.code === 200) {
      list.value = res.data.list
      total.value = res.data.total
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
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

const fetchDepartments = async () => {
  try {
    const res = await departmentApi.getList()
    if (res.code === 200) {
      departments.value = res.data
    }
  } catch (error) {
    console.error(error)
  }
}

const handleSearch = () => {
  page.value = 1
  fetchList()
}

const resetForm = () => {
  searchForm.title = ''
  searchForm.status = ''
  handleSearch()
}

const handleAdd = () => {
  isEdit.value = false
  resetFormData()
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  form.id = row.id
  form.title = row.title
  form.document_number = row.document_number
  form.document_type_id = row.document_type_id
  form.department_id = row.department_id
  form.author = row.author
  form.keyword = row.keyword
  form.file_url = row.file_url
  form.content = row.content
  form.status = row.status
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该档案吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const res = await documentApi.delete(row.id)
    if (res.code === 200) {
      ElMessage.success('删除成功')
      fetchList()
    } else {
      ElMessage.error(res.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

const resetFormData = () => {
  form.id = null
  form.title = ''
  form.document_number = ''
  form.document_type_id = null
  form.department_id = null
  form.author = ''
  form.keyword = ''
  form.file_url = ''
  form.content = ''
  form.status = 1
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitLoading.value = true
  try {
    const data = {
      title: form.title,
      document_number: form.document_number,
      document_type_id: form.document_type_id,
      department_id: form.department_id,
      author: form.author,
      keyword: form.keyword,
      file_url: form.file_url,
      content: form.content,
      status: form.status
    }

    let res
    if (isEdit.value) {
      res = await documentApi.update(form.id, data)
    } else {
      res = await documentApi.create(data)
    }

    if (res.code === 200) {
      ElMessage.success(isEdit.value ? '编辑成功' : '新增成功')
      dialogVisible.value = false
      fetchList()
    } else {
      ElMessage.error(res.message || '操作失败')
    }
  } catch (error) {
    console.error(error)
  } finally {
    submitLoading.value = false
  }
}

onMounted(() => {
  fetchList()
  fetchDocumentTypes()
  fetchDepartments()
})
</script>

<style scoped>
.documents {
  max-width: 1400px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
