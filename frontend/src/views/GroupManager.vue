<template>
  <div class="group-manager">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>服务器组管理 (仅显示常规组)</span>
          <el-button type="success" @click="dialogToken = true">生成 组Key</el-button>
        </div>
      </template>

      <el-table :data="filteredGroups" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="ID" label="SGID" width="80" />
        <el-table-column prop="Name" label="组名称" />
        <el-table-column prop="Type" label="类型" width="100">
          <template #default="scope">
            <el-tag>常规</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" align="right">
          <template #default="scope">
            <el-button size="small" type="primary" @click="openPermDialog(scope.row)">权限管理</el-button>

            <el-popconfirm title="确定删除此权限组吗？此操作不可逆！" @confirm="handleDelete(scope.row)">
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogPerm" :title="`管理权限: ${currentGroup?.Name}`" width="700px">
      <div class="perm-container">
        <div class="perm-list">
          <h4>现有权限</h4>
          <el-table
              :data="currentGroupPerms"
              height="300px"
              style="width: 100%"
              size="small"
              border
              @row-click="selectPerm"
          >
            <el-table-column label="权限名 (PermSID)" show-overflow-tooltip>
              <template #default="scope">
                <span v-if="scope.row.name">{{ scope.row.name }}</span>

                <span v-else-if="permMap[scope.row.permid]" style="color: #409eff;">
                  {{ permMap[scope.row.permid] }}
                </span>

                <span v-else style="color: #999;">ID: {{ scope.row.permid }}</span>
              </template>
            </el-table-column>

            <el-table-column prop="value" label="值" width="80" />
          </el-table>
        </div>

        <div class="perm-form">
          <h4>编辑 / 添加</h4>
          <el-form label-position="top">
            <el-form-item label="权限名">
              <el-input v-model="formPerm.name" placeholder="如 i_channel_join_power" />
            </el-form-item>
            <el-form-item label="权限值">
              <el-input-number v-model="formPerm.value" style="width: 100%" />
            </el-form-item>
            <el-button type="primary" style="width: 100%" @click="updatePerm">保存 / 修改</el-button>
            <div style="margin-top: 10px; font-size: 12px; color: #999;">
              * 点击左侧列表可快速回显原有值
            </div>
          </el-form>
        </div>
      </div>
    </el-dialog>

    <el-dialog v-model="dialogToken" title="生成邀请密钥 (Token)" width="400px">
      <el-form label-width="80px">
        <el-form-item label="组ID">
          <el-select v-model="tokenGid" placeholder="选择组" filterable>
            <el-option v-for="g in filteredGroups" :key="g.ID" :label="g.Name" :value="g.ID" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button type="primary" @click="createToken">生成</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import api from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'
import permMap from '../utils/permMap.js'

const groups = ref([])
const loading = ref(false)

// 权限弹窗相关
const dialogPerm = ref(false)
const currentGroup = ref(null)
const currentGroupPerms = ref([]) // 存储当前组的所有权限
const formPerm = reactive({ name: '', value: 0 })

// Token弹窗相关
const dialogToken = ref(false)
const tokenGid = ref(null)

// 1. 过滤 Query(2) 和 Template(0) 组，只显示 Regular(1)
const filteredGroups = computed(() => {
  return groups.value.filter(g => g.Type === 1)
})

const fetchGroups = async () => {
  loading.value = true
  try {
    const res = await api.get('/api/v1/servergroups')
    groups.value = res.data.data
  } catch(e) {}
  loading.value = false
}

// 2. 删除权限组逻辑
const handleDelete = async (row) => {
  try {
    await api.delete(`/api/v1/servergroup/${row.ID}`)
    ElMessage.success('删除成功')
    fetchGroups()
  } catch (e) {}
}

// 打开权限弹窗时，获取现有权限
const openPermDialog = async (row) => {
  currentGroup.value = row
  formPerm.name = ''
  formPerm.value = 0
  currentGroupPerms.value = []
  dialogPerm.value = true

  // 获取权限列表
  try {
    const res = await api.get(`/api/v1/servergroup/${row.ID}/permissions`)
    currentGroupPerms.value = res.data.data || []
  } catch(e) {}
}

// 3. 点击左侧表格行，回显数据到右侧表单
const selectPerm = (row) => {
  // 优先使用后端返回的 Name，如果为空，则尝试从 permMap 查找，最后才用 ID
  formPerm.name = row.name || permMap[row.permid] || row.permid.toString()
  formPerm.value = row.value
}

const updatePerm = async () => {
  if (!formPerm.name) return ElMessage.warning('请输入权限名')
  try {
    await api.post(`/api/v1/servergroup/${currentGroup.value.ID}/permission`, {
      perm_name: formPerm.name,
      perm_value: formPerm.value
    })
    ElMessage.success('权限已更新')
    // 刷新列表以显示最新值
    const res = await api.get(`/api/v1/servergroup/${currentGroup.value.ID}/permissions`)
    currentGroupPerms.value = res.data.data || []
  } catch(e) {}
}

const createToken = async () => {
  if(!tokenGid.value) return
  try {
    const res = await api.post('/api/v1/token/create', {
      type: 0,
      groupId: tokenGid.value,
      description: 'Created via WebPanel'
    })
    ElMessageBox.alert(res.data.token, '密钥生成成功', { confirmButtonText: '关闭' })
    dialogToken.value = false
  } catch(e) {}
}

onMounted(fetchGroups)
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }

/* 权限弹窗布局 */
.perm-container { display: flex; gap: 20px; }
.perm-list { flex: 1; border-right: 1px solid #eee; padding-right: 10px; }
.perm-form { width: 250px; }
h4 { margin-top: 0; margin-bottom: 10px; }
</style>