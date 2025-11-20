import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import MainLayout from '../views/MainLayout.vue'
import Dashboard from '../views/Dashboard.vue'
import ChannelManager from '../views/ChannelManager.vue'
import GroupManager from '../views/GroupManager.vue'
import BanManager from "../views/BanManager.vue";

const router = createRouter({
    history: createWebHistory(),
    routes: [
        { path: '/login', component: Login },
        { path: '/register', component: Register },
        {
            path: '/',
            component: MainLayout,
            meta: { requiresAuth: true },
            children: [
                { path: '', redirect: '/dashboard' },
                { path: 'dashboard', component: Dashboard, meta: { title: '仪表盘' } },
                { path: 'channels', component: ChannelManager, meta: { title: '频道管理' } },
                { path: 'groups', component: GroupManager, meta: { title: '权限组管理' } },
                { path: 'bans', component: BanManager, meta: { title: '封禁管理' } },
            ]
        }
    ]
})

router.beforeEach((to, from, next) => {
    const token = localStorage.getItem('token')
    if (to.meta.requiresAuth && !token) {
        next('/login')
    } else {
        next()
    }
})

export default router