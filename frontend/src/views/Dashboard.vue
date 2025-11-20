<template>
  <el-container class="layout-container">
    <el-main>
      <el-row :gutter="20" class="mb-4">
        <el-col :span="16">
          <el-card v-loading="loadingServer">
            <template #header>
              <div class="card-header">
                <span>服务器状态</span>
                <el-tag :type="serverInfo.status === 'online' ? 'success' : 'danger'">
                  {{ serverInfo.status || 'Offline' }}
                </el-tag>
              </div>
            </template>
            <el-descriptions border :column="2">
              <el-descriptions-item label="服务器名">{{ serverInfo.name }}</el-descriptions-item>
              <el-descriptions-item label="版本">{{ serverInfo.version }} on {{ serverInfo.platform }}</el-descriptions-item>
              <el-descriptions-item label="在线人数">
                {{ serverInfo.clients_online }} / {{ serverInfo.max_clients }}
              </el-descriptions-item>
              <el-descriptions-item label="运行时间">{{ formatUptime(serverInfo.uptime) }}</el-descriptions-item>
            </el-descriptions>
          </el-card>
        </el-col>

        <el-col :span="8">
          <el-card class="box-card" style="height: 100%">
            <template #header>
              <div class="card-header">
                <span>快捷管理</span>
              </div>
            </template>
            <div class="action-buttons">
              <el-button type="primary" @click="dialogs.createChannel = true">新建频道</el-button>
              <el-button type="warning" @click="dialogs.permission = true">权限操作</el-button>
              <el-button type="success" @click="dialogs.token = true">生成密钥</el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="16">
          <el-card class="box-card">
            <template #header>
              <div class="card-header">
                <span>在线用户 ({{ clients.length }})</span>
                <el-button type="primary" size="small" @click="fetchClients" :icon="Refresh">刷新</el-button>
              </div>
            </template>
            <el-table :data="clients" style="width: 100%" height="400" stripe>
              <el-table-column prop="Nickname" label="昵称" width="180" />
              <el-table-column prop="ID" label="CLID" width="80" />
              <el-table-column prop="DatabaseID" label="DBID" width="80" />
              <el-table-column prop="ChannelID" label="频道ID" width="80" />
              <el-table-column prop="Type" label="类型">
                <template #default="scope">
                  <el-tag size="small" v-if="scope.row.Type === 1">Query</el-tag>
                  <el-tag size="small" type="info" v-else>User</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" align="right">
                <template #default="scope">
                  <el-button size="small" type="danger" @click="openKickDialog(scope.row)">踢出</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-card>

          <el-card class="mt-4">
            <el-input v-model="broadcastMsg" placeholder="发送全服公告..." class="input-with-select">
              <template #append>
                <el-button @click="sendBroadcast">发送</el-button>
              </template>
            </el-input>
          </el-card>
        </el-col>

        <el-col :span="8">
          <el-card class="box-card" body-style="padding: 0">
            <template #header>
              <div class="card-header">
                <span>实时日志</span>
                <div class="live-indicator"></div>
              </div>
            </template>
            <div class="log-console" ref="logBox">
              <div v-for="(log, index) in logs" :key="index" class="log-line">
                <span class="log-time">[{{ log.time }}]</span>
                <span :class="'log-type-' + log.type">{{ log.type.toUpperCase() }}:</span>
                {{ log.data }}
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-main>

    <el-dialog v-model="kickDialogVisible" title="踢出用户" width="30%">
      <span>确定要踢出用户 <b>{{ currentKickUser?.Nickname }}</b> 吗？</span>
      <el-input v-model="kickReason" placeholder="踢出理由" class="mt-2" />
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="kickDialogVisible = false">取消</el-button>
          <el-button type="danger" @click="confirmKick">确定踢出</el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog v-model="dialogs.createChannel" title="新建频道" width="400px">
      <el-form :model="forms.channel" label-width="80px">
        <el-form-item label="频道名称">
          <el-input v-model="forms.channel.name" placeholder="请输入频道名" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="forms.channel.password" placeholder="留空则无密码" show-password />
        </el-form-item>
        <el-form-item label="话题">
          <el-input v-model="forms.channel.topic" placeholder="频道话题/公告" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogs.createChannel = false">取消</el-button>
        <el-button type="primary" @click="handleCreateChannel">创建</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="dialogs.permission" title="权限操作 (高级)" width="500px">
      <el-alert title="修改权限需谨慎，错误的数值可能导致用户无法进入频道或管理功能失效。" type="warning" show-icon :closable="false" class="mb-4"/>
      <el-form :model="forms.perm" label-width="100px">
        <el-form-item label="操作对象">
          <el-radio-group v-model="forms.perm.scope">
            <el-radio-button label="channel">频道</el-radio-button>
            <el-radio-button label="servergroup">服务器组</el-radio-button>
            <el-radio-button label="client">客户端(DB)</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="目标 ID">
          <el-input-number v-model="forms.perm.targetId" :min="0" style="width: 100%" placeholder="例如 ChannelID 或 GroupID" />
          <div class="form-tip">请输入频道ID、服务器组ID或用户数据库ID</div>
        </el-form-item>
        <el-form-item label="权限名">
          <el-input v-model="forms.perm.name" placeholder="如: i_channel_needed_join_power" />
          <div class="form-tip">
            <a href="https://jk.b.g8.sh/" target="_blank" style="color: #409eff">查看 TS3 权限列表</a>
          </div>
        </el-form-item>
        <el-form-item label="权限值">
          <el-input-number v-model="forms.perm.value" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogs.permission = false">取消</el-button>
        <el-button type="primary" @click="handleUpdatePerm">执行修改</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="dialogs.token" title="生成权限密钥 (Token)" width="400px">
      <el-form :model="forms.token" label-width="80px">
        <el-form-item label="类型">
          <el-select v-model="forms.token.type">
            <el-option label="服务器组 (Server Group)" :value="0" />
            <el-option label="频道组 (Channel Group)" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item label="组 ID">
          <el-input-number v-model="forms.token.groupId" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="频道 ID" v-if="forms.token.type === 1">
          <el-input-number v-model="forms.token.channelId" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="forms.token.description" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogs.token = false">取消</el-button>
        <el-button type="success" @click="handleCreateToken">生成</el-button>
      </template>
    </el-dialog>

  </el-container>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { Refresh } from '@element-plus/icons-vue'
