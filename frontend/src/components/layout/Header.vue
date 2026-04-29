<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '../../stores/app'
import { useAuthStore } from '../../stores/auth'
import { useI18n } from 'vue-i18n'
import { computed, ref, onMounted, onUnmounted } from 'vue'

const { t, locale } = useI18n()
const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()

const pageTitle = computed(() => {
  const key = route.meta?.titleKey as string
  return key ? t(key) : t('route.defaultTitle')
})

const uptime = ref(0)
let timer: ReturnType<typeof setInterval>

function formatUptime(seconds: number): string {
  const d = Math.floor(seconds / 86400)
  const h = Math.floor((seconds % 86400) / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  return `${d}d ${h}h ${m}m`
}

onMounted(() => { timer = setInterval(() => { uptime.value++ }, 1000) })
onUnmounted(() => { clearInterval(timer) })

function handleLogout() {
  authStore.logout()
  router.push('/login')
}

function switchLocale(lang: string) {
  locale.value = lang
  localStorage.setItem('aegis_locale', lang)
}
</script>

<template>
  <header class="header">
    <div class="header-left">
      <button class="toggle-btn" @click="appStore.toggleSidebar">&#9776;</button>
      <h2>{{ pageTitle }}</h2>
    </div>
    <div class="header-right">
      <span class="header-stat">{{ t('header.uptime') }} <strong>{{ formatUptime(uptime) }}</strong></span>
      <span class="header-stat">{{ t('header.port') }} <strong>:8080</strong></span>
      <div class="status-badge running">
        <span class="status-dot"></span>
        {{ t('header.status.running') }}
      </div>
      <el-dropdown @command="switchLocale" trigger="click">
        <span class="lang-btn">{{ locale === 'zh-cn' ? '中文' : 'EN' }}</span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="zh-cn" :class="{ active: locale === 'zh-cn' }">中文</el-dropdown-item>
            <el-dropdown-item command="en" :class="{ active: locale === 'en' }">English</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      <button v-if="authStore.isAuthenticated" class="logout-btn" @click="handleLogout" :title="t('header.logout')">
        &#x23FB;
      </button>
    </div>
  </header>
</template>

<style scoped>
.header {
  position: fixed; top: 0; left: var(--sidebar-width); right: 0;
  height: var(--header-height); background: var(--bg-secondary);
  border-bottom: 1px solid var(--border); display: flex;
  align-items: center; justify-content: space-between;
  padding: 0 24px; z-index: 99;
}
.header-left { display: flex; align-items: center; gap: 16px; }
.header-left h2 { font-size: 16px; font-weight: 600; }
.toggle-btn {
  background: none; border: none; color: var(--text-secondary);
  font-size: 20px; cursor: pointer; padding: 4px 8px; border-radius: 4px;
}
.toggle-btn:hover { background: var(--bg-hover); color: var(--text-primary); }
.header-right { display: flex; align-items: center; gap: 16px; }
.header-stat { font-size: 12px; color: var(--text-secondary); }
.header-stat strong { color: var(--text-primary); }
.status-badge {
  display: flex; align-items: center; gap: 6px;
  padding: 4px 12px; border-radius: 20px; font-size: 12px; font-weight: 600;
}
.status-badge.running { background: rgba(34,197,94,0.15); color: var(--success); }
.status-dot {
  width: 8px; height: 8px; border-radius: 50%;
  background: var(--success); animation: pulse 2s infinite;
}
.logout-btn {
  background: none; border: 1px solid var(--border); color: var(--text-secondary);
  width: 32px; height: 32px; border-radius: 6px; cursor: pointer;
  display: flex; align-items: center; justify-content: center; font-size: 14px;
}
.logout-btn:hover { background: rgba(239,68,68,0.15); color: var(--danger); border-color: rgba(239,68,68,0.3); }
.lang-btn {
  font-size: 12px; font-weight: 600; color: var(--text-secondary);
  cursor: pointer; padding: 4px 8px; border-radius: 4px;
  border: 1px solid var(--border); background: var(--bg-tertiary);
}
.lang-btn:hover { color: var(--text-primary); border-color: var(--accent); }
@keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.5; } }
</style>
