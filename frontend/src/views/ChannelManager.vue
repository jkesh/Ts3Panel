<template>
  <div class="channel-manager">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>频道列表</span>
          <el-button type="primary" @click="dialogCreate = true">新建频道</el-button>
        </div>
      </template>

      <el-table :data="channels" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="ID" label="ID" width="80" />
        <el-table-column prop="Name" label="频道名称" />
        <el-table-column prop="TotalClients" label="在线人数" width="100" />
        <el-table-column prop="ParentID" label="父频道ID" width="100" />
        <el-table-column label="操作" width="250">
          <template #default="scope">
            <el-button size="small" @click="openPermDialog(scope.row)">权限</el-button>
            <el-popconfirm title="确定删除此频道吗？" @confirm="handleDelete(scope.row)">
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogCreate" title="新建频道" width="400px">
      <el-form :model="formCreate" label-width="80px">
        <el-form-item label="名称"><el-input v-model="formCreate.name" /></el-form-item>
        <el-form-item label="密码"><el-input v-model="formCreate.password" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createChannel">提交</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="dialogPerm" title="修改频道权限" width="400px">
      <el-form :model="formPerm" label-width="100px">
        <el-form-item label="权限名"><el-input v-model="formPerm.name" placeholder="如 i_channel_needed_join_power"/></el-form-item>
        <el-form-item label="值"><el-input-number v-model="formPerm.value" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button type="primary" @click="updatePerm">修改</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import api from '../api'
import { ElMessage } from 'element-plus'

const channels = ref([])
const loading = ref(false)
const dialogCreate = ref(false)
const dialogPerm = ref(false)

const formCreate = reactive({ name: '', password: '' })
const formPerm = reactive({ cid: 0, name: '', value: 0 })

const fetchChannels = async () => {
  loading.value = true
  try {
    const res = await api.get('/api/v1/channels')
    channels.value = res.data.data
  } catch (e) {}
  loading.value = false
}

const handleDelete = async (row) => {
  try {
    await api.delete(`/api/v1/channel/${row.ID}?force=1`)
    ElMessage.success('删除成功')
    fetchChannels()
  } catch (e) {}
}

const createChannel = async () => {
  try {
    await api.post('/api/v1/channel/create', { channel_name: formCreate.name, channel_password: formCreate.password })
    ElMessage.success('创建成功')
    dialogCreate.value = false
    fetchChannels()
  } catch(e) {}
}

const openPermDialog = (row) => {
  formPerm.cid = row.ID
  formPerm.name = ''
  formPerm.value = 0
  dialogPerm.value = true
}

const updatePerm = async () => {
  try {
    await api.post(`/api/v1/channel/${formPerm.cid}/permission`, { perm_name: formPerm.name, perm_value: formPerm.value })
    ElMessage.success('权限已修改')
    dialogPerm.value = false
  } catch(e) {}
}

onMounted(fetchChannels)
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>