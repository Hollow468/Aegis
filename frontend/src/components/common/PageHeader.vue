<script setup lang="ts">
defineProps<{
  title: string
}>()
</script>

<template>
  <div class="page-header">
    <h3>{{ title }}</h3>
    <div class="actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<style scoped>
.page-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 16px 20px; border-bottom: 1px solid var(--border);
}
.page-header h3 { font-size: 14px; font-weight: 600; }
.actions { display: flex; gap: 8px; }
</style>
