<template>
  <div class="profile">
    <el-card shadow="hover">
      <template #header>
        <span class="card-title">个人信息</span>
      </template>

      <el-form
        :model="form"
        :rules="rules"
        ref="formRef"
        label-width="100px"
        style="max-width: 500px; margin: 0 auto;"
      >
        <el-form-item label="用户名">
          <el-input v-model="form.username" disabled />
        </el-form-item>
        <el-form-item label="角色">
          <el-tag :type="roleType">{{ roleText }}</el-tag>
        </el-form-item>
        <el-form-item label="真实姓名" prop="realName">
          <el-input v-model="form.realName" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="所属部门">
          <el-tag v-if="form.department">{{ form.department.name }}</el-tag>
          <span v-else class="text-gray">暂无</span>
        </el-form-item>
        <el-form-item label="注册时间">
          <span>{{ formatDate(form.created_at) }}</span>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">保存修改</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  username: '',
  realName: '',
  email: '',
  phone: '',
  department: null,
  created_at: ''
})

const rules = {
  email: [
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

const roleType = computed(() => {
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
      return '普通用户'
  }
})

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const loadUserInfo = async () => {
  if (userStore.userInfo) {
    form.username = userStore.userInfo.username
    form.realName = userStore.userInfo.real_name || ''
    form.email = userStore.userInfo.email || ''
    form.phone = userStore.userInfo.phone || ''
    form.department = userStore.userInfo.department
    form.created_at = userStore.userInfo.created_at
  }
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const res = await userStore.updateProfile({
      real_name: form.realName,
      email: form.email,
      phone: form.phone
    })
    if (res.code === 200) {
      ElMessage.success('修改成功')
    } else {
      ElMessage.error(res.message || '修改失败')
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadUserInfo()
})
</script>

<style scoped>
.profile {
  max-width: 800px;
  margin: 0 auto;
}

.card-title {
  font-weight: bold;
  font-size: 16px;
}

.text-gray {
  color: #999;
}
</style>
