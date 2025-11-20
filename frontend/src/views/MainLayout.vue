<template>
  <el-container class="layout-container">
    <el-aside width="220px" class="aside">
      <div class="logo">TS3 面板</div>
      <el-menu
          router
          :default-active="activePath"
          background-color="#1e1e1e"
          text-color="#fff"
          active-text-color="#409eff"
          class="el-menu-vertical"
      >
        <el-menu-item index="/dashboard">
          <el-icon><Odometer /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/channels">
          <el-icon><Monitor /></el-icon>
          <span>频道管理</span>
        </el-menu-item>
        <el-menu-item index="/groups">
          <el-icon><UserFilled /></el-icon>
          <span>权限组管理</span>
        </el-menu-item>
        <el-menu-item index="/bans">
          <el-icon><Lock /></el-icon>
          <span>封禁管理</span>
        </el-menu-item>
      </el-menu>

      <div class="logout-area">
        <el-button type="danger" plain style="width: 100%" @click="logout">退出登录</el-button>
      </div>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="breadcrumb">
          {{ currentRouteName }}
        </div>
        <div class="user">Admin</div>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Odometer, Monitor, UserFilled, Lock } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()

const activePath = computed(() => route.path)
const currentRouteName = computed(() => route.meta.title || 'TS3 Manager')

const logout = () => {
  localStorage.removeItem('token')
  router.push('/login')
}
</script>

<style scoped>
.layout-container { height: 100vh; }
.aside { background-color: #1e1e1e; color: #fff; display: flex; flex-direction: column; }
.logo { height: 60px; line-height: 60px; text-align: center; font-weight: bold; font-size: 20px; border-bottom: 1px solid #333; }
.el-menu-vertical { border-right: none; flex: 1; }
.logout-area { padding: 20px; border-top: 1px solid #333; }
.header { background: #fff; border-bottom: 1px solid #eee; display: flex; align-items: center; justify-content: space-between; padding: 0 20px; }
.breadcrumb { font-size: 16px; font-weight: 500; }
</style>