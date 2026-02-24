<template>
  <div class="server-manager">
    <el-card shadow="never" class="page-header">
      <div class="header-main">
        <div>
          <h2>服务器总管理</h2>
          <p>集中管理服务器配置、临时密码、Query 账号与组成员分配。</p>
        </div>
        <el-button type="primary" :icon="Refresh" @click="initData">刷新全部</el-button>
      </div>
    </el-card>

    <el-row :gutter="16" class="content-grid">
      <el-col :xs="24" :lg="14">
        <el-card>
          <template #header><b>服务器基础配置</b></template>
          <el-form :model="serverForm" label-width="110px">
            <el-form-item label="服务器名称"><el-input v-model="serverForm.name" /></el-form-item>
            <el-form-item label="欢迎消息"><el-input v-model="serverForm.welcome_message" type="textarea" :rows="2" /></el-form-item>
            <el-form-item label="服务器密码"><el-input v-model="serverForm.password" show-password placeholder="留空不修改" /></el-form-item>
            <el-form-item label="最大人数"><el-input-number v-model="serverForm.max_clients" :min="1" :max="2048" style="width:100%" /></el-form-item>
            <el-form-item label="主机消息"><el-input v-model="serverForm.host_message" type="textarea" :rows="2" /></el-form-item>
            <el-form-item label="消息模式">
              <el-select v-model="serverForm.host_message_mode" style="width:100%">
                <el-option label="不显示" :value="0" />
                <el-option label="登录后显示" :value="1" />
                <el-option label="强制显示" :value="2" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveServerSettings">保存配置</el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <el-card class="mt-16">
          <template #header><b>用户组成员管理</b></template>
          <el-form :model="groupClientForm" inline>
            <el-form-item label="服务器组ID"><el-input-number v-model="groupClientForm.server_group_id" :min="1" /></el-form-item>
            <el-form-item label="用户DBID"><el-input-number v-model="groupClientForm.client_db_id" :min="1" /></el-form-item>
            <el-form-item>
              <el-button type="success" @click="changeGroupClient('add')">添加成员</el-button>
              <el-button type="danger" plain @click="changeGroupClient('remove')">移除成员</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="10">
        <el-card>
          <template #header>
            <div class="card-header-inline">
              <b>临时密码</b>
              <el-button type="primary" size="small" @click="createTempPassword">新增</el-button>
            </div>
          </template>
          <el-form :model="tempPasswordForm" label-width="90px" size="small">
            <el-form-item label="密码"><el-input v-model="tempPasswordForm.password" /></el-form-item>
            <el-form-item label="有效秒数"><el-input-number v-model="tempPasswordForm.duration_seconds" :min="1" style="width:100%" /></el-form-item>
            <el-form-item label="描述"><el-input v-model="tempPasswordForm.description" /></el-form-item>
          </el-form>
          <el-table :data="tempPasswords" size="small" max-height="220">
            <el-table-column prop="Password" label="密码" />
            <el-table-column prop="Description" label="描述" />
            <el-table-column label="操作" width="90">
              <template #default="scope">
                <el-button type="danger" link @click="deleteTempPassword(scope.row.Password)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>

        <el-card class="mt-16">
          <template #header>
            <div class="card-header-inline">
              <b>Query 登录管理</b>
              <el-button type="primary" size="small" @click="createQueryLogin">生成</el-button>
            </div>
          </template>
          <el-form :model="queryForm" inline size="small">
            <el-form-item label="用户DBID"><el-input-number v-model="queryForm.client_db_id" :min="1" /></el-form-item>
            <el-form-item label="SID"><el-input-number v-model="queryForm.server_id" :min="0" /></el-form-item>
          </el-form>
          <el-table :data="queryLogins" size="small" max-height="220">
            <el-table-column prop="ClientDBID" label="DBID" width="80" />
            <el-table-column prop="LoginName" label="登录名" />
            <el-table-column label="操作" width="90">
              <template #default="scope">
                <el-button type="danger" link @click="deleteQueryLogin(scope.row.ClientDBID)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import api from '../api'

const serverForm = reactive({
  name: '',
  welcome_message: '',
  password: '',
  max_clients: 32,
  host_message: '',
  host_message_mode: 0
})

const tempPasswordForm = reactive({
  password: '',
  description: '',
  duration_seconds: 3600
})

const queryForm = reactive({ client_db_id: 0, server_id: 0 })
const groupClientForm = reactive({ server_group_id: 0, client_db_id: 0 })
const tempPasswords = ref([])
const queryLogins = ref([])

const loadServerInfo = async () => {
  const { data } = await api.get('/api/v1/server')
  serverForm.name = data.name || ''
  serverForm.max_clients = data.max_clients || 32
}

const loadTempPasswords = async () => {
  const { data } = await api.get('/api/v1/server/temp-passwords')
  tempPasswords.value = data.data || []
}

const loadQueryLogins = async () => {
  const { data } = await api.get('/api/v1/server/query-logins')
  queryLogins.value = data.data || []
}

const initData = async () => {
  await Promise.all([loadServerInfo(), loadTempPasswords(), loadQueryLogins()])
}

const saveServerSettings = async () => {
  await api.put('/api/v1/server/settings', serverForm)
  serverForm.password = ''
  ElMessage.success('服务器配置已保存')
}

const createTempPassword = async () => {
  if (!tempPasswordForm.password) return ElMessage.warning('请先输入密码')
  await api.post('/api/v1/server/temp-passwords', tempPasswordForm)
  ElMessage.success('临时密码创建成功')
  tempPasswordForm.password = ''
  tempPasswordForm.description = ''
  await loadTempPasswords()
}

const deleteTempPassword = async (pw) => {
  await api.delete(`/api/v1/server/temp-passwords/${encodeURIComponent(pw)}`)
  ElMessage.success('临时密码已删除')
  await loadTempPasswords()
}

const createQueryLogin = async () => {
  if (!queryForm.client_db_id) return ElMessage.warning('请输入用户DBID')
  const { data } = await api.post('/api/v1/server/query-logins', queryForm)
  const login = data.data?.LoginName || '-'
  const password = data.data?.Password || '-'
  await ElMessageBox.alert(`登录名：${login}\n密码：${password}`, 'Query 凭据')
  await loadQueryLogins()
}

const deleteQueryLogin = async (cldbid) => {
  await api.delete(`/api/v1/server/query-logins/${cldbid}`)
  ElMessage.success('Query 登录已删除')
  await loadQueryLogins()
}

const changeGroupClient = async (action) => {
  const path = action === 'add' ? '/api/v1/server/group-client/add' : '/api/v1/server/group-client/remove'
  await api.post(path, groupClientForm)
  ElMessage.success(action === 'add' ? '成员添加成功' : '成员移除成功')
}

onMounted(initData)
</script>

<style scoped>
.server-manager { display:flex; flex-direction:column; gap:16px; }
.page-header h2 { margin:0; }
.page-header p { margin:6px 0 0; color:#7a7f8a; }
.header-main { display:flex; justify-content:space-between; align-items:center; }
.content-grid { margin-bottom: 12px; }
.card-header-inline { display:flex; justify-content:space-between; align-items:center; }
.mt-16 { margin-top: 16px; }
</style>
