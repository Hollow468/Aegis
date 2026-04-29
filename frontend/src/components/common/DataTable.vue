<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps<{
  loading?: boolean
  emptyText?: string
}>()
</script>

<template>
  <div class="data-table-wrapper">
    <div v-if="loading" class="table-loading">
      <el-icon class="is-loading" :size="24"><Loading /></el-icon>
      <span>{{ t('common.loading') }}</span>
    </div>
    <table v-else class="data-table">
      <thead>
        <tr><slot name="headers" /></tr>
      </thead>
      <tbody>
        <slot name="body" />
      </tbody>
    </table>
    <div v-if="!loading" class="table-empty">
      <slot name="empty">
        <span class="empty-text">{{ emptyText || t('common.noData') }}</span>
      </slot>
    </div>
  </div>
</template>

<script lang="ts">
import { Loading } from '@element-plus/icons-vue'
</script>

<style scoped>
.data-table-wrapper { position: relative; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table :deep(th),
.data-table :deep(td) {
  padding: 12px 16px; text-align: left; font-size: 13px;
  border-bottom: 1px solid var(--border);
}
.data-table :deep(th) {
  color: var(--text-secondary); font-weight: 600; font-size: 11px;
  text-transform: uppercase; background: var(--bg-tertiary);
}
.table-loading {
  display: flex; align-items: center; justify-content: center;
  gap: 8px; padding: 40px; color: var(--text-secondary);
}
.table-empty:empty { display: none; }
.empty-text { color: var(--text-muted); font-size: 13px; }
</style>
