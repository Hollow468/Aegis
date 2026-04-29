<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoutesStore } from '../../stores/routes'
import { useConfirm } from '../../composables/useConfirm'
import { useI18n } from 'vue-i18n'
import StateTag from '../../components/common/StateTag.vue'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import { ElMessage } from 'element-plus'

const { t } = useI18n()
const routesStore = useRoutesStore()
const { visible, title, message, confirm, handleConfirm, handleCancel } = useConfirm()

onMounted(() => { routesStore.fetchRoutes() })

async function handleDelete(path: string) {
  const ok = await confirm(t('routes.confirm.title'), t('routes.confirm.message', { path }))
  if (!ok) return
  try {
    await routesStore.removeRoute(path)
    ElMessage.success(t('routes.success.deleted'))
  } catch {
    ElMessage.error(t('routes.error.deleteFailed'))
  }
}
</script>

<template>
  <div class="routes-page">
    <div class="panel">
      <div class="panel-header">
        <h3>{{ t('routes.header') }}</h3>
        <button class="btn btn-primary">{{ t('routes.addRoute') }}</button>
      </div>
      <div class="panel-body">
        <div v-if="routesStore.loading" class="loading">{{ t('common.loading') }}</div>
        <table v-else>
          <thead>
            <tr>
              <th>{{ t('routes.col.path') }}</th>
              <th>{{ t('routes.col.method') }}</th>
              <th>{{ t('routes.col.matchType') }}</th>
              <th>{{ t('routes.col.balancer') }}</th>
              <th>{{ t('routes.col.rateLimit') }}</th>
              <th>{{ t('routes.col.upstreams') }}</th>
              <th>{{ t('common.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="r in routesStore.routes" :key="r.path">
              <td><code>{{ r.path }}</code></td>
              <td><StateTag :state="r.method" variant="method" /></td>
              <td><StateTag :state="r.match_type" variant="match" /></td>
              <td>{{ r.balancer }}</td>
              <td>{{ r.rate_limit ? `${r.rate_limit.limit}/s` : '-' }}</td>
              <td>{{ r.upstreams?.length ?? 0 }}</td>
              <td>
                <button class="btn btn-secondary btn-sm">{{ t('common.edit') }}</button>
                <button class="btn btn-danger btn-sm" @click="handleDelete(r.path)">{{ t('common.delete') }}</button>
              </td>
            </tr>
            <tr v-if="routesStore.routes.length === 0">
              <td colspan="7" class="empty">{{ t('routes.empty') }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <ConfirmDialog :visible="visible" :title="title" :message="message" @confirm="handleConfirm" @cancel="handleCancel" />
  </div>
</template>

<style scoped>
.panel { background: var(--bg-secondary); border: 1px solid var(--border); border-radius: var(--radius); }
.panel-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; border-bottom: 1px solid var(--border); }
.panel-header h3 { font-size: 14px; font-weight: 600; }
.panel-body { padding: 20px; }
table { width: 100%; border-collapse: collapse; }
th, td { padding: 12px 16px; text-align: left; font-size: 13px; border-bottom: 1px solid var(--border); }
th { color: var(--text-secondary); font-weight: 600; font-size: 11px; text-transform: uppercase; background: var(--bg-tertiary); }
.loading { padding: 20px; text-align: center; color: var(--text-secondary); }
.empty { text-align: center; color: var(--text-muted); padding: 24px; }
.btn { display: inline-flex; align-items: center; gap: 6px; padding: 8px 16px; border-radius: 6px; font-size: 13px; font-weight: 600; cursor: pointer; border: none; }
.btn-primary { background: var(--accent); color: #fff; }
.btn-secondary { background: var(--bg-tertiary); color: var(--text-primary); border: 1px solid var(--border); }
.btn-danger { background: rgba(239,68,68,0.15); color: var(--danger); }
.btn-sm { padding: 4px 10px; font-size: 12px; }
</style>
