<script setup lang="ts">
import Sidebar from '../components/layout/Sidebar.vue'
import Header from '../components/layout/Header.vue'
import { useAppStore } from '../stores/app'
import { onMounted, onUnmounted } from 'vue'

const appStore = useAppStore()

function onResize() { appStore.checkMobile() }
onMounted(() => { onResize(); window.addEventListener('resize', onResize) })
onUnmounted(() => { window.removeEventListener('resize', onResize) })
</script>

<template>
  <div class="layout" :class="{ collapsed: appStore.sidebarCollapsed, mobile: appStore.isMobile }">
    <Sidebar />
    <div class="layout-overlay" v-if="appStore.isMobile && !appStore.sidebarCollapsed" @click="appStore.toggleSidebar" />
    <div class="layout-main">
      <Header />
      <main class="layout-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<style scoped>
.layout {
  display: flex; height: 100vh; background: var(--bg-primary);
}
.layout-main {
  flex: 1; margin-left: var(--sidebar-width);
  display: flex; flex-direction: column; transition: margin-left 0.3s;
}
.layout.collapsed .layout-main { margin-left: 64px; }
.layout.mobile .layout-main { margin-left: 0; }
.layout-content {
  flex: 1; margin-top: var(--header-height);
  padding: 24px; overflow-y: auto;
}
.layout-overlay {
  position: fixed; inset: 0; background: rgba(0,0,0,0.5);
  z-index: 99; display: none;
}
.layout.mobile .layout-overlay { display: block; }

@media (max-width: 767px) {
  .layout-content { padding: 16px; }
}
</style>
