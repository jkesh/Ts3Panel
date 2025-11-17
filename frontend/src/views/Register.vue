<template>
  <el-container class="layout-container">
    <el-header class="header">
      <div class="logo">TeamSpeak 3 Manager</div>
      <div class="user-info">
        <el-button type="danger" size="small" @click="logout">退出登录</el-button>
      </div>
    </el-header>

    <el-main>
      <el-row :gutter="20" class="mb-4">
        <el-col :span="24">
          <el-card v-loading="loadingServer">
            <template #header>
              <div class="card-header">
                <span>服务器状态</span>
                <el-tag :type="serverInfo.status === 'online' ? 'success' : 'danger'">
                  {{ serverInfo.status || 'Offline' }}
                </el-tag>
              </div>
            </template>
            <el-descriptions border>
              <el-descriptions-item label="服务器名">{{ serverInfo.name }}</el-descriptions-item>
              <el-descriptions-item label="版本">{{ serverInfo.version }} on {{ serverInfo.platform }}</el-descriptions-item>
              <el-descriptions-item label="在线人数">
                {{ serverInfo.clients_online }} / {{ serverInfo.max_clients }}
              </el-descriptions-item>
              <el-descriptions-item label="运行时间">{{ formatUptime(serverInfo.uptime) }}</el-descriptions-item>
            </el-descriptions>
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
            <el-table :data="clients" style="width: 100%" height="500" stripe>
              <el-table-column prop="Nickname" label="昵称" width="180" />
              <el-table-column prop="ID" label="Client ID" width="100" />
              <el-table-column prop="ChannelID" label="频道 ID" width="100" />
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
  </el-container>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { Refresh } from '@element-plus/icons-vue'
import api from '../api'
import { fetchEventSource } from '@microsoft/fetch-event-source'
import { ElMessage } from 'element-plus'

const router = useRouter()

// 数据状态
const serverInfo = ref({})
const clients = ref([])
const logs = ref([])
const loadingServer = ref(false)
const broadcastMsg = ref('')

// 踢人相关
const kickDialogVisible = ref(false)
const currentKickUser = ref(null)
const kickReason = ref('Kicked by WebPanel')

// Ref for log scroll
const logBox = ref(null)

// 格式化时间
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
    // 注意：你的后端返回结构可能是 { data: [...] } 或直接 [...]，请根据实际调整
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
    fetchClients() // 刷新列表
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

// SSE 连接控制器
const ctrl = new AbortController()

const initSSE = () => {
  const token = localStorage.getItem('token')
  // 使用 Microsoft 的 fetchEventSource 以支持 Headers (EventSource 原生不支持 Authorization Header)
  fetchEventSource('http://localhost:8080/api/v1/events/stream', {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${token}`,
    },
    signal: ctrl.signal,
    onmessage(msg) {
      // 处理心跳或其他非 JSON 数据
      if (!msg.data) return

      // 添加日志
      logs.value.push({
        time: new Date().toLocaleTimeString(),
        type: msg.event || 'info', // event 字段对应后端的 type
        data: msg.data
      })

      // 自动滚动到底部
      nextTick(() => {
        if(logBox.value) logBox.value.scrollTop = logBox.value.scrollHeight
      })

      // 如果有人进入，自动刷新列表
      if (msg.event === 'enter') {
        fetchClients()
      }
    },
    onerror(err) {
      console.log('SSE Error:', err)
      // 可以在这里处理重连逻辑
    }
  })
}

onMounted(() => {
  fetchServerInfo()
  fetchClients()
  initSSE()
})

onUnmounted(() => {
  ctrl.abort() // 组件销毁时断开 SSE
})
</script>

<style scoped>
.layout-container { height: 100vh; background: #f5f7fa; }
.header { background: #fff; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid #dcdfe6; }
.logo { font-weight: bold; font-size: 18px; color: #409eff; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.mb-4 { margin-bottom: 16px; }
.mt-2 { margin-top: 8px; }
.mt-4 { margin-top: 16px; }

/* 日志控制台样式 */
.log-console {
  background: #1e1e1e;
  color: #d4d4d4;
  height: 500px;
  overflow-y: auto;
  padding: 10px;
  font-family: 'Consolas', monospace;
  font-size: 12px;
}
.log-line { margin-bottom: 4px; border-bottom: 1px solid #333; padding-bottom: 2px; }
.log-time { color: #569cd6; margin-right: 8px; }
.log-type-message { color: #ce9178; }
.log-type-enter { color: #6a9955; }
.live-indicator { width: 10px; height: 10px; background: #67c23a; border-radius: 50%; animation: blink 2s infinite; }

@keyframes blink { 0% { opacity: 1; } 50% { opacity: 0.4; } 100% { opacity: 1; } }
</style>