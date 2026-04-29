<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '../../stores/app'
import { useRoutesStore } from '../../stores/routes'
import { useUpstreamsStore } from '../../stores/upstreams'
import { useI18n } from 'vue-i18n'
import { computed } from 'vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const routesStore = useRoutesStore()
const upstreamsStore = useUpstreamsStore()

const navSections = computed(() => [
  {
    label: t('sidebar.section.monitor'),
    items: [
      { key: 'dashboard', label: t('sidebar.nav.overview'), icon: '\u{1F4CA}', path: '/dashboard' },
      { key: 'metrics', label: t('sidebar.nav.metrics'), icon: '\u{1F4C8}', path: '/metrics' }
    ]
  },
  {
    label: t('sidebar.section.management'),
    items: [
      { key: 'routes', label: t('sidebar.nav.routes'), icon: '\u{1F6E4}\u{FE0F}', path: '/routes', badge: routesStore.routes.length || undefined },
      { key: 'upstreams', label: t('sidebar.nav.upstreams'), icon: '\u{1F5A5}\u{FE0F}', path: '/upstreams', badge: upstreamsStore.upstreams.length || undefined },
      { key: 'circuit', label: t('sidebar.nav.circuitBreakers'), icon: '⚡', path: '/circuit' }
    ]
  },
  {
    label: t('sidebar.section.security'),
    items: [
      { key: 'ratelimit', label: t('sidebar.nav.rateLimiting'), icon: '\u{1F6A6}', path: '/ratelimit' }
    ]
  },
  {
    label: t('sidebar.section.system'),
    items: [
      { key: 'settings', label: t('sidebar.nav.settings'), icon: '⚙\u{FE0F}', path: '/settings' }
    ]
  }
])

const currentPath = computed(() => route.path)
const collapsed = computed(() => appStore.sidebarCollapsed)

function navigate(path: string) {
  router.push(path)
  if (appStore.isMobile) appStore.sidebarCollapsed = true
}
</script>

<template>
  <aside class="sidebar" :class="{ collapsed }">
    <div class="sidebar-logo">
      <div class="logo-icon">A</div>
      <div v-if="!collapsed">
        <h1>{{ t('sidebar.brand') }}</h1>
        <span>{{ t('sidebar.version') }}</span>
      </div>
    </div>

    <nav class="sidebar-nav">
      <div v-for="section in navSections" :key="section.label">
        <div class="nav-section" v-if="!collapsed">{{ section.label }}</div>
        <div
          v-for="item in section.items"
          :key="item.key"
          class="nav-item"
          :class="{ active: currentPath === item.path }"
          @click="navigate(item.path)"
          :title="collapsed ? item.label : ''"
        >
          <span class="icon">{{ item.icon }}</span>
          <span v-if="!collapsed" class="label">{{ item.label }}</span>
          <span v-if="item.badge && !collapsed" class="badge">{{ item.badge }}</span>
        </div>
      </div>
    </nav>
  </aside>
</template>

<style scoped>
.sidebar {
  position: fixed; left: 0; top: 0; bottom: 0;
  width: var(--sidebar-width); background: var(--bg-secondary);
  border-right: 1px solid var(--border); display: flex;
  flex-direction: column; z-index: 100; transition: width 0.3s;
}
.sidebar.collapsed { width: 64px; }
.sidebar-logo {
  padding: 16px 20px; border-bottom: 1px solid var(--border);
  display: flex; align-items: center; gap: 12px;
}
.logo-icon {
  width: 36px; height: 36px; background: var(--accent);
  border-radius: 8px; display: flex; align-items: center;
  justify-content: center; font-size: 18px; font-weight: 700; color: #fff;
}
.sidebar-logo h1 { font-size: 18px; font-weight: 700; }
.sidebar-logo span { font-size: 11px; color: var(--text-muted); display: block; }
.sidebar-nav { flex: 1; padding: 12px 0; overflow-y: auto; }
.nav-section {
  padding: 8px 16px 4px; font-size: 10px; font-weight: 600;
  text-transform: uppercase; letter-spacing: 1px; color: var(--text-muted);
}
.nav-item {
  display: flex; align-items: center; gap: 12px;
  padding: 10px 20px; cursor: pointer; color: var(--text-secondary);
  transition: all 0.15s; border-left: 3px solid transparent; font-size: 14px;
}
.nav-item:hover { background: var(--bg-hover); color: var(--text-primary); }
.nav-item.active {
  background: rgba(59,130,246,0.1); color: var(--accent);
  border-left-color: var(--accent);
}
.icon { font-size: 18px; width: 24px; text-align: center; }
.badge {
  margin-left: auto; background: var(--danger); color: #fff;
  font-size: 10px; padding: 2px 6px; border-radius: 10px; font-weight: 600;
}
@media (max-width: 767px) {
  .sidebar { transform: translateX(-100%); }
  .sidebar:not(.collapsed) { transform: translateX(0); width: var(--sidebar-width); }
}
</style>
