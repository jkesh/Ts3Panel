<template>
  <div class="bot-manager">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>TS3AudioBot 实例管理</span>
          <div class="header-actions">
            <el-switch v-model="autoRefresh" active-text="自动刷新" style="margin-right: 10px"/>
            <el-button type="primary" @click="dialogAdd = true">添加机器人</el-button>
          </div>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="8" v-for="bot in bots" :key="bot.ID" class="mb-4">
          <el-card shadow="hover" :class="['bot-card', 'status-' + (bot.real_status || 'offline').toLowerCase()]">
            <template #header>
              <div class="bot-card-header">
                <span class="bot-name">{{ bot.name }} #{{ bot.bot_id }}</span>
                <el-tag size="small" :type="getStatusType(bot.real_status)">
                  {{ bot.real_status || 'Unknown' }}
                </el-tag>
              </div>
            </template>

            <div class="control-area">
              <div class="netease-panel mb-2">
                <el-input
                    v-model="neteaseInputs[bot.ID]"
                    placeholder="网易云点歌 (歌名/ID/链接)"
                    size="small"
                    class="input-with-select"
                >
                  <template #prepend>
                    <el-icon><Headset /></el-icon>
                  </template>
                  <template #append>
                    <el-button @click="sendCommand(bot, 'netease', neteaseInputs[bot.ID])">搜索并播放</el-button>
                  </template>
                </el-input>
              </div>

              <el-divider content-position="left">常规播放 (Youtube/直链)</el-divider>

              <el-input
                  v-model="playUrls[bot.ID]"
                  placeholder="http://..."
                  size="small"
                  class="mb-2"
              ></el-input>
              <div v-if="bot.real_status === 'Offline'">
                <el-button type="primary" class="w-100" @click="connectBot(bot)">
                  启动并连接 ({{ bot.server_addr }})
                </el-button>
              </div>

              <div v-else>
                <el-input v-model="playUrls[bot.ID]" placeholder="音乐链接..." size="small" class="mb-2">
                  <template #append>
                    <el-button @click="sendCommand(bot, 'play', playUrls[bot.ID])">播放</el-button>
                  </template>
                </el-input>

                <el-button-group class="w-100 flex-buttons mb-2">
                  <el-button type="warning" size="small" @click="sendCommand(bot, 'pause')">暂停</el-button>
                  <el-button type="danger" size="small" @click="sendCommand(bot, 'stop')">停止</el-button>
                </el-button-group>

                <el-slider v-model="volumes[bot.ID]" :min="0" :max="100" size="small" @change="(val) => sendCommand(bot, 'volume', val)"/>
              </div>

              <div class="mt-2 text-right">
                <el-button type="danger" link size="small" @click="deleteBot(bot)">删除配置</el-button>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-card>

    <el-dialog v-model="dialogAdd" title="添加机器人" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="API 地址"><el-input v-model="form.api_url" placeholder="http://127.0.0.1:58913" /></el-form-item>
        <el-form-item label="Token"><el-input v-model="form.api_token" type="password"/></el-form-item>
        <el-form-item label="默认服务器"><el-input v-model="form.server_addr" placeholder="127.0.0.1"/></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="handleAddBot">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import api from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'

const bots = ref([])
const playUrls = reactive({})
const volumes = reactive({})
const dialogAdd = ref(false)
const autoRefresh = ref(true)
let timer = null

const form = reactive({ name: 'Bot 1', api_url: 'http://127.0.0.1:58913', api_token: '', server_addr: '127.0.0.1' })

const fetchBots = async () => {
  try {
    const res = await api.get('/api/v1/bots')
    bots.value = res.data.data
    // 同步音量初始值
    bots.value.forEach(b => {
      if (volumes[b.ID] === undefined) volumes[b.ID] = 50
    })
  } catch (e) {}
}

const getStatusType = (status) => {
  if (status === 'Playing') return 'success'
  if (status === 'Idle') return 'warning'
  if (status === 'Offline') return 'info'
  return ''
}

const handleAddBot = async () => {
  try {
    await api.post('/api/v1/bot', form)
    dialogAdd.value = false
    fetchBots()
  } catch(e) {}
}

const connectBot = async (bot) => {
  await sendCommand(bot, 'connect')
  setTimeout(fetchBots, 1000) // 1秒后刷新以获取新状态和ID
}

const sendCommand = async (bot, cmd, val) => {
  try {
    await api.post(`/api/v1/bot/${bot.ID}/command`, { command: cmd, value: val ? String(val) : "" })
    ElMessage.success('操作成功')
    fetchBots() // 操作后立即刷新状态
  } catch(e) {
    ElMessage.error('失败: ' + (e.response?.data?.error || '未知错误'))
  }
}

const deleteBot = (bot) => {
  ElMessageBox.confirm('确认删除？').then(async () => {
    await api.delete(`/api/v1/bot/${bot.ID}`)
    fetchBots()
  })
}

onMounted(() => {
  fetchBots()
  timer = setInterval(() => {
    if(autoRefresh.value) fetchBots()
  }, 3000) // 每3秒刷新一次状态
})

onUnmounted(() => clearInterval(timer))
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
.bot-card { transition: all 0.3s; }
.status-playing { border-color: #67C23A; }
.status-offline { filter: grayscale(100%); opacity: 0.8; }
.mb-4 { margin-bottom: 20px; }
.mb-2 { margin-bottom: 10px; }
.w-100 { width: 100%; }
.text-right { text-align: right; }
.flex-buttons { display: flex; }
.flex-buttons .el-button { flex: 1; }
</style>