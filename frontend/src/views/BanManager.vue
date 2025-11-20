<template>
  <div class="ban-manager">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>封禁列表 (Ban List)</span>
          <div>
            <el-button type="danger" plain @click="deleteAll" :icon="Delete">清空列表</el-button>
            <el-button type="primary" @click="dialogAdd = true" :icon="Plus">添加封禁</el-button>
          </div>
        </div>
      </template>

      <el-table :data="bans" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column label="对象 (Name/IP/UID)" min-width="200">
          <template #default="scope">
            <div v-if="scope.row.name">Name: {{ scope.row.name }}</div>
            <div v-if="scope.row.ip">IP: {{ scope.row.ip }}</div>
            <div v-if="scope.row.uid">UID: {{ scope.row.uid }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="理由" show-overflow-tooltip />
        <el-table-column label="时长">
          <template #default="scope">
            {{ scope.row.duration == 0 ? '永久' : (scope.row.duration / 60 + '分钟') }}
          </template>
        </el-table-column>
        <el-table-column prop="invoker" label="操作人" width="120" />
        <el-table-column label="操作" width="100" align="right">
          <template #default="scope">
            <el-button size="small" type="success" @click="unban(scope.row)">解封</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogAdd" title="添加封禁" width="500px">
      <el-form label-width="80px">
        <el-alert title="IP、UID、昵称 至少填写一项" type="info" show-icon :closable="false" class="mb-3"/>
        <el-form-item label="IP地址"><el-input v-model="form.ip" placeholder="可选" /></el-form-item>
        <el-form-item label="UID"><el-input v-model="form.uid" placeholder="用户唯一ID (可选)" /></el-form-item>
        <el-form-item label="昵称"><el-input v-model="form.name" placeholder="支持正则 (可选)" /></el-form-item>
        <el-form-item label="时长(秒)">
          <el-input-number v-model="form.time" :min="0" placeholder="0为永久" />
          <span class="ml-2 text-gray">0 = 永久</span>
        </el-form-item>
        <el-form-item label="理由"><el-input v-model="form.reason" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogAdd = false">取消</el-button>
        <el-button type="danger" @click="addBan">执行封禁</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import api from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete } from '@element-plus/icons-vue'

const bans = ref([])
const loading = ref(false)
const dialogAdd = ref(false)
const form = reactive({ ip: '', uid: '', name: '', time: 0, reason: 'WebPanel Ban' })

const fetchBans = async () => {
  loading.value = true
  try {
    const res = await api.get('/api/v1/bans')
    bans.value = res.data.data
  } catch (e) {}
  loading.value = false
}

const addBan = async () => {
  if (!form.ip && !form.uid && !form.name) return ElMessage.warning('请至少填写一个封禁对象')
  try {
    await api.post('/api/v1/ban', form)
    ElMessage.success('封禁成功')
    dialogAdd.value = false
    fetchBans()
  } catch(e) {}
}

const unban = async (row) => {
  try {
    await api.delete(`/api/v1/ban/${row.id}`)
    ElMessage.success('已解封')
    fetchBans()
  } catch(e) {}
}

const deleteAll = () => {
  ElMessageBox.confirm('确定要清空所有封禁记录吗？', '危险操作', {
    type: 'warning',
    confirmButtonText: '清空',
    cancelButtonText: '取消'
  }).then(async () => {
    try {
      await api.delete('/api/v1/bans/all')
      ElMessage.success('封禁列表已清空')
      fetchBans()
    } catch(e) {}
  })
}

onMounted(fetchBans)
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
.mb-3 { margin-bottom: 12px; }
.ml-2 { margin-left: 8px; }
.text-gray { color: #999; font-size: 12px; }
</style>