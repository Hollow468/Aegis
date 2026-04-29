<script setup lang="ts">
import Sidebar from '../components/layout/Sidebar.vue'
import Header from '../components/layout/Header.vue'
import { useAppStore } from '../stores/app'

const appStore = useAppStore()
</script>

<template>
  <div class="layout" :class="{ collapsed: appStore.sidebarCollapsed }">
    <Sidebar />
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
  display: flex;
  height: 100vh;
  background: var(--bg-primary);
}

.layout-main {
  flex: 1;
  margin-left: var(--sidebar-width);
  display: flex;
  flex-direction: column;
  transition: margin-left 0.3s;
}

.layout.collapsed .layout-main {
  margin-left: 64px;
}

.layout-content {
  flex: 1;
  margin-top: var(--header-height);
  padding: 24px;
  overflow-y: auto;
}
</style>