import api from '../api'
import { fetchEventSource } from '@microsoft/fetch-event-source'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()

// --- 基础数据 ---
const serverInfo = ref({})
const clients = ref([])
const logs = ref([])
const loadingServer = ref(false)
const broadcastMsg = ref('')
const logBox = ref(null)

// --- 弹窗控制 ---
const kickDialogVisible = ref(false)
const currentKickUser = ref(null)
const kickReason = ref('Kicked by WebPanel')

const dialogs = reactive({
  createChannel: false,
  permission: false,
  token: false
})

// --- 表单数据 ---
const forms = reactive({
  channel: { name: '', password: '', topic: '' },
  perm: { scope: 'channel', targetId: 0, name: '', value: 75 },
  token: { type: 0, groupId: 0, channelId: 0, description: '' }
})

// --- 核心逻辑 ---

// 1. 新建频道
const handleCreateChannel = async () => {
  if (!forms.channel.name) return ElMessage.warning('请输入频道名称')
  try {
    await api.post('/api/v1/channel/create', {
      channel_name: forms.channel.name,
      channel_password: forms.channel.password,
      channel_topic: forms.channel.topic
    })
    ElMessage.success('频道创建成功')
    dialogs.createChannel = false
    // 可选：刷新服务器信息或手动通知用户
  } catch (e) {}
}

// 2. 权限操作
const handleUpdatePerm = async () => {
  const { scope, targetId, name, value } = forms.perm
  if (!name || !targetId) return ElMessage.warning('请填写完整参数')

  let url = ''
  // 根据不同范围构造请求 URL (需后端对应实现)
  if (scope === 'channel') {
    url = `/api/v1/channel/${targetId}/permission`
  } else if (scope === 'servergroup') {
    url = `/api/v1/servergroup/${targetId}/permission`
  } else if (scope === 'client') {
    url = `/api/v1/clientdb/${targetId}/permission`
  }

  try {
    await api.post(url, { perm_name: name, perm_value: value })
    ElMessage.success('权限修改指令已发送')
    dialogs.permission = false
  } catch (e) {}
}

