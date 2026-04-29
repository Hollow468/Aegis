<script setup lang="ts">
defineProps<{
  label: string
  value: string | number
  trend?: string
  trendType?: 'up' | 'down' | 'neutral'
}>()
</script>

<template>
  <div class="metric-card">
    <div class="label">{{ label }}</div>
    <div class="value">{{ value }}</div>
    <div v-if="trend" class="trend" :class="trendType">{{ trend }}</div>
  </div>
</template>

<style scoped>
.metric-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 20px;
}
.label {
  font-size: 12px; color: var(--text-secondary);
  text-transform: uppercase; margin-bottom: 8px;
}
.value { font-size: 28px; font-weight: 700; }
.trend { font-size: 12px; margin-top: 4px; }
.trend.up { color: var(--success); }
.trend.down { color: var(--danger); }
.trend.neutral { color: var(--text-secondary); }
</style>
