<script setup lang="ts">
defineProps<{
  status: 'healthy' | 'down' | 'warning'
  label?: string
}>()

const colorMap = { healthy: 'dot-green', down: 'dot-red', warning: 'dot-yellow' }
</script>

<template>
  <span class="status-indicator">
    <span class="dot" :class="colorMap[status]" />
    <span v-if="label">{{ label }}</span>
  </span>
</template>

<style scoped>
.status-indicator { display: inline-flex; align-items: center; gap: 6px; font-size: 12px; font-weight: 600; }
.dot { width: 8px; height: 8px; border-radius: 50%; }
.dot-green { background: var(--success); animation: pulse 2s infinite; }
.dot-red { background: var(--danger); animation: pulse 1s infinite; }
.dot-yellow { background: var(--warning); animation: pulse 1.5s infinite; }
@keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.5; } }
</style>