// 3. 生成 Token
const handleCreateToken = async () => {
  try {
    const res = await api.post('/api/v1/token/create', forms.token)
    // 假设后端返回 { token: "xxx" }
    if (res.data.token) {
      ElMessageBox.alert(res.data.token, '生成成功 (请复制)', {
        confirmButtonText: '关闭'
      })
      dialogs.token = false
    }
  } catch (e) {}
}

// --- 原有功能保持不变 ---

const formatUptime = (seconds) => {
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  return `${h}小时 ${m}分钟`
}

const fetchServerInfo = async () => {
  loadingServer.value = true
  try {
    const res = await api.get('/api/v1/server')
    serverInfo.value = res.data
  } catch(e) {}
  loadingServer.value = false
}

const fetchClients = async () => {
  try {
    const res = await api.get('/api/v1/clients')
    clients.value = res.data.data || res.data
  } catch(e) {}
}

const openKickDialog = (user) => {
  currentKickUser.value = user
  kickReason.value = '违规操作'
  kickDialogVisible.value = true
}

const confirmKick = async () => {
  if (!currentKickUser.value) return
  try {
    await api.post(`/api/v1/client/${currentKickUser.value.ID}/kick`, {
      reason: kickReason.value
    })
    ElMessage.success('踢出成功')
    kickDialogVisible.value = false
    fetchClients()
  } catch (e) {}
}

const sendBroadcast = async () => {
  if(!broadcastMsg.value) return
  try {
    await api.post('/api/v1/broadcast', { message: broadcastMsg.value })
    ElMessage.success('发送成功')
    broadcastMsg.value = ''
  } catch(e) {}
}

const logout = () => {
  localStorage.removeItem('token')
  router.push('/login')
}

// SSE 相关
const ctrl = new AbortController()
const initSSE = () => {
  const token = localStorage.getItem('token')
  fetchEventSource('http://localhost:8080/api/v1/events/stream', {
    method: 'GET',
    headers: { Authorization: `Bearer ${token}` },
    signal: ctrl.signal,
    onmessage(msg) {
      if (!msg.data) return
      logs.value.push({
        time: new Date().toLocaleTimeString(),
        type: msg.event || 'info',
        data: msg.data
      })
      nextTick(() => {
        if(logBox.value) logBox.value.scrollTop = logBox.value.scrollHeight
      })
      if (msg.event === 'enter') fetchClients()
    },
    onerror(err) { console.log('SSE Error:', err) }
  })
}

onMounted(() => {
  fetchServerInfo()
  fetchClients()
  initSSE()
})

onUnmounted(() => {
  ctrl.abort()
})
</script>

<style scoped>
.layout-container { height: 100vh; background: #f5f7fa; }
.header { background: #fff; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid #dcdfe6; padding: 0 20px; }
.logo { font-weight: bold; font-size: 18px; color: #409eff; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.mb-4 { margin-bottom: 16px; }
.mt-2 { margin-top: 8px; }
.mt-4 { margin-top: 16px; }

/* 快捷管理按钮 */
.action-buttons { display: flex; gap: 10px; flex-wrap: wrap; }

/* 表单提示 */
.form-tip { font-size: 12px; color: #909399; margin-top: 4px; line-height: 1.2; }

/* 日志控制台样式 */
.log-console {
  background: #1e1e1e;
  color: #d4d4d4;
  height: 500px;
  overflow-y: auto;
  padding: 10px;
  font-family: 'Consolas', monospace;
  font-size: 12px;
  text-align: left;
}
.log-line { margin-bottom: 4px; border-bottom: 1px solid #333; padding-bottom: 2px; }
.log-time { color: #569cd6; margin-right: 8px; }
.log-type-message { color: #ce9178; }
.log-type-enter { color: #6a9955; }
.live-indicator { width: 10px; height: 10px; background: #67c23a; border-radius: 50%; animation: blink 2s infinite; }

@keyframes blink { 0% { opacity: 1; } 50% { opacity: 0.4; } 100% { opacity: 1; } }
</style>