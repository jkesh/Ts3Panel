<template>
  <div class="login-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <h2>注册新账号</h2>
        </div>
      </template>
      <el-form :model="form" label-width="0">
        <el-form-item>
          <el-input v-model="form.username" placeholder="请输入用户名" prefix-icon="User" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="form.password" type="password" placeholder="请输入密码" prefix-icon="Lock" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="success" @click="handleRegister" :loading="loading" style="width: 100%">立即注册</el-button>
        </el-form-item>
        <div style="text-align: center; margin-top: 10px;">
          <router-link to="/login" style="text-decoration: none; color: #409eff;">
            已有账号？返回登录
          </router-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'
import { ElMessage } from 'element-plus'

const router = useRouter()
const loading = ref(false)
const form = ref({ username: '', password: '' })

const handleRegister = async () => {
  if (!form.value.username || !form.value.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }

  loading.value = true
  try {
    // 调用后端的公开注册接口
    await api.post('/auth/register', form.value)
    ElMessage.success('注册成功，请登录')
    // 注册成功后跳转到登录页
    router.push('/login')
  } catch (error) {
    // 错误通常由 axios 拦截器处理，但也可以在这里补充
    console.error(error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f0f2f5;
}
.box-card {
  width: 400px;
}
.card-header {
  text-align: center;
}
</style>